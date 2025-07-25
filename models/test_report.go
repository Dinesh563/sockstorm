package models

import (
	"fmt"
	"math"
	"time"

	"example.com/test-nats/utils"
)


const TEST_PRINT_TIME = 5 * time.Second // prints every 5 seconds.

type TestReport struct {
	Name                 string
	TotalConnections     int
	ActiveConnections    int
	TotalSubscriptions   int
	ActiveSubscriptions  int
	AvgSubscriptions     int
	TotalPackets         int
	InvalidPackets       int
	ZeroLatencyPkts      int
	OneSecondLatencyPkts int
	LatencyPackets       int
	MinLatency           int
	MaxLatency           int
	AvgLatency           float64
	PktAccumulator       map[int](*ConnectionReport)
}

func (tr *TestReport) ConsumeConnectionReportForTestReport(cr *ConnectionReport) {
	// tr.TotalPackets = cr.NewPackets.TotalPackets
	// tr.InvalidPackets = cr.NewPackets.InvalidPackets
	// tr.ZeroLatencyPkts = cr.NewPackets.ZeroLatencyPkts
	tr.PktAccumulator[cr.Id] = cr
	// fmt.Printf("%+v\n", cr)
}

func (tr *TestReport) Print() {
	ticker := time.NewTicker(TEST_PRINT_TIME)
	defer ticker.Stop()
	for range ticker.C {
		tr.ConstructTestReport()
		PrettyPrintReport(tr)
	}
}

func (tr *TestReport) ConstructTestReport() {
	var (
		minLatency           = math.MaxInt
		maxLatency           = math.MinInt
		totalLatency         = 0
		totalPackets         = 0
		invalidPackets       = 0
		zeroLatencyPkts      = 0
		oneSecondLatencyPkts = 0
		activeConnections    = 0
		latencyPackets       = 0
	)

	for _, cr := range tr.PktAccumulator {
		if cr.Alive {
			activeConnections++
		}
		minLatency = utils.Min(minLatency, cr.MinLatency)
		maxLatency = utils.Max(maxLatency, cr.MaxLatency)
		totalLatency += cr.AvgLatency
		totalPackets += cr.TotalPackets
		invalidPackets += cr.InvalidPackets
		zeroLatencyPkts += cr.ZeroLatencyPkts
		oneSecondLatencyPkts += cr.OneSecondLatencyPkts
		latencyPackets += cr.LatencyPackets
	}
	tr.ActiveConnections = activeConnections
	if latencyPackets > 0 {
		tr.AvgLatency = float64(totalLatency) / float64(latencyPackets)
	} else {
		tr.AvgLatency = 0.0
	}
	tr.LatencyPackets = latencyPackets

	tr.MinLatency = minLatency
	tr.MaxLatency = maxLatency
	tr.TotalPackets = totalPackets
	tr.InvalidPackets = invalidPackets
	tr.ZeroLatencyPkts = zeroLatencyPkts
	tr.OneSecondLatencyPkts = oneSecondLatencyPkts
}

func PrettyPrintReport(tr *TestReport) {
	fmt.Printf("===========================================================================\n")
	fmt.Println("📊 Test Report Summary")
	fmt.Println("----------------------")
	fmt.Printf("Name:                  %s\n", tr.Name)
	fmt.Printf("Total Connections:     %d\n", tr.TotalConnections)
	fmt.Printf("Active Connections:    %d\n", tr.ActiveConnections)
	fmt.Printf("Total Packets:         %d\n", tr.TotalPackets)
	fmt.Printf("Invalid Packets:       %d\n", tr.InvalidPackets)
	fmt.Printf("Zero Latency Packets:  %d\n", tr.ZeroLatencyPkts)
	fmt.Printf("1s Latency Packets:    %d\n", tr.OneSecondLatencyPkts)
	fmt.Printf("Latency Packets:       %d\n", tr.LatencyPackets)

	fmt.Println("\n📈 Latency Stats")
	fmt.Println("----------------------")
	if tr.MinLatency == math.MaxInt64 {
		fmt.Println("Min Latency:           (no data)")
	} else {
		fmt.Printf("Min Latency:           %d s\n", tr.MinLatency)
	}

	if tr.MaxLatency == math.MinInt64 {
		fmt.Println("Max Latency:           (no data)")
	} else {
		fmt.Printf("Max Latency:           %d s\n", tr.MaxLatency)
	}

	fmt.Printf("Avg Latency:           %f s\n", tr.AvgLatency)
	fmt.Printf("===========================================================================\n")
}
