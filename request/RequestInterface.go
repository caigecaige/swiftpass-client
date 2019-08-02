//威富通的支付接口--caige
package request

type Request interface {
	ServiceName() string
	DecodeToXml(string) []byte
}

type CommonParam struct {
	Service    string `xml:"service,omitempty"`
	Version    string `xml:"version,omitempty"`
	Charset    string `xml:"charset,omitempty"`
	SignType   string `xml:"sign_type,omitempty"`
	MchId      string `xml:"mch_id,omitempty"`
	OutTradeNo string `xml:"out_trade_no,omitempty"`
	NonceStr   string `xml:"nonce_str,omitempty"`
	Sign       string `xml:"sign,omitempty"`
}
