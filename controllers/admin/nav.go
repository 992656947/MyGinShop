package admin

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"xiaomiginshop/models"
)

type NavController struct {
	BaseController
}

func (con NavController) Index(c *gin.Context) {
	navList := []models.Nav{}
	models.DB.Find(&navList)
	//fmt.Println(navList)
	c.HTML(200, "admin/nav/index.html", gin.H{
		"navList": navList,
	})

}
func (con NavController) Add(c *gin.Context) {
	c.HTML(200, "admin/nav/add.html", gin.H{})

}
func (con NavController) DoAdd(c *gin.Context) {

	title := c.PostForm("title")
	link := c.PostForm("link")
	position, _ := models.Int(c.PostForm("position"))
	isOpennew, _ := models.Int(c.PostForm("isOpennew"))
	relation := c.PostForm("relation")
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	if title == "" {
		con.Error(c, "标题不能为空", "/admin/nav/add")
		return
	}

	nav := models.Nav{
		Title:     title,
		Link:      link,
		Position:  position,
		IsOpennew: isOpennew,
		Relation:  relation,
		Sort:      sort,
		Status:    status,
		AddTime:   int(models.GetUnix()),
	}
	err := models.DB.Create(&nav).Error
	if err != nil {
		con.Error(c, "增加导航失败 请重试", "/admin/nav/add")
	} else {
		con.Success(c, "增加导航成功", "/admin/nav")
	}

}
func (con NavController) Edit(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传入数据错误", "/admin/nav")
	} else {
		nav := models.Nav{Id: id}
		models.DB.Find(&nav)
		//fmt.Println(nav)
		c.HTML(http.StatusOK, "admin/nav/edit.html", gin.H{
			"nav": nav,
		})
	}
}
func (con NavController) DoEdit(c *gin.Context) {
	id, err := models.Int(c.PostForm("id"))
	if err != nil {
		con.Error(c, "修改失败！", "/admin/nav/edit?id="+models.String(id))
	}
	title := c.PostForm("title")
	link := c.PostForm("link")
	position, _ := models.Int(c.PostForm("position"))
	isOpennew, _ := models.Int(c.PostForm("isOpennew"))
	relation := c.PostForm("relation")
	sort, _ := models.Int(c.PostForm("sort"))
	status, _ := models.Int(c.PostForm("status"))

	nav := models.Nav{Id: id}
	models.DB.Find(&nav)
	nav.Title = title
	nav.Link = link
	nav.Position = position
	nav.IsOpennew = isOpennew
	nav.Relation = relation
	nav.Sort = sort
	nav.Status = status
	nav.AddTime = int(models.GetUnix())
	err1 := models.DB.Save(&nav).Error
	if err1 != nil {
		con.Error(c, "修改失败", "/admin/nav/edit?id="+models.String(id))
		return
	}
	con.Success(c, "修改成功", "/admin/nav")

}

func (con NavController) Delete(c *gin.Context) {
	id, err := models.Int(c.Query("id"))
	if err != nil {
		con.Error(c, "传参错误!", "/admin/nav")
		return
	}
	nav := models.Nav{Id: id}
	models.DB.Delete(&nav)
	con.Success(c, "删除成功！", "/admin/nav")
}
