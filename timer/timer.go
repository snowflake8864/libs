package timer

import (
	//	"constant"
	"github.com/snowflake8864/libs/heap"
	P "github.com/snowflake8864/libs/public"
	"log"
	"time"
)

const MAX_TIMERS = uint32(30)

var timerCount = 0

var timers = &heap.Heap{}

func loopChekTimer() {
	for {
		select {
		//case <-time.AfterFunc(1*time.Second, timer_func):
		case <-time.After(1 * time.Second):
			check_timers()
		}
	}
}

func compare(item1, item2 heap.EleType) int {

	k1, _ := item1.(*ZBtimerRef)
	k2, _ := item2.(*ZBtimerRef)

	//log.Printf("%d-----%d\n", k1.id, k2.id)
	if k1.id < k2.id {
		log.Println("return 1")
		return 1
	} else if k1.id > k2.id {
		log.Println("return -1")
		return -1
	} else {
		log.Println("return 0")
		return 0
	}
}

func InitTimer() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	timers.Init(compare, nil, 100)
	go loopChekTimer()
}

type Func func(expire uint64, data P.Void)

type ZBtimer struct {
	arg      P.Void
	function Func
}

type ZBtimerRef struct {
	expire uint64
	id     uint32
	tmr    *ZBtimer
}

var tmr_id = uint32(0)

func TimerAdd(expire uint64, f Func, arg P.Void) uint32 {

	tref := new(ZBtimerRef)
	t := new(ZBtimer)

	if tref == nil || t == nil {
		tref = nil
		t = nil
		return 0
	}

	/* zero is guard for timers */
	if tmr_id == 0 {
		tmr_id++
	}

	if t == nil {
		return 0
	}

	t.arg = arg
	t.function = f
	tref.expire = uint64(time.Now().Unix()) + expire
	tref.id = tmr_id
	tref.tmr = t
	tmr_id++
	timers.Insert(tref)
	//	timerCount++
	if timers.Size() > MAX_TIMERS {
		log.Println("Warning: I have %d timers", timerCount)
	}

	return tref.id
}

func free(t *ZBtimer) {
	log.Printf("%v", t)

}

func check_timers() {

	tref, _ := timers.First().(*ZBtimerRef)

	tick := uint64(time.Now().Unix())
	for tref != nil && tref.expire < tick {
		timer := tref.tmr
		if timer != nil && timer.function != nil {
			timer.function(tick, timer.arg)
		}

		free(timer)
		timer = nil
		timers.Extract()
		tref, _ = timers.First().(*ZBtimerRef)
	}
}

func TimerCancel(id uint32) {

	//    tref := Timers.tree.()
	if id == 0 {
		return
	}
	for i := uint32(0); i <= timers.Size(); i++ {
		tref := timers.Peek(i).(*ZBtimerRef)
		if tref.id == id {
			free(tref.tmr)
			tref.tmr = nil
			break
		}
	}
}
