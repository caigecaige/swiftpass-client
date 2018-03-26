package request

import "encoding/xml"

type PayAlipayJsPayRequest struct {
	XMLName xml.Name `xml:"xml"`
	CommonParam
	DeviceInfo           string `xml:"device_info"`
	Body                 string `xml:"body"`
	Attach               string `xml:"attach"`
	TotalFee             string `xml:"total_fee"`
	MchCreateIp          string `xml:"mch_create_ip"`
	NotifyUrl            string `xml:"notify_url"`
	LimitCreditPay       string `xml:"limit_credit_pay"`
	TimeStart            string `xml:"time_start"`
	TimeExpire           string `xml:"time_expire"`
	QrCodeTimeoutExpress string `xml:"qr_code_timeout_express"`
	OpUserId             string `xml:"op_user_id"`
	GoodsTag             string `xml:"goods_tag"`
	ProductId            string `xml:"product_id"`
	BuyerLogonId         string `xml:"buyer_logon_id"`
	BuyerId              string `xml:"buyer_id"`
}

func (pwnr PayAlipayJsPayRequest) ServiceName() string {
	return "pay.alipay.jspay"
}

func (pwnr PayAlipayJsPayRequest) DecodeToXml(sign string) []byte {
	pwnr.Sign = sign
	if pwnr.SignType == "" {
		pwnr.SignType = "MD5"
	}
	pwnr.Service = pwnr.ServiceName()
	xmlByte, decodeError := xml.Marshal(pwnr)
	if decodeError == nil {
		return xmlByte
	} else {
		return nil
	}
}
