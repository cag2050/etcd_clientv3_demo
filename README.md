### ectd clientv3 客户端的使用

### 步骤一：运行etcd docker
1.Create a network
```
docker pull bitnami/etcd:3.5
```
2.Launch the Etcd server instance
```
docker run -d --name Etcd-server \
    --network app-tier \
    --publish 2379:2379 \
    --publish 2380:2380 \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    --env ETCD_ADVERTISE_CLIENT_URLS=http://etcd-server:2379 \
    bitnami/etcd:3.5
```
3.Launch your Etcd client instance
```
docker run -it --rm \
    --network app-tier \
    --env ALLOW_NONE_AUTHENTICATION=yes \
    bitnami/etcd:3.5 etcdctl --endpoints http://etcd-server:2379 put /message Hello
```
4.运行命令`docker ps`找到docker实例id

5.进入docker bash
```
docker run -it your_docker_container_id bash
```
6.可以在容器里执行etcdctl命令了
```
etcdctl --endpoints http://etcd-server:2379 get /message
```

### 步骤二：编写etcd clientv3 客户端代码
注意：Endpoints中填写localhost
```
// localhost:2379：本机docker中的etcd服务
Endpoints: []string{"localhost:2379"},
```

资料 | 网址
--- | ---
参考书籍《etcd工作笔记：架构分析、优化与最佳实践》第9章 ectd clientv3 客户端的使用 |
etcd docker 参考 | https://hub.docker.com/r/bitnami/etcd