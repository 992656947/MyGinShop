package shop

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"regexp"
	"xiaomiginshop/models"

	"net/http"

	"github.com/gin-gonic/gin"
)

type PassController struct {
	BaseController
}

// 获取验证码
func (con PassController) Captcha(c *gin.Context) {
	id, b64s, err := models.MakeCaptcha(50, 120, 2)

	if err != nil {
		fmt.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"captchaId":    id,
		"captchaImage": b64s,
	})
}

func (con PassController) Login(c *gin.Context) {
	//生成随机数

	a := models.GetRandomNum()
	// c.HTML(http.StatusOK, "itying/pass/login.html", gin.H{})
	c.JSON(200, gin.H{
		"param": a,
	})
}

func (con PassController) RegisterStep1(c *gin.Context) {
	c.HTML(http.StatusOK, "itying/pass/register_step1.html", gin.H{})

}

func (con PassController) RegisterStep2(c *gin.Context) {
	sign := c.Query("sign")
	verifyCode := c.Query("verifyCode")

	session := sessions.Default(c)
	sessionVerifyCode := session.Get("verifyCode")

	sessionVerifyCodeStr, ok := sessionVerifyCode.(string)
	fmt.Println("---------")
	fmt.Println(verifyCode, sessionVerifyCodeStr)
	if !ok || verifyCode != sessionVerifyCodeStr {
		c.Redirect(302, "/pass/registerStep1")
	}

	//2、获取sign 判断sign是否合法

	var userTemp []models.UserTemp
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.HTML(http.StatusOK, "itying/pass/register_step2.html", gin.H{
			"phone":      userTemp[0].Phone,
			"verifyCode": verifyCode,
			"sign":       sign,
		})
	} else {
		c.Redirect(302, "/pass/registerStep1")
	}
}

func (con PassController) RegisterStep3(c *gin.Context) {

	sign := c.Query("sign")
	smsCode := c.Query("smsCode")
	//1、验证短信验证码是否正确
	session := sessions.Default(c)
	sessionSmsCode := session.Get("smsCode")
	sessionSmsCodeStr, ok := sessionSmsCode.(string)
	if !ok || smsCode != sessionSmsCodeStr {
		c.Redirect(302, "/pass/registerStep1")
	}

	//2、获取sign 判断sign是否合法
	var userTemp []models.UserTemp
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) > 0 {
		c.HTML(http.StatusOK, "itying/pass/register_step3.html", gin.H{
			"sign":    sign,
			"smsCode": smsCode,
		})
	} else {
		c.Redirect(302, "/pass/registerStep1")
	}
}

func (con PassController) DoRegister(c *gin.Context) {

	//1、获取表单传过来的数据
	sign, _ := c.FormFile("sign")
	smsCode, _ := c.FormFile("smsCode")
	password, _ := c.FormFile("password")
	rpassword, _ := c.FormFile("rpassword")
	//2、验证smsCode是否合法

	//3、验证密码是否合法

	//4、验证签名是否合法

	//4、完成注册

	//5、执行登录

}

