package shop

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v2"
	"time"
)

type AlipayController struct{}

func (con AlipayController) Alipay(c *gin.Context) {
	//1、获取订单号 判断此订单号是否值当前用户的

	//2、获取订单里面的支付信息

	//支付宝公钥
	var aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhb/KlxYfhRE8KRp92MQM8ZB8NVjoM9LYFOnPIuNtcMZVA8ld7ybDP2FiA+QEE7wLGqMImwl1Y4xzkrTLCjHVC8fdR8ZvzZR2I3ZOrARerI9+RbkCfT+7YLv55+A+WTHEyiB+v7PfXVTT28s0CHNLPXMyQD1u8UVEQEpbMSs8hH3pJF55Li7kc5VvJpV3RVO9TXZTVAA5mSp9FvO3u+47IJDgFVLnqqHh6ETL1nHVpxiAY2LGer+RWpVYD8v+We+VWsrfJP7bO0xr2pwizldepo8YNYPgcIAIwd7KiveypL1pA0xWgSjUHzrkVh1j/nSnvJgKSdydU/VRcaVt/Mt8wwIDAQAB" // 可选，支付宝提供给我们用于签名验证的公钥，通过支付宝管理后台获取
	//应用私钥
	var privateKey = "MIIEowIBAAKCAQEAmZlUXdYU5qhpSe7bmv+NAEQehXq2wkV2pAhr9zU2fKJSB77V1Z9079SYTs0yGX3XpGKqVYVGZu6F3js9QX1/PPfELHrxOa4mHXQz+TryniXx1e1ng4kDedy1SyS8xTBeGCvReIqAtMUeAo/hedgvdtgYEvY5hZqbV7b3zEpEkq/XgvgaH5nmYR/fDnocwQt5OpDlXNGMrfhU4g0nvQHzb8Kcd2YSjkDB37o/yV9i0Wq8vJqPdWCJOGxCOafy1Qt1VfPF6s9h/9V8mWTGrMGyaUUgz/0ILdlt/bHzRSR9K5qa0x9mEfnGZRj14x3WGIlNfSmS88a4wR0iXP+UtEASXwIDAQABAoIBABoUDV3tNhk/aLjzw/daAh+UcTYqcpMjZhRNlb8gGsMocBL+lKGzdBAwITfn4OSxGAbB9beVbDGXt8TWe/z9iLfaPUVsDj7D0ZbYnuZm2sB9IsU2jIepoJx1G5bJgv9bye4CqorzwQxwFztKIHcmfFCKOfQmN/f2Gv/WgdX+mgvpaNPfHn3gnnNpmhq2Mb8QiNhLvPMhKXORk08a2xGLUEO43Iw7YSg8lZ4L6So/ROSPCtyErAgVTxAkfKoFyJAaHykKPeeLt9Zy5aPndmizIWwCMXb0+Zq7bDmKS3yvq8Nr95ZjOc09GD1nit2vWLCPJxNc2lzkSBxIAmDmXta0+LECgYEA6LMl02AgB47zxXfvhDSpW3SABrzU8GyOnZjcYdf2cZkDtgx0QE3g5E2YeWFcSny3+x3ug6nVjjIscZaBap3DxwMft2PrGX7Bfmecg3hJ+TF7KLFuAkmG+qqzkPeox9YTs/aRBRVBu9oo+Io44HqryIoehCpX59+CjqznhLEX8FcCgYEAqPqR3wJPI38XqGLKoxpI20l99rGD/mj2J3PvFXGtti3e5ZB4irsxO8U38vHEw/2kEXEJZW+2trS486pXaKE8JfrRItjaCUbgUSO+3lCpdkl9gXfimEr0xD2zPf92aAcxqQiG65I+k2Ch5ba2YmteTGi55SB6FuEAvHtGhRf7iTkCgYEAlFQ1pVJduFOwIcx8uaoT1j8hqKnPll2sXtrkh93wsqKV0gKIS8EYvI6VxbGA8d4kLIb81aJ5hUWIPPNyFTLxa7cbDXw8jSjWUCvdgZQ4mwamed73v698weXzxlGHnbJhJtLhx/qvxv2eJid9b+HiBFe+cgLHu/8mKqoefd+g4csCgYAAsOWfz9abAo4KNj015Ymeu/Iz7A3qIGvBRYwYvlpDgHSE485aYuGUqP3NlIeFdagSGjA7pfVNUfffpzasStyAG0J3rgNWPl/0dPz208WdojdNLDxU+xl9I/NzsXO+gSkG0+4ZUIPI/oAq/FBKnr3H+jWoZjWZmlnya16idLKmoQKBgEU5tCv8klIdPLwHu0UOrMWtbPnWABzCu2eJdova0zJblHE5KTgu3U2td/wEfiDJVYpnK3XSZkxhjJWKxxtWSQwUjZwK00pc2J/wxW1wPq7uXI0/jiJT/C4BDwDor20wMm8hqOgdCxPT1yIbD25WUF7QVRBUuX/wQHJfdrsPKzUe" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	//应用id
	var appId = "2021001144694156"
	//最后一个参数 false 表示沙箱环境 true表示正式环境
	client, _ := alipay.New(appId, aliPublicKey, privateKey, true)
	// 加载支付宝公钥证书

	// 将 key 的验证调整到初始化阶段

	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://118.123.14.36:8005/alipayNotify" //post
	p.ReturnURL = "http://127.0.0.1:8005/alipayReturn"
	p.Subject = "测试 公钥证书模式-这是一个gin订单" //标题 修改
	template := "2006-01-02 15:04:05"
	p.OutTradeNo = time.Now().Format(template)
	p.TotalAmount = "0.1" //金额 修改
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	var url, err4 = client.TradePagePay(p)
	if err4 != nil {
		fmt.Println(err4)
	}

	var payURL = url.String()
	fmt.Println(payURL)

	c.Redirect(302, payURL)

}
func (con AlipayController) AlipayNotify(c *gin.Context) {
	fmt.Println("AlipayNotify")

	var aliPublicKey = "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAhb/KlxYfhRE8KRp92MQM8ZB8NVjoM9LYFOnPIuNtcMZVA8ld7ybDP2FiA+QEE7wLGqMImwl1Y4xzkrTLCjHVC8fdR8ZvzZR2I3ZOrARerI9+RbkCfT+7YLv55+A+WTHEyiB+v7PfXVTT28s0CHNLPXMyQD1u8UVEQEpbMSs8hH3pJF55Li7kc5VvJpV3RVO9TXZTVAA5mSp9FvO3u+47IJDgFVLnqqHh6ETL1nHVpxiAY2LGer+RWpVYD8v+We+VWsrfJP7bO0xr2pwizldepo8YNYPgcIAIwd7KiveypL1pA0xWgSjUHzrkVh1j/nSnvJgKSdydU/VRcaVt/Mt8wwIDAQAB"                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           // 可选，支付宝提供给我们用于签名验证的公钥，通过支付宝管理后台获取
	var privateKey = "MIIEowIBAAKCAQEAmZlUXdYU5qhpSe7bmv+NAEQehXq2wkV2pAhr9zU2fKJSB77V1Z9079SYTs0yGX3XpGKqVYVGZu6F3js9QX1/PPfELHrxOa4mHXQz+TryniXx1e1ng4kDedy1SyS8xTBeGCvReIqAtMUeAo/hedgvdtgYEvY5hZqbV7b3zEpEkq/XgvgaH5nmYR/fDnocwQt5OpDlXNGMrfhU4g0nvQHzb8Kcd2YSjkDB37o/yV9i0Wq8vJqPdWCJOGxCOafy1Qt1VfPF6s9h/9V8mWTGrMGyaUUgz/0ILdlt/bHzRSR9K5qa0x9mEfnGZRj14x3WGIlNfSmS88a4wR0iXP+UtEASXwIDAQABAoIBABoUDV3tNhk/aLjzw/daAh+UcTYqcpMjZhRNlb8gGsMocBL+lKGzdBAwITfn4OSxGAbB9beVbDGXt8TWe/z9iLfaPUVsDj7D0ZbYnuZm2sB9IsU2jIepoJx1G5bJgv9bye4CqorzwQxwFztKIHcmfFCKOfQmN/f2Gv/WgdX+mgvpaNPfHn3gnnNpmhq2Mb8QiNhLvPMhKXORk08a2xGLUEO43Iw7YSg8lZ4L6So/ROSPCtyErAgVTxAkfKoFyJAaHykKPeeLt9Zy5aPndmizIWwCMXb0+Zq7bDmKS3yvq8Nr95ZjOc09GD1nit2vWLCPJxNc2lzkSBxIAmDmXta0+LECgYEA6LMl02AgB47zxXfvhDSpW3SABrzU8GyOnZjcYdf2cZkDtgx0QE3g5E2YeWFcSny3+x3ug6nVjjIscZaBap3DxwMft2PrGX7Bfmecg3hJ+TF7KLFuAkmG+qqzkPeox9YTs/aRBRVBu9oo+Io44HqryIoehCpX59+CjqznhLEX8FcCgYEAqPqR3wJPI38XqGLKoxpI20l99rGD/mj2J3PvFXGtti3e5ZB4irsxO8U38vHEw/2kEXEJZW+2trS486pXaKE8JfrRItjaCUbgUSO+3lCpdkl9gXfimEr0xD2zPf92aAcxqQiG65I+k2Ch5ba2YmteTGi55SB6FuEAvHtGhRf7iTkCgYEAlFQ1pVJduFOwIcx8uaoT1j8hqKnPll2sXtrkh93wsqKV0gKIS8EYvI6VxbGA8d4kLIb81aJ5hUWIPPNyFTLxa7cbDXw8jSjWUCvdgZQ4mwamed73v698weXzxlGHnbJhJtLhx/qvxv2eJid9b+HiBFe+cgLHu/8mKqoefd+g4csCgYAAsOWfz9abAo4KNj015Ymeu/Iz7A3qIGvBRYwYvlpDgHSE485aYuGUqP3NlIeFdagSGjA7pfVNUfffpzasStyAG0J3rgNWPl/0dPz208WdojdNLDxU+xl9I/NzsXO+gSkG0+4ZUIPI/oAq/FBKnr3H+jWoZjWZmlnya16idLKmoQKBgEU5tCv8klIdPLwHu0UOrMWtbPnWABzCu2eJdova0zJblHE5KTgu3U2td/wEfiDJVYpnK3XSZkxhjJWKxxtWSQwUjZwK00pc2J/wxW1wPq7uXI0/jiJT/C4BDwDor20wMm8hqOgdCxPT1yIbD25WUF7QVRBUuX/wQHJfdrsPKzUe" // 必须，上一步中使用 RSA签名验签工具 生成的私钥
	var appId = "2021001144694156"

	var client, err1 = alipay.New(appId, aliPublicKey, privateKey, true)
	if err1 != nil {
		fmt.Println(err1)
		return
	}

	req := c.Request
	req.ParseForm()
	ok, _ := client.VerifySign(req.Form)

	fmt.Println(ok)

	fmt.Println(req.Form)

	c.String(200, "ok")
}
func (con AlipayController) AlipayReturn(c *gin.Context) {
	c.String(200, "支付成功")
}
