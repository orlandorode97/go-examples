package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Response struct {
	Jobs []*Job `json:"jobs"`
}
type Job struct {
	ID string `json:"id"`
}

func main() {
	jsonBody := `{
	"jobs": [
		{
			"links": {
				"job_spec": "orders_report",
				"schedule": null
			},
			"options": {
				"dealer_id": "2",
				"sale_code": "JQSVN"
			}
		}
	]
    }`

	request, err := http.NewRequest(http.MethodPost, "http://localhost/job/jobs/", strings.NewReader(jsonBody))
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer here")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	jsonResponse, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	var jobResult Response
	err = json.Unmarshal(jsonResponse, &jobResult)
	if err != nil {
		panic(err)
	}
	fmt.Println("Job created successfully", response.StatusCode, jobResult.Jobs[0].ID)

	time.Sleep(30 * time.Millisecond)

	request, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost/job/jobs/%v", jobResult.Jobs[0].ID), nil)
	if err != nil {
		panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer here")

	response, err = http.DefaultClient.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	fmt.Println("Job cancelled successfully", response.StatusCode)

}
