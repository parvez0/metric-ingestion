package test_scripts

import (
	"github.com/parvez0/go-requests/requests"
	"github.com/parvez0/metric-ingestion/objects"
	"net/http"
	"strconv"
	"testing"
)

// initializing http client with default settings and base path
var client = requests.NewClient(requests.GlobalOptions{
	BasePath: "http://localhost:5000",
})

var dummyData = objects.MetricsList{
	{
		CpuUsed: 30,
		MemoryUsed: 40,
	},
	{
		CpuUsed: 60,
		MemoryUsed: 70,
	},
	{
		CpuUsed: 30,
		MemoryUsed: 70,
	},
	{
		CpuUsed: 55,
		MemoryUsed: 80,
	},
	{
		CpuUsed: 10,
		MemoryUsed: 90,
	},
}

// TestServer creates a server and makes a get request to health check
func TestServer(t *testing.T) {
	options := requests.Options{
		Url:     "/health-check",
		Method:  "GET",
	}
	client.NewRequest(options)
	res, err := client.Send()
	if err != nil{
		t.Fatalf("failed make a get request to server - %+v", err)
	}
	t.Logf("server is wokring statusCode - %d", res.GetStatusCode())
}

// TestMetricsPush pushes dummy metrics to the server for testing
func TestMetricsPush(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	options := requests.Options{
		Url:     "/metrics",
		Method:  http.MethodPost,
		Headers: headers,
	}
	for i, metric := range dummyData {
		t.Run("pushing metric - " + strconv.Itoa(i), func(t *testing.T) {
			options.Body = metric
			client.NewRequest(options)
			res, err := client.Send()
			if err != nil{
				t.Fatalf("failed make a get request to server - %+v", err)
			}
			want := http.StatusOK
			got := res.GetStatusCode()
			if want != got{
				t.Fatalf("failed to push metrics with status code - %d", got)
			}
		})
	}
}

// TestVerifyReport tests route verify
func TestVerifyReport(t *testing.T) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	options := requests.Options{
		Url:     "/report",
		Method:  http.MethodGet,
		Headers: headers,
	}
	client.NewRequest(options)
	res, err := client.Send()
	if err != nil{
		t.Fatalf("get /report failed with error - %+v", err)
	}
	want := http.StatusOK
	got := res.GetStatusCode()
	if want != got {
		t.Fatalf("test failed wanted : %d, got : %d", want, got)
	}
}

