package models

type TestReport struct {
	Name                        string
	TotalConnections            int
	ActiveConnections           int
	TotalSubscriptions          int
	ActiveSubscriptions         int
	AvgSubscriptions            int
	TotalPackets                int
	InvalidPackets              int
	ZeroLatencyPkts             int
	OneSecondLatencyPkts        int
	LatencyPackets              int
	AvgLatencyPktsPerConnection int
	MinLatency                  int
	MaxLatency                  int
	AvgLatency                  int
}

func (tr *TestReport) GenerateConnectionReport(cr *ConnectionReport) {
    tr.TotalPackets = cr.NewPackets.TotalPackets
	tr.InvalidPackets = cr.NewPackets.InvalidPackets
	tr.ZeroLatencyPkts = cr.NewPackets.ZeroLatencyPkts
}
