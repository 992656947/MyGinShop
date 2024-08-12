package models

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
)

type ginCookie struct {
}

//写入

func (cookie ginCookie) Set(c *gin.Context, key string, value interface{}) {
	bytes, _ := json.Marshal(value)
	c.SetCookie(key, string(bytes), 3600, "/", "localhost", false, true)
}

// 获取
func (cookie ginCookie) Get(c *gin.Context, key string, obj interface{}) bool {

	valueStr, err1 := c.Cookie(key)
	if err1 == nil && valueStr != "" && valueStr != "[]" {
		err2 := json.Unmarshal([]byte(valueStr), obj)
		return err2 == nil
	}
	return false
}

// 实例化结构体
var Cookie = &ginCookie{}
