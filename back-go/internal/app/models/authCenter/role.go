package authCenter

import (
	"fmt"
	mysqlDB "gin-web/internal/initialize/mysql"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	Id         int64     `json:"id" gorm:"column:id;type:bigint;primaryKey;not null"`
	Name       string    `json:"name" gorm:"column:name;type:varchar(20);not null"`          //角色名称
	Desc       string    `json:"desc" gorm:"column:desc;type:varchar(20)"`                   //详情描述
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime;index"` //创建time
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:(-)"`          //修改time
	Apis       []Api     `gorm:"many2many:role_apis;"`
	//Users      []User    `gorm:"many2many:user_roles"`
}

func (Role) tableName() string {
	return "role"
}

// 查询
func (r *Role) GetList(skip, limit int, startTime, endTime string) ([]Role, int64, error) {
	//总数
	var total int64
	countTx := mysqlDB.DB.Model(&Role{})
	if startTime != "" && endTime != "" {
		countTx = countTx.Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if r.Name != "" {
		countTx = countTx.Where("name LIKE ?", "%"+r.Name+"%") // 模糊查询 name
	}
	if err := countTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//子查询
	subQuery := mysqlDB.DB.Model(&Role{}).Select("id").Order("create_time DESC")
	if startTime != "" && endTime != "" {
		subQuery = subQuery.Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if r.Name != "" {
		subQuery = subQuery.Where("name LIKE ?", "%"+r.Name+"%") // 模糊查询 name
	}
	subQuery = subQuery.Offset(skip).Limit(limit)

	var resDB []Role

	if err := mysqlDB.DB.Model(&Role{}).
		//Select("id", "ip", "name", "url", "method", "desc", "createTime", "updateTime").
		//Preload("Apis").
		Joins("JOIN (?) AS tmp ON tmp.id = role.id", subQuery).
		Order("create_time DESC").
		Find(&resDB).Error; err != nil {
		return nil, 0, err
	}
	return resDB, total, nil
}

// 添加Role
func (r *Role) Insert(apiIds []int) error {
	r.CreateTime = time.Now()
	return mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		// 创建 Role 记录
		res := tx.Create(r)
		if res.Error != nil {
			return res.Error
		}
		// 查找所有指定的 Api 记录
		var apis []Api
		if err := tx.Find(&apis, apiIds).Error; err != nil {
			return err
		}
		// 确保所有 apiIds 都存在
		if len(apis) != len(apiIds) {
			return fmt.Errorf("api数量不匹配")
		}
		// 关联 Api 到 Role
		if err := tx.Model(&Role{Id: r.Id}).Association("Apis").Append(apis); err != nil {
			return err
		}
		return nil
	})
}

// 删除Role
func (r *Role) Remove(id int64) error {
	//删除role,受制于user/api
	return mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		// 清除 Api 与 Roles 的关联关系
		err := tx.Model(&Role{Id: id}).Association("Apis").Clear()
		if err != nil {
			return err
		}
		// 删除 Role 记录
		err = tx.Where("id = ?", id).Delete(&Role{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (r *Role) Edit(addApis, deletedApis []int) error {
	r.UpdateTime = time.Now()
	err := mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		// 更新角色基本信息
		if err := tx.Model(&Role{Id: r.Id}).Updates(r).Error; err != nil {
			return fmt.Errorf("更新角色信息失败: %w", err)
		}
		// 删除关联
		if len(deletedApis) > 0 {
			if err := tx.Table("role_apis").Where("role_id = ? AND api_id IN ?", r.Id, deletedApis).Delete(nil).Error; err != nil {
				return fmt.Errorf("删除关联失败: %w", err)
			}
		}
		// 添加关联
		if len(addApis) > 0 {
			records := make([]map[string]interface{}, len(addApis))
			for i, apiID := range addApis {
				records[i] = map[string]interface{}{
					"role_id": r.Id,
					"api_id":  apiID,
				}
			}
			if err := tx.Table("role_apis").Create(records).Error; err != nil {
				return fmt.Errorf("添加关联失败: %w", err)
			}
		}
		return nil
	})
	return err
}

// 查看是否存在
func (r *Role) IsExist() (bool, error) {
	var exists bool
	err := mysqlDB.DB.Model(&Role{}).
		Select("1").
		Where("name = ?", r.Name).
		Limit(1).
		Scan(&exists).Error
	if err != nil {
		return false, err
	}
	return true, nil // 记录存在
}

func GetApisByRoleID(roleID int64) ([]Api, error) {
	var role Role
	err := mysqlDB.DB.Model(&Role{}).
		Preload("Apis", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		Where("id = ?", roleID).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return role.Apis, nil
}