func (con PassController) SendCode(c *gin.Context) {
	phone := c.Query("phone")
	verifyCode := c.Query("verifyCode")
	captchaId := c.Query("captchaId")
	fmt.Println(captchaId, verifyCode)

	// 1、验证图形验证码是否正确 保存
	if captchaId == "resend" {
		// 1、注册第二个页面发送验证码的时候需要验证图形验证码
		sessionDefault := sessions.Default(c)
		sessionVerifyCode := sessionDefault.Get("verifyCode")
		sessionVerifyCodeStr, ok := sessionVerifyCode.(string)
		if !ok || verifyCode != sessionVerifyCodeStr {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "非法请求",
			})
			return
		}
	} else {
		// 1、验证图形验证码是否正确 保存图形验证码
		if flag := models.VerifyCaptcha(captchaId, verifyCode); !flag {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "验证码输入错误，请重试",
			})
			return
		}
		//保存图形验证码
		sessionDefault := sessions.Default(c)
		sessionDefault.Set("verifyCode", verifyCode)
		sessionDefault.Save()
	}

	//2、判断手机格式是否合法

	pattern := `^[\d]{11}$`
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(phone) {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号格式不合法",
		})
		return
	}

	//3、验证手机号是否注册过
	var userList []models.User
	models.DB.Where("phone = ?", phone).Find(&userList)
	if len(userList) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "手机号已经注册，请直接登录",
		})
		return
	}

	//4、判断当前ip地址今天发送短信的次数
	ip := c.ClientIP()
	currentDay := models.GetDay() //20211211
	var sendCount int64
	models.DB.Table("user_temp").Where("ip=? AND add_day=?", ip, currentDay).Count(&sendCount)
	if sendCount > 4 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "此ip今天发送短信的次数已经达到上限，请明天再试",
		})
		return
	}

	//5、验证当前手机号今天发送的次数是否合法
	var userTemp []models.UserTemp
	smsCode := models.GetRandomNum()
	sign := models.Md5(phone + currentDay) //签名：主要用于页面跳转传值
	models.DB.Where("phone = ? AND add_day=?", phone, currentDay).Find(&userTemp)
	if len(userTemp) > 0 {
		if userTemp[0].SendCount > 2 {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "此手机号今天发送短信的次数已经达到上限，请明天再试",
			})
			return
		} else {
			//1、生成短信验证码  发送验证码
			fmt.Println("----------自己集成发送短信的接口--------")
			// 发送短信
			//import ypclnt "github.com/yunpian/yunpian-go-sdk/sdk"
			//clnt := ypclnt.New("apikey")
			//param := ypclnt.NewParam(2)
			//param[ypclnt.MOBILE] = "18616020610"
			//param[ypclnt.TEXT] = "【云片网】您的验证码是1234"
			//r := clnt.Sms().SingleSend(param)
			//
			////账户:clnt.User() 签名:clnt.Sign() 模版:clnt.Tpl() 短信:clnt.Sms() 语音:clnt.Voice() 流量:clnt.Flow()
			fmt.Println(smsCode)
			//2、服务器保持验证码
			session := sessions.Default(c)
			session.Set("smsCode", smsCode)
			session.Save()

			//3、更新发送短信的次数
			oneUserTemp := models.UserTemp{}
			models.DB.Where("id=?", userTemp[0].Id).Find(&oneUserTemp)
			oneUserTemp.SendCount = oneUserTemp.SendCount + 1
			models.DB.Save(&oneUserTemp)

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "发送短信成功",
				"sign":    sign,
			})
			return
		}

	} else {
		//1、生成短信验证码  发送验证码  调用前面课程的接口
		fmt.Println("----------自己集成发送短信的接口--------")
		fmt.Println(smsCode)
		//2、服务器保持验证码
		session := sessions.Default(c)
		session.Set("smsCode", smsCode)
		session.Save()

		//3、记录发送短信的次数

		oneUserTemp := models.UserTemp{
			Ip:        ip,
			Phone:     phone,
			SendCount: 1,
			AddDay:    currentDay,
			AddTime:   int(models.GetUnix()),
			Sign:      sign,
		}
		models.DB.Create(&oneUserTemp)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "发送短信成功",
			"sign":    sign,
		})
		return

	}

	c.JSON(200, gin.H{
		"SendCode": "SendCode",
	})

}

// 验证验证码
func (con PassController) ValidateSmsCode(c *gin.Context) {

	sign := c.Query("sign")
	smsCode := c.Query("smsCode")
	fmt.Println(sign)
	fmt.Println(smsCode)
	//1、验证数据是否合法

	var userTemp []models.UserTemp
	models.DB.Where("sign=?", sign).Find(&userTemp)
	if len(userTemp) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "非法请求",
		})
		return
	}

	//2、验证短信验证码是否正确
	sessionDefault := sessions.Default(c)
	sessionSmsCode := sessionDefault.Get("smsCode")
	sessionSmsCodeStr, ok := sessionSmsCode.(string)
	if !ok || smsCode != sessionSmsCodeStr {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "短信验证码输入错误",
		})
		return
	}
	//3、判断验证码有没有过期   15分

	nowTime := models.GetUnix()
	if (nowTime-int64(userTemp[0].AddTime))/1000/60 > 15 {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "短信验证码已过期",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "验证码输入正确",
	})
}
