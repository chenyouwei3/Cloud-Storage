package test

import (
	"context"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func TestRedis(t *testing.T) {
	// Redis集群连接参数
	clusterOptions := &redis.ClusterOptions{
		Addrs: []string{
			server106_ip + ":6379",
			server106_ip + ":6380",
			server106_ip + ":6381",
			server106_ip + ":6382",
			server106_ip + ":6383",
			server106_ip + ":6384",
		},
	}

	// 创建Redis集群客户端
	client := redis.NewClusterClient(clusterOptions)

	// 连接Redis集群
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping Redis以确认连接
	err := client.Ping(ctx).Err()
	if err != nil {
		t.Fatal("连接失败:", err)
	}

	// 执行 CLUSTER INFO 命令
	info, err := client.ClusterInfo(ctx).Result()
	if err != nil {
		t.Fatal("Failed to execute CLUSTER INFO:", err)
	}
	t.Log("Cluster Info:", info) // 使用 t.Log 输出集群信息

	// 设置并获取值
	client.Set(ctx, "hhh", 5369, 0)
	val, err := client.Get(ctx, "hhh").Result()
	if err != nil {
		t.Fatal("Failed to get value:", err)
	}

	// 验证返回值
	if val != "5369" {
		t.Errorf("Expected value '5369', but got '%s'", val)
	}
}

func TestRedisClusterMaster(t *testing.T) {
	// 创建 Redis 集群客户端实例
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			server106_ip + ":6379", // 节点1的主机和端口
			server106_ip + ":6381", // 节点2的主机和端口
			server106_ip + ":6383", // 节点3的主机和端口

		},
	})

	// 测试连接
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		t.Fatal("连接失败:", err)
	}
	t.Log("连接成功:", pong) // 使用 t.Log 输出连接状态

	// 设置并获取值
	client.Set(context.Background(), "aaaa", "eeeee", 0)
	val, err := client.Get(context.Background(), "aaaa").Result()
	if err != nil {
		t.Fatal("Failed to get value:", err)
	}

	// 验证返回值
	if val != "eeeee" {
		t.Errorf("Expected value 'eeeee', but got '%s'", val)
	}

	// 关闭客户端连接
	err = client.Close()
	if err != nil {
		t.Fatal("关闭连接失败:", err)
	}
}
