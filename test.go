package main

import (
	"fmt"
	"math"
	"time"

	"example.com/test-nats/models"
)

func TestExponentially() {
	fmt.Println("Running Test1")

	tr := models.TestReport{Name: "My test", TotalConnections: N, PktAccumulator: make(map[int]*models.ConnectionReport, N)}

	go tr.Print()

	ch := make(chan *models.ConnectionReport, 1000)

	go func() {
		for {
			tr.ConsumeConnectionReportForTestReport(<-ch)
		}
	}()

	x := 0.0
	totalConnectionsSpawned := 0
	i := 0
	for totalConnectionsSpawned < N {
		rate := int(math.Floor(math.Exp(x)))

		for k := i ; k < (i + rate) && k < N; k++ {
			go ConnectUser(ch, k)
		}
		i += rate
		x += 0.2 // e^x = 1 when x = 0
		totalConnectionsSpawned += rate
	}

	time.Sleep(10 * time.Minute)
}
