package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
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
	pageSize := 3

	//获取当前用户
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	//获取当前用户下面订单信息
	var orderList []models.Order

	//获取keyword
	keyword := c.Query("keyword")
	if keyword == "" {
		models.DB.Where("uid = ? ", user.Id).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&orderList)
	} else {
		models.DB.Where("uid = ? AND order_id = ?", user.Id, keyword).Offset((page - 1) * pageSize).Limit(pageSize).Preload("OrderItem").Order("add_time desc").Find(&orderList)
	}

	//获取总数量
	var count int64
	models.DB.Where("uid = ?", user.Id).Table("order").Count(&count)

	if len(orderList) > 0 {
		c.HTML(http.StatusOK, "admin/order/index.html", gin.H{
			"order": orderList,
			//注意float64类型
			"totalPages": math.Ceil(float64(count) / float64(pageSize)),
			"page":       page,
			"keyword":    keyword,
		})
	} else {
		if page != 1 {
			c.Redirect(302, "/admin/order")
		} else {
			c.HTML(200, "admin/order/index.html", gin.H{
				"order":      orderList,
				"page":       page,
				"totalPages": math.Ceil(float64(count) / float64(pageSize)),

				"keyword": keyword,
			})
		}

	}

}

func (con OrderController) Edit(c *gin.Context) {

	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "查询失败！", "/admin/order")
	}

	var OrderItems []models.OrderItem
	var user models.User

	models.Cookie.Get(c, "userinfo", &user)
	models.DB.Where("order_id = ? AND uid = ?", id, user.Id).Find(&OrderItems)

	order := models.Order{Id: id}
	models.DB.Find(&order)
	//c.JSON(200, gin.H{
	//	"orderItems": OrderItems,
	//})
	c.HTML(200, "admin/order/edit.html", gin.H{
		"order":     order,
		"orderItem": OrderItems,
	})
}

func (con OrderController) DoEdit(c *gin.Context) {
	// 获取订单ID
	orderId, _ := strconv.Atoi(c.PostForm("orderId"))
	fmt.Println(orderId)

	// 查询要修改的订单
	order := models.Order{}
	if err := models.DB.Where("order_id = ?", orderId).First(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "订单不存在"})
		return
	}

	// 更新订单信息
	order.Name = c.PostForm("name")
	order.Phone = c.PostForm("phone")
	order.Address = c.PostForm("address")
	order.AllPrice, _ = strconv.ParseFloat(c.PostForm("allPrice"), 64)
	order.PayStatus, _ = strconv.Atoi(c.PostForm("payStatus"))
	order.OrderStatus, _ = strconv.Atoi(c.PostForm("orderStatus"))

	// 保存订单信息
	if err := models.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "订单更新失败"})
		return
	}

	// 获取并更新订单项信息
	var orderItems []models.OrderItem
	if err := models.DB.Where("order_id = ?", orderId).Find(&orderItems).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取订单项失败"})
		return
	}

	for i := range orderItems {
		orderItemId := strconv.Itoa(orderItems[i].Id)

		orderItems[i].ProductTitle = c.PostForm("orderItems[" + orderItemId + "][productTitle]")
		orderItems[i].ProductImg = c.PostForm("orderItems[" + orderItemId + "][productImg]")
		orderItems[i].ProductPrice, _ = strconv.ParseFloat(c.PostForm("orderItems["+orderItemId+"][productPrice]"), 64)
		orderItems[i].ProductNum, _ = strconv.Atoi(c.PostForm("orderItems[" + orderItemId + "][productNum]"))
		orderItems[i].GoodsVersion = c.PostForm("orderItems[" + orderItemId + "][goodsVersion]")
		orderItems[i].GoodsColor = c.PostForm("orderItems[" + orderItemId + "][goodsColor]")

		// 保存每个订单项信息
		if err := models.DB.Save(&orderItems[i]).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "订单项更新失败"})
			return
		}
	}

	// 返回到订单列表页面
	c.Redirect(http.StatusFound, "/admin/order")
}

func (con OrderController) View(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "查询失败！", "/admin/order")
	}

	var OrderItems []models.OrderItem
	var user models.User

	models.Cookie.Get(c, "userinfo", &user)
	models.DB.Where("order_id = ? AND uid = ?", id, user.Id).Find(&OrderItems)

	order := models.Order{Id: id}
	models.DB.Find(&order)
	//c.JSON(200, gin.H{
	//	"orderItems": OrderItems,
	//})
	c.HTML(200, "admin/order/view.html", gin.H{
		"order":      order,
		"orderItems": OrderItems,
	})

}

// 主页面删除
func (con OrderController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "删除失败！", "/admin/order")
	}

	err1 := models.DB.Transaction(func(tx *gorm.DB) error {
		//删除子表
		var OrderItem []models.OrderItem
		var user models.User
		models.Cookie.Get(c, "userinfo", &user)
		if err2 := tx.Where("order_id = ? AND uid = ?", id, user.Id).Delete(&OrderItem).Error; err2 != nil {
			return err2
		}

		//删除父表
		Order := models.Order{Id: id}
		if err3 := tx.Delete(&Order).Error; err3 != nil {
			return err3
		}

		return nil
	})
	//事务失败处理
	if err1 != nil {
		c.String(500, "订单处理失败，请重试")
		return

	}
}

func (con OrderController) Check(c *gin.Context) {

}
