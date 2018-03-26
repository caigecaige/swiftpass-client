package swiftpass

import (
	"caige/swiftpass/request"
	"caige/swiftpass/response"
	"caige/utils"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

const (
	SWIFTPASS_GETWAY = "https://pay.swiftpass.cn/pay/gateway"
)

const (
	RESULT_STATUS_SUCC = "0"
	RESULT_CODE_SUCC   = "0"
)

const (
	TRADE_STATE_SUCCESS = "SUCCESS"
	TRADE_STATE_REFUND  = "REFUND"
	TRADE_STATE_NOTPAY  = "NOTPAY"
	TRADE_STATE_CLOSED  = "CLOSED"
	TRADE_STATE_REVERSE = "REVERSE"
	TRADE_STATE_REVOKED = "REVOKED"
)

type SwiftClient struct {
	Key    string
	GetWay string
	Sign   string
}

func NewDefaultClient() SwiftClient {
	var client SwiftClient
	client.GetWay = SWIFTPASS_GETWAY
	return client
}

func (swiftClient *SwiftClient) GenerateSign(req request.Request) string {
	xmlData := req.DecodeToXml("")
	signParams := make(map[string]string)
	signParams = utils.Xml2Map(strings.NewReader(string(xmlData)))
	for k, v := range signParams {
		if v == "" {
			delete(signParams, k)
		}
	}
	sortParams := sortParams(signParams)
	paramStr := ""
	for _, v := range sortParams {
		for ck, cv := range v {
			paramStr = paramStr + ck + "=" + cv + "&"
		}
	}
	paramStr = paramStr + "key=" + swiftClient.Key
	sign := strings.ToUpper(utils.Md5String(paramStr))
	return sign
}

func (swiftClient *SwiftClient) Execute(swpReq request.Request, swpResp response.Response) (bool, error) {
	swiftClient.GetWay = SWIFTPASS_GETWAY
	sign := swiftClient.GenerateSign(swpReq)
	data := swpReq.DecodeToXml(sign)
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("POST", swiftClient.GetWay, strings.NewReader(string(data)))
	req.Header.Set("Accept", "text/xml,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Charset", "UTF-8")
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()
	content, rerr := ioutil.ReadAll(res.Body)
	if rerr != nil {
		return false, rerr
	}
	encodeErr := xml.Unmarshal(content, swpResp)
	if encodeErr != nil {
		return false, encodeErr
	}
	return true, nil
}

//排序
func sortParams(m map[string]string) []map[string]string {
	nm := make([]map[string]string, 0)
	pm := make([]string, 0)
	for k, _ := range m {
		pm = append(pm, k)
	}
	sort.Strings(pm)
	for _, sv := range pm {
		nm = append(nm, map[string]string{sv: m[sv]})
	}
	return nm
}
