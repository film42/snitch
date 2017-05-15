package main

import (
	"fmt"
	"net/http"
	"time"
)

type HostCheck struct {
	HostPort  string
	IsHealthy bool
}

func NewHostCheck(hostPort string) *HostCheck {
	return &HostCheck{
		HostPort: hostPort,
		// default to healthy
		IsHealthy: true,
	}
}

func (this *HostCheck) UrlString() string {
	return "http://" + this.HostPort
}

func (this *HostCheck) PerformHealthCheck() {
	httpClient := &http.Client{
		Timeout: time.Second * 5,
	}

	urlString := this.UrlString()
	response, err := httpClient.Get(urlString)
	if err != nil {
		this.IsHealthy = false
		fmt.Println("Error checking", urlString, "got error", err)
		return
	}
	response.Body.Close()

	if response.StatusCode >= 300 {
		this.IsHealthy = false
		fmt.Println("ALERT!", urlString, "returned status", response.StatusCode)
		return
	}

	this.IsHealthy = true
}

type CheckHostService struct {
	HostChecks []*HostCheck
}

func NewCheckHostService(checkHosts []string) *CheckHostService {
	checks := []*HostCheck{}
	for _, checkHost := range checkHosts {
		checks = append(checks, NewHostCheck(checkHost))
	}

	return &CheckHostService{
		HostChecks: checks,
	}
}

func (this *CheckHostService) SuccessRate() float64 {
	if len(this.HostChecks) == 0 {
		return 1.0
	}

	successCount := 0.0
	for _, check := range this.HostChecks {
		if check.IsHealthy {
			successCount += 1.0
		}
	}

	return successCount / float64(len(this.HostChecks))
}

func (this *CheckHostService) Start() {
	// This is not a great check sevice, but it does the job.
	for {
		for _, check := range this.HostChecks {
			check.PerformHealthCheck()
		}

		time.Sleep(time.Second * 10)
	}
}
