package init_db

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// gorm初始化
func InitGorm(MysqlDataSourece string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(MysqlDataSourece),
		&gorm.Config{})
	if err != nil {
		panic("连接mysql数据库失败, error=" + err.Error())
	} else {
		fmt.Println("连接mysql数据库成功")
	}
	return db
}

func InitRedisCluster(redisCluster []string) *redis.ClusterClient {
	// Redis集群连接参数
	clusterOptions := &redis.ClusterOptions{
		Addrs: redisCluster,
	}
	// 创建Redis集群客户端
	rdb := redis.NewClusterClient(clusterOptions)

	// 连接Redis集群
	//创建一个带有超时的上下文 (ctx) 以控制连接操作的最大等待时间，超时时间设置为 5 秒。cancel 是一个函数，可以在操作完成后被调用来释放资源。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic("连接redis失败, error=" + err.Error())
	}
	fmt.Println("redis连接成功")
	return rdb
}

// kafka日志初始化
func InitKafkaLog() {
	//将日志写入kafka
	logx.SetWriter(*LogXKafka())
}
