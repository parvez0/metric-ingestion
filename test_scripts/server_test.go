package test_scripts

import (
	"github.com/parvez0/go-requests/requests"
	"github.com/parvez0/whatsapp-provisioner/server"
	"testing"
)

// TestServer creates a server and makes a get request to health check
func TestServer(t *testing.T) {
	go server.CreateServer()
	client := requests.NewClient()
	options := requests.Options{
		Url:     "http://localhost:5000/health-check",
		Method:  "GET",
	}
	client.NewRequest(options)
	res, err := client.Send()
	if err != nil{
		t.Fatalf("failed make a get request to server - %+v", err)
	}
	t.Logf("server is wokring statusCode - %d", res.GetStatusCode())
}
