package lease_keep_alive

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
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

func keepLease(lease clientv3.Lease, leaseId clientv3.LeaseID) {
	var err error = nil
	keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId)
	if err != nil {
		fmt.Printf("KeepAlive err: %+v\n", err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Printf("租约已经失效了\n")
					goto END
				} else {
					fmt.Printf("收到自动续租应答：%+v\n", keepResp.ID)
				}
			}
		}
	END:
	}()
}

func TestKeepLease(t *testing.T) {
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

	leaseGrantResp, err = lease.Grant(context.TODO(), 3)
	if err != nil {
		fmt.Printf("Grant err: %+v\n", err)
		return
	}
	leaseId = leaseGrantResp.ID

	keepLease(lease, leaseId)
	// 让主线程别马上退出
	time.Sleep(20 * time.Second)
}
