package etcdclient

import (
	"strings"
	"testing"
	"time"
)

func TestEtcdClient(t *testing.T) {

	etcdc, err := NewEtcdClient(strings.Split("192.168.1.52:2379", ","))
	if err != nil {
		t.Fatalf("NewEtcdClient error:%s", err.Error())
	}
	etcdc.WatchAction("xue")
	for {

		time.Sleep(time.Second * 3)
		etcdc.Get("xue")
		etcdc.Set("xue", "hello")
	}
}
