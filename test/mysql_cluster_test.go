package test

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestMySQLConnection(t *testing.T) {
	// MySQL 连接配置
	dsn := "root:root@tcp(" + server106_ip + ":6446)/"
	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("无法打开数据库连接: %v", err)
	}
	defer db.Close()

	// 测试数据库连接
	err = db.Ping()
	if err != nil {
		t.Fatalf("无法连接到数据库: %v", err)
	}
	t.Log("成功连接到 MySQL 数据库!")

	// 可选：执行简单的查询
	// 例如，列出当前数据库中的所有数据库
	rows, err := db.Query("SHOW DATABASES;")
	if err != nil {
		t.Fatalf("查询失败: %v", err)
	}
	defer rows.Close()

	t.Log("数据库列表:")
	for rows.Next() {
		var dbName string
		if err := rows.Scan(&dbName); err != nil {
			t.Fatalf("读取数据失败: %v", err)
		}
		t.Log(" -", dbName)
	}

	// 检查迭代过程中是否有错误
	if err = rows.Err(); err != nil {
		t.Fatalf("行迭代错误: %v", err)
	}
}
