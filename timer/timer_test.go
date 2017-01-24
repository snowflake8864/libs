package timer

import (
	P "github.com/snowflake8864/libs/public"
	"log"
	"testing"
	"time"
)

type Demo struct {
	keepalive_tmr uint32
	Data          string
}

func tcp_keepalive(now uint64, arg P.Void) {
	t, ok := arg.(*Demo)

	if ok {
		log.Println("Value=", t.Data)
	} else {
		log.Println("Value=", arg)

	}
	t.keepalive_tmr = TimerAdd(1, tcp_keepalive, t)
}

func TestTimer(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	InitTimer()
	t1 := new(Demo)
	t1.Data = "hello Wyy, where are you now!"
	t1.keepalive_tmr = TimerAdd(1, tcp_keepalive, t1)
	for {
		time.Sleep(5 * time.Second)
	}
}
