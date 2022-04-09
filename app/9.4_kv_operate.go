package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	client, err := clientv3.New(clientv3.Config{
		// localhost:2379：本机docker中的etcd服务
		Endpoints: []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if client == nil || err == context.DeadlineExceeded {
		// handle errors
		fmt.Println(err)
		panic("invalid connection!")
	}

	kv := clientv3.KV(client)
	putResp, putErr := kv.Put(context.TODO(), "/aa", "hello-world2!")
	if putErr != nil {
		fmt.Printf("putErr:\n %+v\n", putErr)
	} else {
		fmt.Printf("putResp:\n %+v\n", putResp)
		getResp, getErr := kv.Get(context.TODO(), "/aa")
		if getErr != nil {
			fmt.Printf("getErr:\n %+v\n", getErr)
		} else {
			fmt.Printf("getResp:\n %+v\n", getResp)
		}
	}

	deleteResp, deleteErr := kv.Delete(context.TODO(), "/aa")
	if deleteErr != nil {
		fmt.Printf("deleteErr:\n %+v\n", deleteErr)
	} else {
		fmt.Printf("deleteResp:\n %+v\n", deleteResp)
	}

	client.Close()
}