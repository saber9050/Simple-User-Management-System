package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrPictureFile = errors.New("只允许上传图片")      //上传文件类型错误
	ErrPictureSize = errors.New("图片大小不能超过5MB") //上传文件大小错误
)

/*	保存头像	*/

// 保存上传的头像，返回可访问路径
func SaveUploadedFile(r *http.Request, fileName string, userID int) (string, error) {
	file, handler, err := r.FormFile(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	//	验证文件类型
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".png" && ext != ".gif" && ext != ".jpeg" {
		return "", ErrPictureFile
	}
	//	限制文件大小	5MB
	if handler.Size > 5*1024*1024 {
		return "", ErrPictureSize
	}
	//	生成新文件，名字：id_当前时间戳.ext	时间是指从January 1, 1970 UTC到现在经过的总时间(单位纳秒)
	newFileName := fmt.Sprintf("%d_%d%s", userID, time.Now().UnixNano(), ext)
	savePath := filepath.Join("uploads", newFileName)
	//	创建目录，如果有就跳过
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		return "", err
	}
	//	写入文件
	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}
	//	返回web访问路径
	return "/" + savePath, nil
}

// 删除旧的头像
func DeleteUploadedFile(avatarPath string) error {
	// 默认头像不删除
	if avatarPath == "/uploads/default.png" {
		return nil
	}

	// 将 Web 路径转换为文件系统路径，删除前缀”/“
	filePath := strings.TrimPrefix(avatarPath, "/")
	if filePath == "" {
		return nil
	}

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // 文件不存在，无需删除
	}

	// 删除文件
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("删除旧头像失败: %w", err)
	}
	return nil
}
