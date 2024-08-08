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
	topNavList := []models.Nav{}
	models.DB.Where("status=1 AND position=1").Find(&topNavList)
	//2、获取轮播图数据
	focusList := []models.Focus{}
	models.DB.Where("status=1 AND focus_type=1").Find(&focusList)
	//3、获取分类的数据
	goodsCateList := []models.GoodsCate{}
	//
	models.DB.Where("pid = 0 AND status=1").Order("sort DESC").Preload("GoodsCateItems", func(db *gorm.DB) *gorm.DB {
		return db.Where("goods_cate.status=1").Order("goods_cate.sort DESC")
	}).Find(&goodsCateList)

	//4. 中间导航
	middleNavList := []models.Nav{}
	models.DB.Where("status=1 AND position=2").Find(&middleNavList)
	for i := 0; i < len(middleNavList); i++ {
		relation := strings.ReplaceAll(middleNavList[i].Relation, "，", ",") //21，22,23,24
		relationIds := strings.Split(relation, ",")
		goodsList := []models.Goods{}
		models.DB.Where("id in ?", relationIds).Select("id,title,goods_img,price").Find(&goodsList)
		middleNavList[i].GoodsItems = goodsList
	}

	c.HTML(http.StatusOK, "itying/index/index.html", gin.H{
		"topNavList":    topNavList,
		"focusList":     focusList,
		"goodsCateList": goodsCateList,
		"middleNavList": middleNavList,
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
