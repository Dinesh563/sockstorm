package models

import (
	// "fmt

	"example.com/test-nats/decoder"
)

type ConnectionReport struct {
	Id                   int
	TotalPackets         int
	InvalidPackets       int
	TotalSubscriptions   int
	ZeroLatencyPkts      int
	OneSecondLatencyPkts int
	LatencyPackets       int
	MinLatency           int
	MaxLatency           int
	AvgLatency           int
	Alive                bool
	NewPackets           *ConnectionReport
}

func (cr *ConnectionReport) GenerateConnectionReport(packet *decoder.CompactMarketData) {
	switch v := packet.DiffWithCurrrentTime; {
	case v == 0:
		cr.ZeroLatencyPkts++
	case v == 1:
		cr.OneSecondLatencyPkts++
	case v == 43200, v >= 18_000:
		// invalid timestamp packet
		cr.InvalidPackets++
	default:
		// high latency packet.
		AppendLatencyPacket(packet, cr)
	}
	cr.TotalPackets++
}

func AppendLatencyPacket(packet *decoder.CompactMarketData, cr *ConnectionReport) {
	// fmt.Printf("%+v \n", packet)
	cr.AvgLatency = (cr.AvgLatency*cr.LatencyPackets + packet.DiffWithCurrrentTime) / (cr.LatencyPackets + 1)
	cr.LatencyPackets++

	if packet.DiffWithCurrrentTime > cr.MaxLatency {
		cr.MaxLatency = packet.DiffWithCurrrentTime
	}

	if packet.DiffWithCurrrentTime < cr.MinLatency {
		cr.MinLatency = packet.DiffWithCurrrentTime
	}
}
