package main

import (
	"gin-web/initialize/mysql"
	"gin-web/models/authcCenter"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const layout = "2006-01-02 15:04:05" // 时间格式

func main() {
	// 插入的数据
	apiData := []authcCenter.Api{
		{1, "删除user", "user/deleted", "POST", "添加user功能-修改测试", parseTime(layout), parseTime(layout), nil},
		{2, "添加user", "user/add", "POST", "添加user功能", parseTime(layout), parseTime(layout), nil},
		{3, "更新user", "user/update", "PUT", "更改user功能", parseTime(layout), parseTime(layout), nil},
		{4, "查询user", "user/getAll", "GET", "查询user功能", parseTime(layout), parseTime(layout), nil},
		{5, "添加role", "role/add", "POST", "添加role功能", parseTime(layout), parseTime(layout), nil},
		{6, "删除role", "role/deleted", "DELETE", "删除role功能", parseTime(layout), parseTime(layout), nil},
		{7, "更新role", "role/update", "PUT", "更新role功能", parseTime(layout), parseTime(layout), nil},
		{8, "查询role", "role/getAll", "GET", "查询role功能", parseTime(layout), parseTime(layout), nil},
		{9, "添加api", "api/add", "POST", "添加api功能", parseTime(layout), parseTime(layout), nil},
		{10, "删除api", "api/deleted", "DELETE", "删除api功能", parseTime(layout), parseTime(layout), nil},
		{11, "更新api", "api/update", "PUT", "更新api功能", parseTime(layout), parseTime(layout), nil},
		{12, "查询api", "api/getAll", "GET", "查询api功能", parseTime(layout), parseTime(layout), nil},
	}
	// 插入数据
	for _, data := range apiData {
		mysql.DB.Model(&authcCenter.Api{}).Create(&data)
	}
	roleData := []authcCenter.Role{
		{1, "SuperAdmin", "", parseTime(layout), parseTime(layout), nil, nil},
		{2, "Normal", "", parseTime(layout), parseTime(layout), nil, nil},
	}
	for _, data := range roleData {
		if data.Name == "SuperAdmin" {
			data.Add([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
		} else {
			data.Add([]int{})
		}

	}
	//user := []authcCenter.User{
	//	{1,},
	//}
}

// 辅助函数，用于解析字符串时间
func parseTime(timeStr string) time.Time {
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		log.Fatal("时间解析失败:", err)
	}
	return t
}
