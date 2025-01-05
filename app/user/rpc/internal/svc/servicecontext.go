package svc

import (
	"cloud-storage/app/user/model"
	"cloud-storage/app/user/rpc/internal/config"
	"cloud-storage/common/init_db"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	UserModel      model.UserModel
	RedisClusterDB *redis.ClusterClient
	MysqlClusterDB *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.MysqlCluster.DataSource)           //创建mysql连接
	mysqlMaster := init_db.InitGorm(c.MysqlCluster.DataSource) //初始化gorm连接
	mysqlMaster.AutoMigrate(                                   //数据库迁移
		&model.User{},
		&model.UserAuth{},
	)
	redisCluster := []string{
		c.RedisCluster.Cluster1,
		c.RedisCluster.Cluster2,
		c.RedisCluster.Cluster3,
		c.RedisCluster.Cluster4,
		c.RedisCluster.Cluster5,
		c.RedisCluster.Cluster6,
	}
	redisDb := init_db.InitRedisCluster(redisCluster)
	//初始化gorm数据库连接
	return &ServiceContext{
		Config:         c,
		UserModel:      model.NewUserModel(conn, c.CacheRedis),
		MysqlClusterDB: mysqlMaster,
		RedisClusterDB: redisDb,
	}

}
