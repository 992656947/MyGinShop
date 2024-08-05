package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaomiginshop/models"
)

type SettingController struct {
	BaseController
}

func (con SettingController) Index(c *gin.Context) {
	settingList := models.Setting{}
	models.DB.First(&settingList)
	c.HTML(http.StatusOK, "admin/setting/index.html", gin.H{
		"setting": settingList,
	})

}

func (con SettingController) DoEdit(c *gin.Context) {
	setting := models.Setting{Id: 1}
	models.DB.Find(&setting)
	err := c.ShouldBind(&setting)
	if err != nil {
		con.Error(c, "修改数据成功", "/admin/setting")
		return
	} else {
		//上传图片
		siteLogo, err1 := models.UploadImg(c, "site_logo")
		if len(siteLogo) > 0 && err1 == nil {
			setting.SiteLogo = siteLogo
		}
		noPicture, err2 := models.UploadImg(c, "no_picture")
		if len(noPicture) > 0 && err2 == nil {
			setting.NoPicture = noPicture
		}
		err3 := models.DB.Save(&setting).Error
		if err3 != nil {
			con.Error(c, "修改数据成功", "/admin/setting")
			return
		}
	}

}
