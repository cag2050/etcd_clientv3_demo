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
	} else {
		fmt.Println(client.Cluster.MemberList(context.TODO()))
	}
	client.Close()
}