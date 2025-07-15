package decoder

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"example.com/test-nats/utils"
)

// MarketData represents the decoded market data
type CompactMarketData struct {
	Mode                 uint8
	Exchange             int8
	InstrumentToken      int32
	LTP                  int32
	Change               int32
	LTT                  int32
	LTTIST               time.Time
	LowDPR               int32
	HighDPR              int32
	CurrentOpenInterest  int32
	InitialOpenInterest  int32
	BidPrice             int32
	AskPrice             int32
	DiffWithCurrrentTime int
}

// DecodeMessage decodes binary WebSocket messages based on mode
func DecodeMessage(data []byte) interface{} {
	if len(data) < 1 {
		return "Invalid data"
	}

	// Read the mode (first byte)
	mode := data[0]
	
	reader := bytes.NewReader(data)

	switch mode {
	case 254:
		return DecodeMarketInfo(mode, reader)
	case 58:
		return DecodePositionUpdate(mode, reader)
	case 50:
		return DecodeOrderUpdate(mode, reader)
	case 2:
		return DecodeCompactMarketData(data, reader)
	case 4:
		return DecodeFullSnapQuote(mode, reader)
	case 14:
		return DecodeTBTSnapQuote(mode, reader)
	default:
		return fmt.Sprintf("Unhandled mode %d", mode)
	}
}

// decodeMarketData decodes binary market data message
func DecodeCompactMarketData(data []byte, reader *bytes.Reader) CompactMarketData {
	if len(data) < 42 {
		return CompactMarketData{}
	}

	// Initialize MarketData struct
	var decoded CompactMarketData

	// Read values from the binary buffer
	binary.Read(reader, binary.BigEndian, &decoded.Mode)
	binary.Read(reader, binary.BigEndian, &decoded.Exchange)
	binary.Read(reader, binary.BigEndian, &decoded.InstrumentToken)
	binary.Read(reader, binary.BigEndian, &decoded.LTP)
	binary.Read(reader, binary.BigEndian, &decoded.Change)
	binary.Read(reader, binary.BigEndian, &decoded.LTT)
	binary.Read(reader, binary.BigEndian, &decoded.LowDPR)
	binary.Read(reader, binary.BigEndian, &decoded.HighDPR)
	binary.Read(reader, binary.BigEndian, &decoded.CurrentOpenInterest)
	binary.Read(reader, binary.BigEndian, &decoded.InitialOpenInterest)
	binary.Read(reader, binary.BigEndian, &decoded.BidPrice)
	binary.Read(reader, binary.BigEndian, &decoded.AskPrice)

	// Convert UNIX timestamp to human-readable IST time
	decoded.LTTIST = utils.FormatTimestamp(decoded.LTT)
	decoded.DiffWithCurrrentTime = utils.TimeDifference(time.Now(), decoded.LTTIST)
	return decoded
}

// Dummy functions for other decoders
func DecodeMarketInfo(mode uint8, reader *bytes.Reader) string {
	return "Market Info Decoding Not Implemented"
}

func DecodePositionUpdate(mode uint8, reader *bytes.Reader) string {
	return "Position Update Decoding Not Implemented"
}

func DecodeOrderUpdate(mode uint8, reader *bytes.Reader) string {
	return "Order Update Decoding Not Implemented"
}

func DecodeFullSnapQuote(mode uint8, reader *bytes.Reader) string {
	return "Full Snap Quote Decoding Not Implemented"
}

func DecodeTBTSnapQuote(mode uint8, reader *bytes.Reader) string {
	return "TBT Snap Quote Decoding Not Implemented"
}
