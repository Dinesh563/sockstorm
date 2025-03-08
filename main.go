package main

import (
	"fmt"
	"math"
	"time"

	"example.com/test-nats/decoder"
	"example.com/test-nats/models"
	"github.com/gorilla/websocket"
)

const url = "ws://localhost:3240/ws/v1/feeds"

func Test1() {
	fmt.Println("Running Test1")

	ch := make(chan *models.ConnectionReport, 1000)

	go func() {

		for {
			msg := <-ch
			fmt.Printf("Report Received from a connection, %+v \n", *msg)
		}

	}()

	ConnectUser(ch)

}

func ConnectUser(ch chan *models.ConnectionReport) {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		fmt.Println("Unable to connect to websocket", err)
		return
	}

	cr := models.ConnectionReport{
		MinLatency: math.MaxInt,
		MaxLatency: math.MinInt,
		NewPackets: &models.ConnectionReport{},
	}

	var reportLastSent = time.Now().Unix()

	// subscribe n number of connections
	Subscribe(conn, 1)

	// send heartbeart periodically
	go Heartbeat(conn)

	// defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Error while Reading message from websocket", err)
			return
		}

		switch v := decoder.DecodeMessage(msg).(type) {
		case decoder.CompactMarketData:
			cr.GenerateConnectionReport(&v)
		default:
			//increment invalid packet
			cr.InvalidPackets++
			cr.TotalPackets++
		}

		if time.Now().Unix()-reportLastSent > 5 {
			reportLastSent = time.Now().Unix()
			//send report to the parent go routine.

			cr.NewPackets.ZeroLatencyPkts = cr.ZeroLatencyPkts - cr.NewPackets.ZeroLatencyPkts
			cr.NewPackets.OneSecondLatencyPkts = cr.OneSecondLatencyPkts - cr.NewPackets.OneSecondLatencyPkts
			cr.NewPackets.InvalidPackets = cr.InvalidPackets - cr.NewPackets.InvalidPackets
			cr.NewPackets.TotalPackets = cr.TotalPackets - cr.NewPackets.TotalPackets
			cr.NewPackets.ZeroLatencyPkts = cr.ZeroLatencyPkts - cr.NewPackets.ZeroLatencyPkts

			ch <- &cr
		}

	}
}

func Subscribe(conn *websocket.Conn, noInstruments int) {
	msg := map[string]interface{}{"a": "subscribe", "v": [][]int{{1, 26024}}, "m": "compact_marketdata"}
	conn.WriteJSON(msg)
}

func Heartbeat(conn *websocket.Conn) {
	// periodically send hearbeat msg every 30 seconds
	for {
		msg := map[string]interface{}{"a": "h", "v": []int{}, "m": ""}
		conn.WriteJSON(msg)
		time.Sleep(25 * time.Second)
	}
}

func main() {
	fmt.Println("Testing NATS with GO..")
	Test1()
}
