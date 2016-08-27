package dbc

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

type A struct {
	B int `db:"b"`
	C int
	D int `json:"d"`
	E
}

type E struct {
	I int64 `db:"i" sql:"ignore"`
	F
}

type F struct {
	F int `db:"e"`
}

type G struct {
	F *F
}

func TestCollectColumn(t *testing.T) {
	c := strings.Join(CollectColumn(&A{}), "")
	if c != "be" {
		t.Errorf("Expected \"b\", but it was %v instead.", c)
	}

	c = strings.Join(CollectColumn(nil), "")
	if c != "" {
		t.Errorf("Expected [], but it was %v instead.", c)
	}

	g := &G{}
	g.F = &F{F: 1}
	c = strings.Join(CollectColumn(g), "")
	if c != "e" {
		t.Errorf("Expected [], but it was %v instead.", c)
	}
}

func TestTimeStamp(t *testing.T) {
	ts := TimeStamp{Time: time.Unix(1467362677, 0)}
	b, err := json.Marshal(ts)
	if s := string(b); err != nil || s != "1467362677" {
		t.Errorf("Expected 1467362677, but it was %v instead.", s)
	}

	var ts0 TimeStamp
	json.Unmarshal([]byte("1467362677"), &ts0)
	if ts := ts0.Unix(); ts != 1467362677 {
		t.Errorf("Expected 1467362677, but it was %v instead.", ts)
	}
}
