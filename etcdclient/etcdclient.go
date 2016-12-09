package etcdclient

import (
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"golang.org/x/net/context"
	// "google.golang.org/grpc"
	"fmt"
	"time"
)

const (
	DialTimeout = time.Second * 3
)

type EtcdClient struct {
	client *clientv3.Client
}

func NewEtcdClient(addrs []string) (*EtcdClient, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   addrs,
		DialTimeout: DialTimeout,
	})
	if err != nil {
		return nil, errors.Annotatef(err, "addrs:%+v", addrs)
	}
	return &EtcdClient{client: cli}, nil
}

//get value from etcd
func (e *EtcdClient) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	resp, err := clientv3.NewKV(e.client).Get(ctx, key)
	cancel()
	if err != nil {
		return "", errors.Trace(err)
	}

	if len(resp.Kvs) == 0 {
		log.Debugf("key:%s value not found", key)
		return "0", nil
	}

	log.Debugf("find key:%s, value:%s", key, string(resp.Kvs[0].Value))
	return string(resp.Kvs[0].Value), nil
}

func (e *EtcdClient) Set(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), DialTimeout)
	_, err := e.client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return errors.Trace(err)
	}
	log.Debugf("set key:%s, value:%s", key, value)
	return nil
}

func (e *EtcdClient) watch(key string, event mvccpb.Event_EventType) {
	watcher := clientv3.NewWatcher(e.client)
	defer watcher.Close()

	for {
		watchChan := watcher.Watch(e.client.Ctx(), key)
		for resp := range watchChan {
			if resp.Canceled {
				return
			}
			for _, ev := range resp.Events {
				if ev.Type == event {
					return
				}
			}
		}
	}
}

//WaitKeyDelete 等待Key被删除
func (e *EtcdClient) WaitKeyDelete(key string) {
	e.watch(key, mvccpb.DELETE)
}

//WaitKeyPut 等待Key被修改
func (e *EtcdClient) WaitKeyPut(key string) {
	e.watch(key, mvccpb.PUT)
}

//WaitKeyPut 等待Key被修改
func (e *EtcdClient) Onlywatch(key string) {
	watcher := clientv3.NewWatcher(e.client)
	defer watcher.Close()

	for {
		watchChan := watcher.Watch(e.client.Ctx(), key)
		for resp := range watchChan {
			if resp.Canceled {
				return
			}
			for _, ev := range resp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}

}

func (e *EtcdClient) WatchAction(key string) {

	go func() {
		e.Onlywatch(key)
	}()

}
