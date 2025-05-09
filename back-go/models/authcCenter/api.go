package authcCenter

import (
	"errors"
	mysqlDB "gin-web/initialize/mysql"
	"time"

	"gorm.io/gorm"
)

type Api struct {
	Id         int64     `json:"id" gorm:"column:id;type:bigint;primaryKey;not null"`
	Name       string    `json:"name" gorm:"column:name;type:varchar(20);not null"` //api名称
	Url        string    `json:"url" gorm:"column:url;type:varchar(20);not null"`   //地址
	Method     string    `json:"method" gorm:"column:method;type:varchar(8);not null"`
	Desc       string    `json:"desc" gorm:"column:desc;type:varchar(20)"`                 //详情描述
	CreateTime time.Time `json:"createTime" gorm:"column:createTime;autoCreateTime;index"` //创建time
	UpdateTime time.Time `json:"updateTime" gorm:"column:updateTime;default:(-)"`          //修改time
	Roles      []Role    `gorm:"many2many:role_apis;"`                                     //gorm结构体
}

// 查询
func (a *Api) GetList(skip, limit int, startTime, endTime string) ([]Api, int64, error) {
	//总数
	var total int64
	countTx := mysqlDB.DB.Model(&Api{}).Select("")
	if startTime != "" && endTime != "" {
		countTx = countTx.Where("createTime >= ? AND createTime <= ?", startTime, endTime)
	}
	if a.Name != "" {
		countTx = countTx.Where("name = ?", a.Name)
	}
	if a.Url != "" {
		countTx = countTx.Where("url = ?", a.Url)
	}
	if err := countTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//子查询
	subQuery := mysqlDB.DB.Model(&Api{}).Select("id").Order("createTime DESC")
	if startTime != "" && endTime != "" {
		subQuery = subQuery.Where("createTime >= ? AND createTime <= ?", startTime, endTime)
	}
	if a.Name != "" {
		subQuery = subQuery.Where("name = ?", a.Name)
	}
	if a.Url != "" {
		subQuery = subQuery.Where("url = ?", a.Url)
	}
	subQuery = subQuery.Offset(skip).Limit(limit)
	var resDB []Api
	if err := mysqlDB.DB.Model(&Api{}).
		//Select("id", "ip", "name", "url", "method", "desc", "createTime", "updateTime").
		Joins("JOIN (?) AS tmp ON tmp.id = api.id", subQuery).
		Order("createTime DESC").
		Find(&resDB).Error; err != nil {
		return nil, 0, err
	}
	return resDB, total, nil
}

// 添加Api
func (a *Api) Insert() error {
	a.CreateTime = time.Now()
	if err := mysqlDB.DB.Create(a).Error; err != nil {
		return err
	}
	return nil
}

// 删除Api
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

// 修改Api
func (a *Api) Edit() error {
	//修改参数
	a.UpdateTime = time.Now()
	if err := mysqlDB.DB.Updates(a).Error; err != nil {
		return err
	}
	return nil
}

// 查看是否存在
func (a *Api) IsExist() (bool, error) {
	// 查重
	var api Api
	err := mysqlDB.DB.Model(&Api{}).Where("name = ? OR url = ?", a.Name, a.Url).Take(&api).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil // 记录不存在
	}
	if err != nil {
		return false, err // 其他错误
	}
	return true, nil // 记录存在
}
