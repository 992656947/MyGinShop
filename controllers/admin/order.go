package admin

import (
	"github.com/gin-gonic/gin"
	"math"
	"xiaomiginshop/models"
)

type OrderController struct {
	BaseController
}

func (con OrderController) Order(c *gin.Context) {
	// 当前页
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	//一页显示多少条
	pageSize := 8

	//获取当前用户
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	//获取当前用户下面订单信息
	var orderList []models.Order
	models.DB.Where("uid = ?", user.Id).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&orderList)

	//获取总数量
	var count int64
	models.DB.Where("uid = ?", user.Id).Table("order").Count(&count)

	c.HTML(200, "admin/order/index.html", gin.H{
		"order":      orderList,
		"page":       page,
		"totalPages": math.Ceil(float64(count) / float64(pageSize)),
	})
}

func (con OrderController) DoEdit(c *gin.Context) {
	c.String(200, "修改成功")
}

func (con OrderController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "删除失败！", "/admin/order")
	}
	orde

}
