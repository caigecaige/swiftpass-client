//威富通的支付接口--caige
package request

type Request interface {
	ServiceName() string
	DecodeToXml(string) []byte
}

type CommonParam struct {
	Service    string `xml:"service"`
	Version    string `xml:"version"`
	Charset    string `xml:"charset"`
	SignType   string `xml:"sign_type"`
	MchId      string `xml:"mch_id"`
	OutTradeNo string `xml:"out_trade_no"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
}
