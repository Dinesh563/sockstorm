package models

import "example.com/test-nats/decoder"

type ConnectionReport struct {
	TotalPackets         int
	InvalidPackets       int
	TotalSubscriptions   int
	ZeroLatencyPkts      int
	OneSecondLatencyPkts int
	LatencyPackets       int
	MinLatency           int
	MaxLatency           int
	AvgLatency           int
	NewPackets           *ConnectionReport
}

func (cr *ConnectionReport) GenerateConnectionReport(packet *decoder.CompactMarketData) {
	switch packet.DiffWithCurrrentTime {
	case 0:
		cr.ZeroLatencyPkts++
	case 1:
		cr.OneSecondLatencyPkts++
	case 43200:
		// invalid timestamp packet
		cr.InvalidPackets++
	default:
		// high latency packet.
		AppendLatencyPacket(packet, cr)

	}
	cr.TotalPackets++
}

func AppendLatencyPacket(packet *decoder.CompactMarketData, cr *ConnectionReport) {
	cr.AvgLatency = (cr.AvgLatency*cr.LatencyPackets + packet.DiffWithCurrrentTime) / (cr.LatencyPackets + 1)
	cr.LatencyPackets++

	if packet.DiffWithCurrrentTime > cr.MaxLatency {
		cr.MaxLatency = packet.DiffWithCurrrentTime
	}

	if packet.DiffWithCurrrentTime < cr.MinLatency {
		cr.MinLatency = packet.DiffWithCurrrentTime
	}
}
