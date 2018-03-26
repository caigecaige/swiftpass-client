package request

import "encoding/xml"

type PayWeixinNativeRequest struct {
	XMLName xml.Name `xml:"xml"`
	CommonParam
	DeviceInfo  string `xml:"device_info"`
	Body        string `xml:"body"`
	Attach      string `xml:"attach"`
	TotalFee    string `xml:"total_fee"`
	NotifyUrl   string `xml:"notify_url"`
	TimeStart   string `xml:"time_start"`
	TimeExpire  string `xml:"time_expire"`
	OpUserId    string `xml:"op_user_id"`
	GoodsTag    string `xml:"goods_tag"`
	ProductId   string `xml:"product_id"`
	MchCreateIp string `xml:"mch_create_ip"`
}

func (pwnr PayWeixinNativeRequest) ServiceName() string {
	return "pay.weixin.native"
}

func (pwnr PayWeixinNativeRequest) DecodeToXml(sign string) []byte {
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
