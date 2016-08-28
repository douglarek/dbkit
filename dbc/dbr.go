package dbc

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"time"

	"github.com/gocraft/dbr"
)

var _ dbr.EventReceiver = (*simpleEventReceiver)(nil)

// SimpleEventReceiver implements dbr.EventReceiver interface.
type simpleEventReceiver struct {
	out io.Writer
}

// Event implements dbr.EventReceiver.Event.
func (r *simpleEventReceiver) Event(eventName string) {
	fmt.Fprintf(r.out, "[%v]\n", eventName)
}

// EventKv implements dbr.EventReceiver.EventKv.
func (r *simpleEventReceiver) EventKv(eventName string, kvs map[string]string) {}

// EventErr implements dbr.EventReceiver.EventErr.
func (r *simpleEventReceiver) EventErr(eventName string, err error) error {
	fmt.Fprintf(r.out, "[%v] %v\n", eventName, err)
	return err
}

// EventErrKv implements dbr.EventReceiver.EventErrKv.
func (r *simpleEventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	fmt.Fprintf(r.out, "[%v] |%s %v %s| %v\n", eventName, color("red"), err, color("reset"), kvs)
	return err
}

// Timing implements dbr.EventReceiver.Timing.
func (r *simpleEventReceiver) Timing(eventName string, nanoseconds int64) {}

// TimingKv implements dbr.EventReceiver.TimingKv.
func (r *simpleEventReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	ns := time.Duration(nanoseconds)
	ct, rs := color("reset"), color("reset")
	if nanoseconds/int64(time.Millisecond) > 10 {
		ct = color("red")
	}
	fmt.Fprintf(r.out, "[%v] |%s %v %s| %v\n", eventName, ct, ns, rs, kvs)
}

// NewEventReceiver returns a new dbr EventReceiver.
func NewEventReceiver(out ...io.Writer) dbr.EventReceiver {
	if len(out) > 0 {
		return &simpleEventReceiver{out: out[0]}
	}
	return &simpleEventReceiver{out: os.Stdout}
}

// New returns a new dbr session.
func New(c *dbr.Connection, er ...dbr.EventReceiver) *dbr.Session {
	er0 := (dbr.EventReceiver)(nil)
	if len(er) > 0 {
		er0 = er[0]
	}
	return c.NewSession(er0)
}

// CollectColumn returns auto mapping columns with table in dbr insert.
func CollectColumn(table interface{}) (columns []string) {
	v := reflect.ValueOf(table)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		columns = append(columns, CollectColumn(v.Elem().Interface())...)
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			switch f.Type.Kind() {
			case reflect.Ptr, reflect.Struct:
				columns = append(columns, CollectColumn(v.Field(i).Interface())...)
			}
			if f.PkgPath != "" && !f.Anonymous {
				continue
			}
			t := f.Tag.Get("db")
			if t == "-" || t == "" || f.Tag.Get("sql") == "ignore" {
				continue
			}
			columns = append(columns, t)
		}
	}
	return columns
}

// color for terminal output.
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

func color(c string) string {
	switch c {
	case "green":
		return green
	case "white":
		return white
	case "yellow":
		return yellow
	case "red":
		return red
	case "blue":
		return blue
	case "magenta":
		return magenta
	case "cyan":
		return cyan
	case "reset":
		return reset
	}
	panic("unknown color!")
}
