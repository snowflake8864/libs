package libs

import (
	"strings"
	"testing"
	"time"
)

func TestEtcdClient(t *testing.T) {

	etcdc, err := NewEtcdClient(strings.Split("127.0.0.1:2379,10.211.55.74:2379", ","))
	if err != nil {
		t.Fatalf("NewEtcdClient error:%s", err.Error())
	}
	etcdc.WatchAction("/xue")
	for {

		time.Sleep(time.Second * 3)
		etcdc.Get("/xue")
	}
}
