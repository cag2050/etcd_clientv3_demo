package app

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.etcd.io/etcd/clientv3"
	"testing"
	"time"
)

func TestKV(t *testing.T) {
	rootContext := context.Background()
	client, err := clientv3.New(clientv3.Config{
		// localhost:2379：本机docker中的etcd服务
		Endpoints: []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if client == nil || err == context.DeadlineExceeded {
		fmt.Println(err)
		panic("invalid connection!")
	}
	defer client.Close()

	kv := clientv3.KV(client)
	ctx, cancelFunc := context.WithTimeout(rootContext, time.Duration(2) * time.Second)
	response, err := kv.Get(ctx, "cc")
	cancelFunc()
	if err != nil {
		fmt.Printf("Get err: %+v\n", err)
	}
	kvs := response.Kvs
	if len(kvs) > 0 {
		fmt.Printf("last value is: %s\n", string(kvs[0].Value))
	} else {
		fmt.Printf("empty key for %s\n", "cc")
	}

	uuidStr := uuid.New().String()
	fmt.Printf("new value is: %s\n", uuidStr)
	ctx2, cancelFunc2 := context.WithTimeout(rootContext, time.Duration(2) * time.Second)

	_, err = kv.Put(ctx2, "cc", uuidStr)

	delRes, err := kv.Delete(ctx2, "cc")
	if err != nil {
		fmt.Printf("Delete err: %+v\n", err)
	} else {
		fmt.Printf("delete %s for %t\n", "cc", delRes.Deleted > 0)
	}

	cancelFunc2()
}