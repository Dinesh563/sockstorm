package subjects

import (
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type subjects struct {
	NseSubjects   [][2]int
	BseSubjects   [][2]int
	TotalSubjects [][2]int
}

var Subjects = subjects{
	NseSubjects: NseSubjects,
	BseSubjects: BseSubjects,
}

func init() {
	Subjects.TotalSubjects = append(Subjects.BseSubjects, Subjects.NseSubjects...)
}

func (s *subjects) Subscribe(conn *websocket.Conn, n int) {

	if n > len(s.TotalSubjects) {
		n = len(s.TotalSubjects)
	}

	// Shuffle the slice
	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(n, func(i, j int) {
		s.TotalSubjects[i], s.TotalSubjects[j] = s.TotalSubjects[j], s.TotalSubjects[i]
	})

	s.TotalSubjects[0] = [2]int{1, 26004}
	// subscribe
	msg := map[string]interface{}{"a": "subscribe", "v": s.TotalSubjects[:n], "m": "compact_marketdata"}
	conn.WriteJSON(msg)

	// msg := map[string]interface{}{"a": "subscribe", "v": [][2]int{{1, 26024}}, "m": "compact_marketdata"}

	// conn.WriteJSON(msg)
}
