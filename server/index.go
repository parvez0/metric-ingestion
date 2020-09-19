package server

import (
	"encoding/json"
	"github.com/parvez0/whatsapp-provisioner/custom_logger"
	"net/http"
)

var clog = custom_logger.NewLogger()

// CreateServer creates a http service with default handlers
func CreateServer()  {
	http.HandleFunc("/health-check", HealthCheck)
	if err := http.ListenAndServe(":5000", nil); err != nil {
		clog.Panicf("failed to start server - %+v", err)
	}
	clog.Info("go server started and is listening on : 5000")
}

// HealthCheck returns a json object to
func HealthCheck(writer http.ResponseWriter, request *http.Request) {
	if method := request.Method; method != http.MethodGet{
		writer.WriteHeader(404)
		writer.Write([]byte("resource not found"))
		return
	}
	resp := GenericResponse{
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