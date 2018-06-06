## swiftpass(威富通)的支付客户端sdk
##### 已完成
- 微信扫码支付
- 微信扫码支付
- 微信js支付
- 微信h5支付
- 支付宝扫码支付
- 支付宝js支付
- 订单查询
- 订单关闭
- 订单退款

##### 用法:

   ```go
 import  "github.com/caigecaige/swiftpass-client"


req := request.PayWeixinNativeRequest{}
req.MchId = "xxxxxxxx"  //商户号
req.TotalFee = "1"  //金额
req.NonceStr = ""
req.OutTradeNo = "" //订单号
req.Body = ""
req.MchCreateIp = "127.0.0.1"
req.NonceStr = ""
req.NotifyUrl = ""
client := swiftpass.NewDefaultClient()
client.Key = "xxxxxxxx" //商户支付秘钥
resp := new(response.PayWeixinNativeResponse)
client.Execute(req, resp)
```