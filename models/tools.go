package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

// 时间戳转换成日期
func UnixToTime(timestamp int) string {
	t := time.Unix(int64(timestamp), 0)
	return t.Format("2006-01-02 15:04:05")
}

// 日期转换成时间戳 2020-05-02 15:04:05
func DateToUnix(str string) int64 {
	template := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(template, str, time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// 获取时间戳
func GetUnix() int64 {
	return time.Now().Unix()
}

// 获取纳秒
func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

// 获取当前的日期
func GetDate() string {
	template := "2006-01-02 15:04:05"
	return time.Now().Format(template)
}

// 获取年月日
func GetDay() string {
	template := "20060102"
	return time.Now().Format(template)
}

// md5加密
func Md5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// 把字符串解析成html
func Str2Html(str string) template.HTML {
	return template.HTML(str)
}

// 表示把string转换成int
func Int(str string) (int, error) {
	n, err := strconv.Atoi(str)
	return n, err
}

// 表示把string转换成Float64
func Float(str string) (float64, error) {
	n, err := strconv.ParseFloat(str, 64)
	return n, err
}

// 表示把int转换成string
func String(n int) string {
	str := strconv.Itoa(n)
	return str
}

////上传图片
//func UploadImg(c *gin.Context, picName string) (string, error) {
//	ossStatus := GetOssStatus()
//	if ossStatus == 1 {
//		return OssUploadImg(c, picName)
//	} else {
//		return LocalUploadImg(c, picName)
//	}
//
//}
//

// Oss上传
//func OssUplod(file *multipart.FileHeader, dst string) (string, error) {
//
//	f, err := file.Open()
//	if err != nil {
//		return "", err
//	}
//	defer f.Close()
//
//	// 创建OSSClient实例。
//	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "GJoqWHXB2c9S9gwP", "Lgf3weXuWITUUb17vDJfveg1jmKEe9")
//	if err != nil {
//		return "", err
//	}
//
//	// 获取存储空间。
//	bucket, err := client.Bucket("beego")
//	if err != nil {
//		return "", err
//	}
//
//	// 上传文件流。
//	err = bucket.PutObject(dst, f)
//	if err != nil {
//		return "", err
//	}
//	return dst, nil
//}

////上传图片到Oss
//func OssUploadImg(c *gin.Context, picName string) (string, error) {
//	// 1、获取上传的文件
//	file, err := c.FormFile(picName)
//
//	if err != nil {
//		return "", err
//	}
//
//	// 2、获取后缀名 判断类型是否正确  .jpg .png .gif .jpeg
//	extName := path.Ext(file.Filename)
//	allowExtMap := map[string]bool{
//		".jpg":  true,
//		".png":  true,
//		".gif":  true,
//		".jpeg": true,
//	}
//
//	if _, ok := allowExtMap[extName]; !ok {
//		return "", errors.New("文件后缀名不合法")
//	}
//
//	// 3、定义图片保存目录  static/upload/20210624
//
//	day := GetDay()
//	dir := "static/upload/" + day
//
//	// 4、生成文件名称和文件保存的目录   111111111111.jpeg
//	fileName := strconv.FormatInt(GetUnixNano(), 10) + extName
//
//	// 5、执行上传
//	dst := path.Join(dir, fileName)
//
//	OssUplod(file, dst)
//	return dst, nil
//
//}

// 上传图片
func UploadImg(c *gin.Context, picName string) (string, error) {
	// 1、获取上传的文件
	file, err := c.FormFile(picName)
	if err != nil {
		return "", err
	}

	// 2、获取后缀名 判断类型是否正确  .jpg .png .gif .jpeg
	extName := path.Ext(file.Filename)
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
	}

	if _, ok := allowExtMap[extName]; !ok {
		return "", errors.New("文件后缀名不合法")
	}

	// 3、创建图片保存目录  static/upload/20210624

	day := GetDay()
	dir := "./static/upload/" + day

	err1 := os.MkdirAll(dir, 0666)
	if err1 != nil {
		fmt.Println(err1)
		return "", err1
	}

	// 4、生成文件名称和文件保存的目录   111111111111.jpeg
	fileName := strconv.FormatInt(GetUnixNano(), 10) + extName

	// 5、执行上传
	dst := path.Join(dir, fileName)
	c.SaveUploadedFile(file, dst)
	return dst, nil

}

func Sub(a int, b int) int {
	return a - b
}

// Substr截取字符串
func Substr(str string, start int, end int) string {
	rs := []rune(str)
	rl := len(rs)
	if start < 0 {
		start = 0
	}
	if start > rl {
		start = 0
	}

	if end < 0 {
		end = rl
	}
	if end > rl {
		end = rl
	}
	if start > end {
		start, end = end, start
	}

	return string(rs[start:end])

}
