package routers

import (
	"xiaomiginshop/controllers/shop"

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

	}
}
