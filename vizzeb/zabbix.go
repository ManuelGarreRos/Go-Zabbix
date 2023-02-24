package vizzeb

import (
	"fmt"
	"github.com/essentialkaos/go-zabbix"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	address = "zabbix.visualtis.com"
	ip      = "87.106.52.15:10051"
	//Hostname of the client
	//Your pc if u'r testing or the server in production
	//This hostname and the IP should be listed at the Zabbix client
	hostname = "visualtis-ThinkPad-P1"
)

func StartZabbix() {
	client, err := zabbix.NewClient(ip, hostname)
	if err != nil {
		fmt.Printf("zabbix.NewClient error: %v \n: ", err)
		return
	}
	fmt.Printf("conection successfully to: %v \n", address)
	for {
		cpuUsage, mAlloc, mFrees, mPointers, mSystem := getMetrics()
		client.Add("go.cpu", cpuUsage)
		client.Add("go.status", 1)
		client.Add("go.memoryallocation", mAlloc)
		client.Add("go.memoryfress", mFrees)
		client.Add("go.memorypointers", mPointers)
		client.Add("go.memorysystem", mSystem)
		client.Add("go.intprueba", 1)
		//Test only
		res, err := client.Send()
		//Production
		//_, err = client.Send()
		if err != nil {
			fmt.Printf("Send error: %v \n", err)
		}
		//Comment metrics sended in production
		fmt.Printf("Metrics sended (processed: %d | failed: %d | total: %d) \n", res.Processed, res.Failed, res.Total)
		ZabbixErrorLog("log test")
		ZabbixPanicLog("Panic log test")
		dt := time.Now()
		fmt.Println("Current date and time is: ", dt.String())
		time.Sleep(300 * time.Second)
	}
}

func ZabbixErrorLog(log string) {
	client, err := zabbix.NewClient(ip, hostname)
	if err != nil {
		fmt.Printf("zabbix.NewClient error: %v \n: ", err)
		return
	}
	client.Add("go.log", log)
	_, err = client.Send()
	if err != nil {
		fmt.Printf("Send error: %v \n", err)
	}
}

func ZabbixPanicLog(panic string) {
	client, err := zabbix.NewClient(ip, hostname)
	if err != nil {
		fmt.Printf("zabbix.NewClient error: %v \n: ", err)
		return
	}
	client.Add("go.panic", panic)
	_, err = client.Send()
	if err != nil {
		fmt.Printf("Send error: %v \n", err)
	}
}

func getMetrics() (float64, uint64, uint64, uint64, uint64) {
	cmd := exec.Command("bash", "-c", "top -b -n 1 | awk '/^%Cpu/{print $2}'")
	output, err := cmd.Output()
	cpuUsage := 0.0
	if err != nil {
		fmt.Println(err)
	} else {
		cpuUsageString := strings.TrimSpace(string(output))
		cpuUsageString = strings.ReplaceAll(cpuUsageString, ",", ".")
		cpuUsage, err = strconv.ParseFloat(cpuUsageString, 64)
		if err != nil {
			fmt.Printf("Parse to float err: %s \n", err)
		}
	}
	var mStats runtime.MemStats
	runtime.ReadMemStats(&mStats)
	//comment metrics readings on production
	fmt.Printf("CPU usage: %v \n", cpuUsage)
	fmt.Printf("Memory allocation: %v \n", mStats.Alloc)
	fmt.Printf("Memory frees: %v \n", mStats.Frees)
	fmt.Printf("Memory threads: %v \n", mStats.Lookups)
	fmt.Printf("Memory system: %v \n", mStats.Sys)

	return cpuUsage, mStats.Alloc, mStats.Frees, mStats.Lookups, mStats.Sys
}
