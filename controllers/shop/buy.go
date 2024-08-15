package shop

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"xiaomiginshop/models"
)

type BuyController struct {
	BaseController
}

func (con BuyController) Checkout(c *gin.Context) {

	var cartList []models.Cart
	models.Cookie.Get(c, "cartList", &cartList)

	var orderList []models.Cart

	var allPrice float64
	var allNum int

	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
			allNum += cartList[i].Num
		}
	}
	//fmt.Println(orderList)
	//fmt.Println(orderList[0].GoodsImg)

	var user models.User
	models.Cookie.Get(c, "userinfo", &user)
	//fmt.Println(user.Id)
	var addressList []models.Address
	models.DB.Where("uid = ?", user.Id).Find(&addressList)
	//fmt.Println(addressList)

	//签名
	orderSign := models.Md5(models.GetRandomNum())
	session := sessions.Default(c)
	session.Set("orderSign", orderSign)
	session.Save()

	//orderList是否存在
	if len(orderList) == 0 {
		c.Redirect(302, "/")
		return
	}

	con.Render(c, "itying/buy/checkout.html", gin.H{
		"orderList":   orderList,
		"allPrice":    allPrice,
		"allNum":      allNum,
		"addressList": addressList,
		"orderSign":   orderSign,
	})
}

func (con BuyController) DoCheckout(c *gin.Context) {
	//防止重复提交
	orderSignClient := c.PostForm("orderSign")
	session := sessions.Default(c)
	orderSignSession := session.Get("orderSign")
	orderSignServer, ok := orderSignSession.(string)
	fmt.Println("执行到OK")
	if !ok {
		c.Redirect(302, "/")
		return
	}
	fmt.Println("ok执行结束，判断orderSignClient != orderSignServer开始")
	fmt.Println(orderSignClient)
	fmt.Println(orderSignServer)
	if orderSignClient != orderSignServer {
		c.Redirect(302, "/")
		return
	}
	fmt.Println("判断结束")
	session.Delete("orderSign")
	session.Save()

	// 1、获取用户信息 获取用户的收货地址信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	var addressResult []models.Address
	models.DB.Where("uid = ? AND default_address=1", user.Id).Find(&addressResult)

	if len(addressResult) == 0 {
		c.Redirect(302, "/buy/checkout")
		return
	}

	// 2、获取购买商品的信息
	var cartList []models.Cart
	models.Cookie.Get(c, "cartList", &cartList)
	var orderList []models.Cart
	var allPrice float64
	for i := 0; i < len(cartList); i++ {
		if cartList[i].Checked {
			allPrice += cartList[i].Price * float64(cartList[i].Num)
			orderList = append(orderList, cartList[i])
		}
	}

	// 开启事务
	var myOrderID int
	err1 := models.DB.Transaction(func(tx *gorm.DB) error {
		// 3、把订单信息放在订单表，把商品信息放在商品表
		order := models.Order{
			OrderId:     models.GetOrderId(),
			Uid:         user.Id,
			AllPrice:    allPrice,
			Phone:       addressResult[0].Phone,
			Name:        addressResult[0].Name,
			Address:     addressResult[0].Address,
			PayStatus:   0,
			PayType:     0,
			OrderStatus: 0,
			AddTime:     int(models.GetUnix()),
		}
		// 插入订单信息
		if err2 := tx.Create(&order).Error; err2 != nil {
			return err2
		}
		myOrderID = order.Id
		// 插入商品信息
		for i := 0; i < len(orderList); i++ {
			orderItem := models.OrderItem{
				OrderId:      order.Id,
				Uid:          user.Id,
				ProductTitle: orderList[i].Title,
				ProductId:    orderList[i].Id,
				ProductImg:   orderList[i].GoodsImg,
				ProductPrice: orderList[i].Price,
				ProductNum:   orderList[i].Num,
				GoodsVersion: orderList[i].GoodsVersion,
				GoodsColor:   orderList[i].GoodsColor,
			}
			// 返回nil表示事务提交
			if err3 := models.DB.Create(&orderItem).Error; err3 != nil {
				return err3
			}
		}

		return nil
	})

	// 事务失败处理
	if err1 != nil {
		c.String(500, "订单处理失败，请重试")
		return
	}

	// 4、删除购物车里面的选中数据
	var noSelectCartList []models.Cart
	for i := 0; i < len(cartList); i++ {
		if !cartList[i].Checked {
			noSelectCartList = append(noSelectCartList, cartList[i])
		}
	}
	models.Cookie.Set(c, "cartList", noSelectCartList)

	if myOrderID == 0 {
		c.Redirect(302, "/")
		return
	}

	c.Redirect(302, "/buy/pay?orderId="+models.String(myOrderID))
}

func (con BuyController) Pay(c *gin.Context) {
	orderId, err := models.Int(c.Query("orderId"))
	if err != nil {
		c.Redirect(302, "/")
	}

	//获取用户
	var user models.User
	models.Cookie.Get(c, "userinfo", &user)

	order := models.Order{}
	models.DB.Where("id = ?", orderId).Find(&order)

	if user.Id != order.Uid {
		c.Redirect(302, "/")
		return
	}

	var orderItems []models.OrderItem
	models.DB.Where("order_id = ?", order.Id).Find(&orderItems)
	var allPrice float64
	for i := 0; i < len(orderItems); i++ {
		allPrice += orderItems[i].ProductPrice
	}

	con.Render(c, "itying/buy/pay.html", gin.H{
		"orderItems": orderItems,
		"order":      order,
		"allPrice":   allPrice,
	})

}
