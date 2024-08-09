package shop

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"xiaomiginshop/models"
)

type DefaultController struct{}

func (con DefaultController) Index(c *gin.Context) {
	//1、获取顶部导航
	var topNavList []models.Nav
	if hasTopNavList := models.CacheDb.Get("topNavList", &topNavList); !hasTopNavList {
		models.DB.Where("status=1 AND position=1").Find(&topNavList)
		models.CacheDb.Set("topNavList", topNavList, 60*60)
	}

	//2、获取轮播图数据
	var focusList []models.Focus
	if hasFocusList := models.CacheDb.Get("focusList", &focusList); !hasFocusList {
		models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
		models.CacheDb.Set("focusList", focusList, 60*60)
	}

	//3、获取分类的数据
	var goodsCateList []models.GoodsCate

	if hasGoodsCateList := models.CacheDb.Get("goodsCateList", &goodsCateList); !hasGoodsCateList {
		models.DB.Where("pid = 0 AND status=1").Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
			return db.Where("goods_cate.status=1").Order("goods_cate.sort DESC")
		}).Find(&goodsCateList)

		models.CacheDb.Set("goodsCateList", goodsCateList, 60*60)
	}

	//4. 中间导航
	var middleNavList []models.Nav
	if hasMiddleNavList := models.CacheDb.Get("middleNavList", &middleNavList); !hasMiddleNavList {
		models.DB.Where("status=1 AND position=2").Find(&middleNavList)
		for i := 0; i < len(middleNavList); i++ {
			relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",") //21，22,23,24
			relationIds := strings.Split(relation, ",")
			var goodsList []models.Goods
			models.DB.Where("id in ?", relationIds).Select("id,title,goods_img,price").Find(&goodsList)
			middleNavList[i].GoodsItems = goodsList
		}
		models.CacheDb.Set("middleNavList", middleNavList, 60*60)
	}

	//手机
	var phoneList []models.Goods
	if hasPhoneList := models.CacheDb.Get("phoneList", &phoneList); !hasPhoneList {
		phoneList = models.GetGoodsByCategory(1, "best", 8)
		models.CacheDb.Set("phoneList", phoneList, 60*60)
	}

	//配件
	var otherList []models.Goods
	if hasOtherList := models.CacheDb.Get("otherList", &otherList); !hasOtherList {
		otherList = models.GetGoodsByCategory(9, "all", 1)
		models.CacheDb.Set("otherList", otherList, 60*60)
	}

	c.HTML(http.StatusOK, "itying/index/index.html", gin.H{
		"topNavList":    topNavList,
		"focusList":     focusList,
		"goodsCateList": goodsCateList,
		"middleNavList": middleNavList,
		"phoneList":     phoneList,
		"otherList":     otherList,
	})
}

//func (con DefaultController) Thumbnail1(c *gin.Context) {
//	//按宽度进行比例缩放，输入输出都是文件
//	//filename string, savepath string, width int
//	filename := "static/upload/0.png"
//	savepath := "static/upload/0_600.png"
//	err := ScaleF2F(filename, savepath, 600)
//	if err != nil {
//		c.String(200, "生成图片失败")
//		return
//	}
//	c.String(200, "Thumbnail1 成功")
//}
//
//func (con DefaultController) Thumbnail2(c *gin.Context) {
//	filename := "static/upload/tao.jpg"
//	savepath := "static/upload/tao_400.png"
//	//按宽度和高度进行比例缩放，输入和输出都是文件
//	err := ThumbnailF2F(filename, savepath, 400, 400)
//	if err != nil {
//		c.String(200, "生成图片失败")
//		return
//	}
//	c.String(200, "Thumbnail2 成功")
//}
//
//func (con DefaultController) Qrcode1(c *gin.Context) {
//	var png []byte
//	png, err := qrcode.Encode("https://www.itying.com", qrcode.Medium, 256)
//	if err != nil {
//		c.String(200, "生成二维码失败")
//		return
//	}
//	c.String(200, string(png))
//}
//
//func (con DefaultController) Qrcode2(c *gin.Context) {
//	savepath := "static/upload/qrcode.png"
//	err := qrcode.WriteFile("https://www.itying.com", qrcode.Medium, 556, savepath)
//	if err != nil {
//		c.String(200, "生成二维码失败")
//		return
//	}
//	file, _ := ioutil.ReadFile(savepath)
//	c.String(200, string(file))
//}
