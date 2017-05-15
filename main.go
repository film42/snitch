package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var checkHostService *CheckHostService
var checkProcessService *CheckProcessService

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if checkHostService == nil || checkProcessService == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	successRate := checkHostService.SuccessRate()
	fmt.Println("HEALTH CHECK - Current host success rate:", successRate)
	if successRate <= 0.7 {
		fmt.Println("Error rate:", successRate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	successRate = checkProcessService.SuccessRate()
	fmt.Println("HEALTH CHECK - Current process success rate:", successRate)
	if successRate <= 0.7 {
		fmt.Println("Error rate:", successRate)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type StringList []string

func (this *StringList) String() string {
	return strings.Join((*this)[:], ", ")
}

func (this *StringList) Set(value string) error {
	*this = append(*this, value)
	return nil
}

func main() {
	var checkHosts StringList
	var processes StringList

	flag.Var(&checkHosts, "check", "List of host:port combos to check. Example: --check localhost:1234 --check localhost:3422")
	flag.Var(&processes, "process", "List of process substring to check. Example: --process sidekiq --process puma")
	portPtr := flag.Int("port", 9999, "Port to listen on")
	flag.Parse()

	if len(checkHosts) == 0 && len(processes) == 0 {
		fmt.Println("Must provide at least one process or host to check:")
		flag.PrintDefaults()
		return
	}

	checkHostService = NewCheckHostService(checkHosts)
	go checkHostService.Start()

	checkProcessService = NewCheckProcessService(processes)
	go checkProcessService.Start()

	http.HandleFunc("/", healthCheckHandler)
	http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil)
}
