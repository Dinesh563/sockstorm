package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"example.com/test-nats/decoder"
	"example.com/test-nats/models"
	"example.com/test-nats/subjects"
	"github.com/gorilla/websocket"
)

// const url = "ws://localhost:3240/ws/v1/feeds"
const (
	url                            = "wss://uat1.tradelab.ltd/ws/v1/feeds"
	N                              = 1000 // N concurrent connections
	CONNECTION_REPORT_TIME         = 5 * time.Second
	HeartbeatTime                  = 25 * time.Second
	NoOfSubscriptionsPerConnection = 100
)

func Test1() {
	fmt.Println("Running Test1")

	tr := models.TestReport{Name: "My test", TotalConnections: N, PktAccumulator: make(map[int]*models.ConnectionReport, N)}

	go tr.Print()

	ch := make(chan *models.ConnectionReport, 1000)

	go func() {

		for {
			tr.ConsumeConnectionReportForTestReport(<-ch)
			// TODO : process each connection report and create a test report
			// fmt.Printf("Report Received from a connection, %+v \n", *msg)
		}

	}()

	for i := 0; i < N; i++ {
		go ConnectUser(ch, i)
	}
	time.Sleep(10 * time.Minute)
}

// func isWebSocketAlive(conn *websocket.Conn) bool {
// 	if conn == nil {
// 		return false
// 	}

// 	err := conn.WriteMessage(websocket.PingMessage, nil)
// 	return err == nil
// }



func ConnectUser(ch chan *models.ConnectionReport, id int) {

	cr := models.ConnectionReport{
		Id:         id,
		MinLatency: math.MaxInt,
		MaxLatency: math.MinInt,
		Alive:      true,
		NewPackets: &models.ConnectionReport{},
	}

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		fmt.Println("Unable to connect to websocket", err)
		cr.Alive = false
		ch <- &cr
		return
	}

	// subscribe n number of suscriptions randomly
	subjects.Subjects.Subscribe(conn, NoOfSubscriptionsPerConnection)

	// send heartbeart periodically
	go Heartbeat(conn)

	ReadFirstN_SecondsPkts(5, conn)
	// defer conn.Close()

	go func() {
		ticker := time.NewTicker(CONNECTION_REPORT_TIME)
		for range ticker.C {
			ch <- &cr
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("error while Reading message from websocket", err)
			cr.Alive = false
			ch <- &cr
			return
		}

		switch v := decoder.DecodeMessage(msg).(type) {
		// type check
		case decoder.CompactMarketData:
			cr.GenerateConnectionReport(&v)
		default:
			// increment invalid packet
			cr.InvalidPackets++
			cr.TotalPackets++
		}

	}

}

func Heartbeat(conn *websocket.Conn) {
	// periodically send hearbeat msg every 30 seconds
	for {
		msg := map[string]interface{}{"a": "h", "v": []int{}, "m": ""}
		if conn.WriteJSON(msg) != nil {
			break
		}
		time.Sleep(HeartbeatTime)
	}
}

// read_packets after n seconds
func ReadFirstN_SecondsPkts(n int, conn *websocket.Conn) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(n)*time.Second)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Read first %v second packets \n", n)
			return
		default:
			conn.ReadMessage()
		}
	}
}

func main() {
	fmt.Println("Testing NATS with GO..")
	// Test1()
	TestExponentially()
}
