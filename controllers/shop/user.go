package shop

import (
	"math"
	"xiaomiginshop/models"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	BaseController
}

func (con UserController) Index(c *gin.Context) {
	// c.String(http.StatusOK, "首页")
	var tpl = "itying/user/welcome.html"
	con.Render(c, tpl, gin.H{})
}
func (con UserController) OrderList(c *gin.Context) {
	// 当前页
	page, _ := models.Int(c.Query("page"))
	if page == 0 {
		page = 1
	}
	//一页显示多少条
	pageSize := 2

	//获取当前用户
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	//获取当前用户下面订单信息
	var orderList []models.Order
	models.DB.Where("uid = ?", user.Id).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&orderList)

	//获取总数量
	var count int64
	models.DB.Where("uid = ?", user.Id).Table("order").Count(&count)

	var tpl = "itying/user/order.html"
	con.Render(c, tpl, gin.H{
		"order":      orderList,
		"page":       page,
		"totalPages": math.Ceil(float64(count) / float64(pageSize)),
	})
}
func (con UserController) OrderInfo(c *gin.Context) {

	id, err := models.Int(c.Query("id"))
	if err != nil {
		c.Redirect(302, "/user/order")
	}
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)
	var order []models.Order
	models.DB.Where("id=? And uid=?", id, user.Id).Preload("OrderItem").Find(&order)

	if len(order) == 0 {
		c.Redirect(302, "/user/order")
		return
	}
	var tpl = "itying/user/order_info.html"
	con.Render(c, tpl, gin.H{
		"order": order[0],
	})
}
