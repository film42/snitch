package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type ProcessCheck struct {
	Name      string
	IsPresent bool
}

func NewProcessCheck(name string) *ProcessCheck {
	return &ProcessCheck{
		Name:      name,
		IsPresent: true, // defaults to true
	}
}

func pgrepContains(name string) bool {
	outputBuffer, err := exec.Command("pgrep", "-f", name).Output()
	if err != nil {
		fmt.Println("Could not execute pgrep for name", name, "-", err)
		return false
	}

	currentPidStr := strconv.Itoa(os.Getpid())
	pids := []string{}
	scanner := bufio.NewScanner(bytes.NewReader(outputBuffer))
	for scanner.Scan() {
		pid := scanner.Text()
		if pid == currentPidStr {
			// Don't include snitch's pid.
			continue
		}

		pids = append(pids, pid)
	}

	if len(pids) == 0 {
		return false
	}

	return true
}

func (this *ProcessCheck) PerformHealthCheck() {
	if pgrepContains(this.Name) {
		this.IsPresent = true
		return
	}

	this.IsPresent = false
}

type CheckProcessService struct {
	ProcessChecks []*ProcessCheck
}

func NewCheckProcessService(processes []string) *CheckProcessService {
	checks := []*ProcessCheck{}
	for _, process := range processes {
		checks = append(checks, NewProcessCheck(process))
	}
	return &CheckProcessService{
		ProcessChecks: checks,
	}
}

func (this *CheckProcessService) SuccessRate() float64 {
	if len(this.ProcessChecks) == 0 {
		return 1.0
	}

	successCount := 0.0
	for _, check := range this.ProcessChecks {
		if check.IsPresent {
			successCount += 1.0
		}
	}

	return successCount / float64(len(this.ProcessChecks))
}

func (this *CheckProcessService) Start() {
	// This is not a great check sevice, but it does the job.
	for {
		for _, check := range this.ProcessChecks {
			check.PerformHealthCheck()
		}

		time.Sleep(time.Second * 10)
	}
}
