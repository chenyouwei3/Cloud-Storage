package authCenter

import (
	"errors"
	"fmt"
	mysqlDB "gin-web/internal/initialize/mysql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	Id         int64     `json:"id" gorm:"column:id;type:bigint;primaryKey;not null"`
	Name       string    `json:"name" gorm:"column:name;type:varchar(35);not null"`          //用户名
	Account    string    `json:"account" gorm:"column:account;type:varchar(35);not null"`    //账号
	Password   string    `json:"password" gorm:"column:password;type:varchar(100);not null"` //密码
	AvatarUrl  string    `json:"avatarUrl" gorm:"column:avatarUrl;type:varchar(50)"`         //头像Url
	Sex        string    `json:"sex" gorm:"column:sex;type:varchar(3);not null"`             //性别
	Email      string    `json:"email" gorm:"column:email;type:varchar(35);not null"`        //邮箱
	Salt       string    `json:"salt" gorm:"column:salt;type:varchar(35);not null"`          //盐加密
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;autoCreateTime;index"` //创建time
	UpdateTime time.Time `json:"update_time" gorm:"column:update_time;default:(-)"`          //修改time
	Roles      []Role    `gorm:"many2many:user_roles"`                                       //外键role
}

// 查询
func (u *User) GetList(skip, limit int, startTime, endTime string) ([]User, int64, error) {
	//总数
	var total int64
	countTx := mysqlDB.DB.Model(&User{})
	if startTime != "" && endTime != "" {
		countTx = countTx.Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if u.Name != "" {
		countTx = countTx.Where("name LIKE ?", "%"+u.Name+"%") // 模糊查询 name
	}
	if u.Email != "" {
		countTx = countTx.Where("email LIKE ?", "%"+u.Email+"%") // 模糊查询 email
	}
	if err := countTx.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//子查询
	subQuery := mysqlDB.DB.Model(&User{}).Select("id").Order("create_time DESC")
	if startTime != "" && endTime != "" {
		subQuery = subQuery.Where("create_time >= ? AND create_time <= ?", startTime, endTime)
	}
	if u.Name != "" {
		subQuery = subQuery.Where("name LIKE ?", "%"+u.Name+"%") // 模糊查询 name
	}
	if u.Email != "" {
		subQuery = subQuery.Where("email LIKE ?", "%"+u.Email+"%") // 模糊查询 email
	}
	subQuery = subQuery.Offset(skip).Limit(limit)
	var resDB []User
	if err := mysqlDB.DB.Model(&User{}).
		//Select("id", "ip", "name", "url", "method", "desc", "createTime", "updateTime").
		Joins("JOIN (?) AS tmp ON tmp.id = user.id", subQuery).
		Order("create_time DESC").
		Find(&resDB).Error; err != nil {
		return nil, 0, err
	}
	return resDB, total, nil
}

func (u *User) Insert(roleIds []int) error {
	u.CreateTime = time.Now()
	return mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		res := tx.Create(u)
		if res.Error != nil {
			return res.Error
		}
		// 查找所有指定的 roles 记录
		var roles []Role
		if err := tx.Find(&roles, roleIds).Error; err != nil {
			return err
		}
		// 确保所有 roles 都存在
		if len(roles) != len(roleIds) {
			return fmt.Errorf("role数量不匹配")
		}
		// 关联 Role 到 User
		if err := tx.Model(&User{Id: u.Id}).Association("Roles").Append(roles); err != nil {
			return err
		}
		return nil
	})
}

func (u *User) Remove(id int64) error {
	//删除role,受制于user/api
	return mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		// 清除 User 与 Roles 的关联关系
		err := tx.Model(&User{Id: id}).Association("Roles").Clear()
		if err != nil {
			return err
		}
		// 删除 User 记录
		err = tx.Where("id = ?", id).Delete(&User{}).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (u *User) Edit(addRoles, deletedRoles []int) error {
	u.UpdateTime = time.Now()
	err := mysqlDB.DB.Transaction(func(tx *gorm.DB) error {
		//更新用户基本信息
		if err := tx.Model(&User{Id: u.Id}).Updates(u).Error; err != nil {
			return fmt.Errorf("更新用户信息失败: %w", err)
		}
		// 删除关联
		if len(deletedRoles) > 0 {
			if err := tx.Table("user_roles").Where("user_id = ? AND role_id IN ?", u.Id, deletedRoles).Delete(nil).Error; err != nil {
				return fmt.Errorf("删除关联失败: %w", err)
			}
		}
		// 添加关联
		if len(addRoles) > 0 {
			records := make([]map[string]interface{}, len(addRoles))
			for i, roleId := range addRoles {
				records[i] = map[string]interface{}{
					"user_id": u.Id,
					"role_id": roleId,
				}
			}
			if err := tx.Table("user_roles").Create(records).Error; err != nil {
				return fmt.Errorf("添加关联失败: %w", err)
			}
		}
		return nil
	})
	return err
}

func (u *User) GetOne() (*User, error) {
	var user User
	query := mysqlDB.DB.Model(&User{}).Preload("Roles", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name") // 只查 Role 的 id 和 name 字段
	})
	if u.Account != "" {
		query = query.Where("account = ?", u.Account)
	} else if u.Name != "" {
		query = query.Where("name = ?", u.Name)
	} else if u.Email != "" {
		query = query.Where("email = ?", u.Email)
	} else {
		return nil, errors.New("account 和 name 和 email 不能同时为空")
	}
	err := query.Take(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *User) IsExist() (bool, error) {
	var exists bool
	query := mysqlDB.DB.Model(&User{}).Select("1")
	if u.Account != "" {
		query.Where("account = ?", u.Account)
	} else if u.Email != "" {
		query.Where("email = ?", u.Email)
	} else if u.Name != "" {
		query.Where("name = ?", u.Email)
	}
	if err := query.Limit(1).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}
