package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

func main() {
	router := gin.Default()

	// 上传文件接口
	router.POST("/upload", UploadFile)

	router.Run(":8080")
}

// UploadFile 处理文件上传
func UploadFile(c *gin.Context) {
	// 从表单中获取文件，字段名为 "file"
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "获取文件失败",
			"error":   err.Error(),
		})
		return
	}

	// 保存路径（这里保存到当前项目的 upload 目录下）
	savePath := filepath.Join("upload", file.Filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "保存文件失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "上传成功",
		"filename": file.Filename,
		"path":     savePath,
	})
}
