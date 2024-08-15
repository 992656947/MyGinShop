package routers

import (
	"xiaomiginshop/controllers/shop"
	"xiaomiginshop/middlewares"

	"github.com/gin-gonic/gin"
)

func DefaultRoutersInit(r *gin.Engine) {
	defaultRouters := r.Group("/")
	{
		defaultRouters.GET("/", shop.DefaultController{}.Index)
		//defaultRouters.GET("/thumbnail1", shop.DefaultController{}.Thumbnail1)
		//defaultRouters.GET("/thumbnail2", shop.DefaultController{}.Thumbnail2)
		//defaultRouters.GET("/qrcode1", shop.DefaultController{}.Qrcode1)
		//defaultRouters.GET("/qrcode2", shop.DefaultController{}.Qrcode2)
		defaultRouters.GET("/category:id", shop.ProductController{}.Category)
		defaultRouters.GET("/detail", shop.ProductController{}.Detail)
		defaultRouters.GET("/product/getImgList", shop.ProductController{}.GetImgList)
		defaultRouters.GET("/cart", shop.CartController{}.Get)
		defaultRouters.GET("/cart/addCart", shop.CartController{}.AddCart)

		defaultRouters.GET("/cart/successTip", shop.CartController{}.AddCartSuccess)

		defaultRouters.GET("/cart/decCart", shop.CartController{}.DecCart)
		defaultRouters.GET("/cart/incCart", shop.CartController{}.IncCart)

		defaultRouters.GET("/cart/changeOneCart", shop.CartController{}.ChangeOneCart)
		defaultRouters.GET("/cart/changeAllCart", shop.CartController{}.ChangeAllCart)
		defaultRouters.GET("/cart/delCart", shop.CartController{}.DelCart)

		defaultRouters.GET("/pass/login", shop.PassController{}.Login)
		defaultRouters.POST("/pass/doLogin", shop.PassController{}.DoLogin)
		defaultRouters.GET("/pass/loginOut", shop.PassController{}.LoginOut)
		defaultRouters.GET("/pass/captcha", shop.PassController{}.Captcha)

		defaultRouters.GET("/pass/registerStep1", shop.PassController{}.RegisterStep1)
		defaultRouters.GET("/pass/registerStep2", shop.PassController{}.RegisterStep2)
		defaultRouters.GET("/pass/registerStep3", shop.PassController{}.RegisterStep3)
		defaultRouters.GET("/pass/sendCode", shop.PassController{}.SendCode)
		defaultRouters.GET("/pass/validateSmsCode", shop.PassController{}.ValidateSmsCode)
		defaultRouters.POST("/pass/doRegister", shop.PassController{}.DoRegister)

		defaultRouters.GET("/buy/checkout", middlewares.InitUserAuthMiddleware, shop.BuyController{}.Checkout)
		defaultRouters.POST("/buy/doCheckout", middlewares.InitUserAuthMiddleware, shop.BuyController{}.DoCheckout)
		defaultRouters.GET("/buy/pay", middlewares.InitUserAuthMiddleware, shop.BuyController{}.Pay)

		defaultRouters.POST("/address/addAddress", middlewares.InitUserAuthMiddleware, shop.AddressController{}.AddAddress)
		defaultRouters.POST("/address/editAddress", middlewares.InitUserAuthMiddleware, shop.AddressController{}.EditAddress)
		defaultRouters.GET("/address/changeDefaultAddress", middlewares.InitUserAuthMiddleware, shop.AddressController{}.ChangeDefaultAddress)
		defaultRouters.GET("/address/getOneAddressList", middlewares.InitUserAuthMiddleware, shop.AddressController{}.GetOneAddressList)

	}
}
