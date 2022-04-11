package app

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/google/uuid"
	"testing"
	"time"
)

var (
	leaseId        clientv3.LeaseID
	getResp        *clientv3.GetResponse
	leaseGrantResp *clientv3.LeaseGrantResponse
	kv             clientv3.KV
	keepResp       *clientv3.LeaseKeepAliveResponse
	keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
)

func TestLease(t *testing.T) {
	rootContext := context.Background()
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if client == nil || err == context.DeadlineExceeded {
		fmt.Printf("clientv3.New err: %+v\n", err)
		panic("invalid connection!")
	}
	defer client.Close()

	lease := clientv3.NewLease(client)

	leaseGrantResp, err = lease.Grant(context.TODO(), 10)
	if err != nil {
		fmt.Printf("Grant err: %+v\n", err)
		return
	}
	leaseId = leaseGrantResp.ID

	kv = clientv3.NewKV(client)
	uuid := uuid.New().String()
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
	// 将租期Id绑定到指定的key
	_, err = kv.Put(ctx, "dd", uuid, clientv3.WithLease(leaseId))
	if err != nil {
		fmt.Printf("Put err: %+v\n", err)
		return
	}
	cancelFunc()

	for {
		ctx2, cancelFunc2 := context.WithTimeout(rootContext, time.Duration(2)*time.Second)
		getResp, err = kv.Get(ctx2, "dd")
		if err != nil {
			fmt.Printf("Get err: %+v\n", err)
			return
		}
		cancelFunc2()

		if getResp.Count == 0 {
			fmt.Printf("kv 过期了\n")
			break
		}
		fmt.Printf("还没过期：%+v\n", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}
