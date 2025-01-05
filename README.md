# Cloud-Storage 基于Go-Zero的微服务存储项目

## 架构



## 技术栈

| 功能               | 实现                                     |
| :----------------- |----------------------------------------|
| http框架           | gozero                                 |
| rpc框架            | gozero                                 |
| orm框架            | gorm                                   |
| 数据库             | Innodb-cluster,redis-cluster           |
| 对象存储           | minio集群                      |
| 部署               | Docker,docer-compose                   |
| 服务发现与配置中心 | etcd                                   |
| 消息队列           | kafka                                  |
| 链路追踪           | jaeger                                 |
| 服务监控           | prometheus，grafana                     |
| 日志搜集           | filebeat，go-stash，elasticsearch，kibana |
| 网关               | traefik                                |


### 微服务架构

使用`go-zero`框架，将整个项目拆封为五个服务

1.



### 高可用
