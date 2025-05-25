# Cloud-Storage 基于Gin的网盘服务项目

## 🚀技术栈

| 功能     | 实现                    |
|:-------|-----------------------|
| http框架 | gin                   |
| orm框架  | gorm                  |
| 数据库    | Innodb-cluster        |
| 部署     | Docker,docer-compose  |
| Web前端  | Vue3 / ant-design-vue |
## 开发环境

Go v1.20

Node.js v18.18.0


## 数据初始化

进入back-go/migrate 运行脚本

## 启动

在internal/initialize/config当中配置config.yml格式如下

```yaml
// FileName: config.yml
APP:
  name:    #服务名称
  ip: 127.0.0.1 
  port: 8080   
  mode: run   #运行模式
  staticFS: true  #是否开静态资源访问


MySQL:
  driverName: 
  host: 
  port: 
  database: 
  username: 
  password: 
  charset: utf8mb4

```
## 部署

进入front-vue过后运行build_dist.sh 过后进行打包

后续进入back-go/cmd运行build_docker_image.sh打包成image进行运行

可以根据需求运行clean_docker.sh清楚容器与镜像

