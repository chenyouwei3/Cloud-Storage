package authCenter

import (
	"fmt"
	mysqlDB "gin-web/internal/initialize/mysql"
	"time"

	"gorm.io/gorm"
)

type Api struct {
	Id         int64     `json:"id" gorm:"column:id;type:bigint;primaryKey;not null"`
	Name       string    `json:"name" gorm:"column:name;type:varchar(20);not null"` //api名称
	Url        string    `json:"url" gorm:"column:url;type:varchar(50);not null"`   //地址
	Method     string    `json:"method" gorm:"column:method;type:varchar(8);not null"`
	Desc       string    `json:"desc" gorm:"column:desc;type:varchar(20)"`                   //详情描述
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime;index"` //创建time
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:(-)"`          //修改time
	//Roles      []Role    `gorm:"many2many:role_apis;"`                                       //gorm结构体
}

func (Api) tableName() string {
	return "api"
}

// 查询
func (a *Api) GetList(skip, limit int, startTime, endTime string) ([]Api, int64, error) {
	//总数
	var total int64
	countTx := mysqlDB.DB.Model(&Api{})
	if startTime != "" && endTime != "" {
		countTx = countTx.Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if a.Name != "" {
		countTx = countTx.Where("name LIKE ?", "%"+a.Name+"%") // 模糊查询 name
	}
	if a.Url != "" {
		countTx = countTx.Where("url LIKE ?", "%"+a.Url+"%") // 模糊查询 url
	}
	if a.Method != "" {
		countTx = countTx.Where("method = ?", a.Method)
	}
	if err := countTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//子查询
	subQuery := mysqlDB.DB.Model(&Api{}).Select("id").Order("create_time DESC")
	if startTime != "" && endTime != "" {
		subQuery = subQuery.Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if a.Name != "" {
		subQuery = subQuery.Where("name LIKE ?", "%"+a.Name+"%") // 模糊查询 name
	}
	if a.Url != "" {
		subQuery = subQuery.Where("url LIKE ?", "%"+a.Url+"%") // 模糊查询 url
	}
	if a.Method != "" {
		subQuery = subQuery.Where("method = ?", a.Method)
	}
	subQuery = subQuery.Offset(skip).Limit(limit)
	var resDB []Api
	if err := mysqlDB.DB.Model(&Api{}).
		//Select("id", "ip", "name", "url", "method", "desc", "create_time", "updateTime").
		Joins("JOIN (?) AS tmp ON tmp.id = api.id", subQuery).
		Order("create_time DESC").
		Find(&resDB).Error; err != nil {
		return nil, 0, err
	}
	return resDB, total, nil
}

// 添加
func (a *Api) Insert() error {
	a.CreateTime = time.Now()
	fmt.Println("api_testing", a)
	if err := mysqlDB.DB.Model(&Api{}).Create(a).Error; err != nil {
		return err
	}
	return nil
}

// 删除
func (a *Api) Remove(id int64) error {
	//启动事务
	return mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		// 删除 Api 记录
		err := tx.Where("id = ?", id).Delete(&Api{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// 修改
func (a *Api) Edit() error {
	//修改参数
	a.UpdateTime = time.Now()
	if err := mysqlDB.DB.Updates(a).Error; err != nil {
		return err
	}
	return nil
}

// exist
func (a *Api) IsExist() (bool, error) {
	var exists bool
	err := mysqlDB.DB.Model(&Api{}).
		Select("1").
		Where("name = ? OR url = ?", a.Name, a.Url).
		Limit(1).
		Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return exists, nil
}
