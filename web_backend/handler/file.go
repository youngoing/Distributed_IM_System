package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// 静态文件配置结构体
type StaticConfig struct {
	UploadPath   string   // 上传文件保存路径
	AllowedTypes []string // 允许的文件类型
	MaxFileSize  int64    // 最大文件大小（字节）
}

// 默认配置
var DefaultStaticConfig = StaticConfig{
	UploadPath:   "../uploads",
	AllowedTypes: []string{".jpg", ".jpeg", ".png", ".gif"},
	MaxFileSize:  5 << 20, // 5MB
}

// 检查文件类型是否允许
func isAllowedImageType(ext string) bool {
	for _, allowedType := range DefaultStaticConfig.AllowedTypes {
		if ext == allowedType {
			return true
		}
	}
	return false
}

// 处理文件上传
func HandleImgUpload(c *gin.Context) {
	id := c.Query("id")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// 检查文件大小
	if file.Size > DefaultStaticConfig.MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// 检查文件类型
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isAllowedImageType(ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	// 生成唯一文件名
	filename := generateUniqueFilename(id)
	filepath := filepath.Join(DefaultStaticConfig.UploadPath, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// 返回文件访问URL
	fileURL := "/uploads/" + filename
	c.JSON(http.StatusOK, gin.H{
		"url":      fileURL,
		"filename": filename,
	})
}

// 生成唯一文件名，使用传入的 id 拼接 "avatar" 和 ".jpg"
func generateUniqueFilename(id string) string {
	// 拼接 id 和 "avatar"，然后使用 ".jpg" 作为文件扩展名
	filename := fmt.Sprintf("%s_avatar", id)
	return fmt.Sprintf("%s.jpg", filename)
}

// // 删除文件
// func Delete(id string) error {
// 	// 拼接文件名
// 	filename := fmt.Sprintf("%s_avatar.jpg", id)
// 	// 检查文件是否存在
// 	if _, err := os.Stat(filename); err != nil {
// 		if os.IsNotExist(err) {
// 			// 文件不存在
// 			fmt.Printf("文件 %s 不存在\n", filename)
// 			return nil
// 		}
// 		// 如果检查文件时发生其他错误，返回错误
// 		return err
// 	}
// 	// 文件存在，执行删除
// 	err := os.Remove(filename)
// 	if err != nil {
// 		return fmt.Errorf("删除文件 %s 失败: %w", filename, err)
// 	}

// 	// 删除成功
// 	fmt.Printf("文件 %s 已删除\n", filename)
// 	return nil
// }
