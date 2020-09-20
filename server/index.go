package server

import (
	"encoding/json"
	"fmt"
	"github.com/parvez0/metric-ingestion/custom_logger"
	"github.com/parvez0/metric-ingestion/objects"
	sqlite_db "github.com/parvez0/metric-ingestion/sqlite-db"
	"io"
	"net"
	"net/http"
	"os"
)

var clog = custom_logger.NewLogger()
var db = sqlite_db.CreateDbConnection()
var statusMessages = map[int]string{
	500: "Internal server error",
	400: "Bad request",
}

// CreateServer creates a http service with default handlers
func CreateServer()  {
	table := "metrics"
	if os.Getenv("GO_ENV") == "testing"{
		table = "testing_table_metrics"
		clog.Debugf("server is running in testing environment creating a dummy table - %s", table)
	}

	// creating the schema and tables for first time initialization
	db.PopulateDB(table)

	// http route handlers
	http.HandleFunc("/health-check", HealthCheck)
	http.HandleFunc("/metrics", MetricIngestion)
	http.HandleFunc("/report", FetchReport)

	// creating http server and listening on port 5000
	clog.Info("go server starting and is listening on : 5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		clog.Panicf("failed to start server - %+v", err)
	}
}

// CheckError helps verifying the error response with appropriate error code
func CheckError(err error, message string, status int)  {
	if err != nil{
		clog.Errorf(fmt.Sprintf(message + " - %+v", err))
		panic(objects.ApiError{
			Status: status,
			Error: statusMessages[status],
		})
		return
	}
}

// RecoverPanic helps in recovering from the panic during api request
func RecoverPanic(writer *http.ResponseWriter)  {
	if err := recover(); err != nil{
		// converting the err type interface to type error
		status := 500
		message := statusMessages[status]
		switch err.(type) {
		case objects.ApiError:
			data := err.(objects.ApiError)
			status = data.Status
			message = data.Error
		case error:
			e := err.(error)
			message = e.Error()
		}
		(*writer).WriteHeader(status)
		(*writer).Write([]byte(message))
		return
	}
}

// ParseIp helps in getting the ip from remote address
func ParseIp(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		clog.Warnf("userip: %q is not IP:port", remoteAddr)
		return ""
	}
	userIP := net.ParseIP(host)
	if userIP == nil {
		clog.Warnf("userip: %q is not IP:port", remoteAddr)
		return ""
	}
	return host
}

// HealthCheck returns a json object to
func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	if method := request.Method; method != http.MethodGet{
		writer.WriteHeader(404)
		writer.Write([]byte("resource not found"))
		return
	}
	resp := objects.GenericResponse{
		Success: true,
		Message: "Happy GO",
		Data: map[string]string{"message": "Go server is working and ready to accept connections"},
	}
	writer.Header().Set("Content-Type", "application/json")
	buf, err := json.Marshal(resp)
	if err != nil{
		clog.Errorf("failed json marshall - %+v", err)
	}
	writer.WriteHeader(200)
	writer.Write(buf)
}

// MetricIngestion captures the metrics pushed and save it data base
func MetricIngestion(writer http.ResponseWriter, request *http.Request)  {
	if method := request.Method; method != http.MethodPost {
		writer.WriteHeader(404)
		writer.Write([]byte("resource not found"))
		return
	}

	// recovering from panic to close the connection with appropriate response
	defer RecoverPanic(&writer)

	// reading the request body from socket
	decoder := json.NewDecoder(request.Body)

	metrics := objects.Metrics{}
	err := decoder.Decode(&metrics)
	if err == io.EOF {
		CheckError(err, "failed to read body params", 400)
	} else{
		CheckError(err, "failed to parse json body", 500)
	}
	// if percentage_cpu_used and percentage_memory_used are empty return with bad request error
	if metrics.CpuUsed == 0 && metrics.MemoryUsed == 0{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(400)
		resp := objects.GenericResponse{
			Success: false,
			Message: statusMessages[400],
			Data: map[string]string{
				"message": "percentage_cpu_used and percentage_memory_used are required",
			},
		}
		resBody, err := json.Marshal(resp)
		CheckError(err, "failed to create response json object", 500)
		writer.Write(resBody)
		return
	}
	// Remote address contains the ip address of the caller client
	metrics.Ip = ParseIp(request.RemoteAddr)
	// inserting the record to db
	_, err = db.Insert("", &metrics)
	CheckError(err, "failed to insert data into db", 500)
	writer.Write([]byte("metric recorded"))
	return
}

// FetchReport returns the max cpu and memory objects recorded
func FetchReport(writer http.ResponseWriter, request *http.Request)  {
	if method := request.Method; method != http.MethodGet {
		writer.WriteHeader(404)
		writer.Write([]byte("resource not found"))
		return
	}
	defer RecoverPanic(&writer)
	queries := request.URL.Query()
	filterBy := "cpu"
	if filter := queries.Get("filter"); filter != ""{
		filterBy = filter
	}
	selector := ""
	switch filterBy {
		case "cpu":
			selector = "select max(percentage_cpu_used), percentage_memory_used, "
		case "memory":
			selector = "select percentage_cpu_used, max(percentage_memory_used), "
		default:
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(400)
			resp := objects.GenericResponse{
				Success: false,
				Message: statusMessages[400],
				Data: map[string]string{
					"message": "please provide a valid filter either cpu or memory",
				},
			}
			resBody, err := json.Marshal(resp)
			CheckError(err, "failed to create response in /report json object", 500)
			writer.Write(resBody)
			return
	}
	query := fmt.Sprintf("%s ip, date from %s group by ip", selector, db.Table)
	rows, err := db.Select("", query)
	CheckError(err, "failed to fetch data from db ", 500)
	var metrics objects.MetricsList
	for rows.Next(){
		metric := objects.Metrics{}
		err := rows.Scan(&metric.CpuUsed, &metric.MemoryUsed, &metric.Ip, &metric.Date)
		if err != nil{
			clog.Errorf("failed to parse row - %+v", err)
			continue
		}
		metrics = append(metrics, metric)
	}
	buf, err := json.Marshal(metrics)
	CheckError(err, "failed to parse db output", 500)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(buf)
	return
}