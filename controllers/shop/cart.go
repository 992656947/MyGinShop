package shop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"xiaomiginshop/models"
)

type CartController struct {
	BaseController
}

func (con CartController) AddCart(c *gin.Context) {
	/*
	   购物车数据保持到哪里？：

	           1、购物车数据保存在本地    （cookie）

	           2、购物车数据保存到服务器(mysql)   （必须登录）

	           3、没有登录 购物车数据保存到本地 ， 登录成功后购物车数据保存到服务器  （用的最多）


	   增加购物车的实现逻辑：

	           1、获取增加购物车的数据  （把哪一个商品加入到购物车）

	           2、判断购物车有没有数据   （cookie）

	           3、如果购物车没有任何数据  直接把当前数据写入cookie

	           4、如果购物车有数据

	              1、判断购物车有没有当前数据

	                       有当前数据让当前数据的数量加1，然后写入到cookie

	              2、如果没有当前数据直接写入cookie
	*/

	// 1、获取增加购物车的数据,放在结构体里面  （把哪一个商品加入到购物车）
	colorId, _ := models.Int(c.Query("color_id"))
	goodsId, err := models.Int(c.Query("goods_id"))
	if err != nil {
		c.Redirect(302, "/")
	}
	goods := models.Goods{Id: goodsId}
	goodsColor := models.GoodsColor{Id: colorId}
	models.DB.Find(&goods)
	models.DB.Find(&goodsColor)
	currentData := models.Cart{
		Id:           goodsId,
		Title:        goods.Title,
		Price:        goods.Price,
		GoodsVersion: goods.GoodsVersion,
		Num:          1,
		GoodsColor:   goodsColor.ColorName,
		GoodsImg:     goods.GoodsImg,
		GoodsGift:    goods.GoodsGift, /*赠品*/
		GoodsAttr:    "",              //根据自己的需求拓展
		Checked:      true,            /*默认选中*/
	}

	// 2、判断购物车有没有数据   （cookie）
	var cartList []models.Cart
	models.Cookie.Get(c, "cartList", &cartList)

	if len(cartList) > 0 {
		//4 有数据
		if models.HasCartData(cartList, currentData) {
			for i := 0; i < len(cartList); i++ {
				if cartList[i].Id == currentData.Id && cartList[i].GoodsColor == currentData.GoodsColor && cartList[i].GoodsAttr == currentData.GoodsAttr {
					cartList[i].Num += 1
				}
			}
		} else {
			cartList = append(cartList, currentData)
		}
		models.Cookie.Set(c, "cartList", cartList)
	} else {
		// 3、如果购物车没有任何数据  直接把当前数据写入cookie
		cartList = append(cartList, currentData)
		models.Cookie.Set(c, "cartList", cartList)
	}

	models.Cookie.Get(c, "cartList", &cartList)
	fmt.Println(cartList)

	c.String(200, "加入购物车成功!")
	//models.Cookie.Get(c, "cartList", &cartList)
	//
	//c.JSON(200, gin.H{
	//	"success":  true,
	//	"cartList": cartList,
	//})

}
func (con CartController) Cart(c *gin.Context) {
	//获取购物车数据 显示购物车数据
	var cartList []models.Cart
	models.Cookie.Get(c, "cartList", &cartList)
	fmt.Println(cartList)

	c.JSON(200, gin.H{
		"cartList": cartList,
	})

}
