package shop

import (
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

	con.Render(c, "itying/buy/checkout.html", gin.H{
		"orderList":   orderList,
		"allPrice":    allPrice,
		"allNum":      allNum,
		"addressList": addressList,
	})
}

func (con BuyController) DoCheckout(c *gin.Context) {
	// 1、获取用户信息 获取用户的收货地址信息
	user := models.User{}
	models.Cookie.Get(c, "userinfo", &user)

	addressResult := models.Address{}
	models.DB.Where("uid = ? AND default_address=1", user.Id).Find(&addressResult)

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
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 3、把订单信息放在订单表，把商品信息放在商品表
		order := models.Order{
			OrderId:     models.GetOrderId(),
			Uid:         user.Id,
			AllPrice:    allPrice,
			Phone:       addressResult.Phone,
			Name:        addressResult.Name,
			Address:     addressResult.Address,
			PayStatus:   0,
			PayType:     0,
			OrderStatus: 0,
			AddTime:     int(models.GetUnix()),
		}

		// 插入订单信息
		if err := tx.Create(&order).Error; err != nil {
			return err // 如果插入订单信息失败，回滚事务
		}

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
			if err := tx.Create(&orderItem).Error; err != nil {
				return err // 如果插入商品信息失败，回滚事务
			}
		}

		return nil // 返回nil表示事务提交
	})

	if err != nil {
		// 事务失败处理
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

	c.Redirect(302, "/buy/pay")
}
