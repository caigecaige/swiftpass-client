package swiftpass

import (
	"caige/swiftpass/request"
	"caige/swiftpass/response"
	"caige/utils"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
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

const (
	SignTypeMd5 = "md5"
	SignTypeRsa = "rsa"
)

type SwiftClient struct {
	Key         string
	GetWay      string
	SignType    string
	RequestXml  string
	ResponseXml string
	PrivateKey  []byte
	PublicKey   []byte
}

func NewDefaultClient() SwiftClient {
	var client SwiftClient
	client.GetWay = SWIFTPASS_GETWAY
	client.SignType = SignTypeMd5
	return client
}

//加载rsa秘钥文件
func (swiftClient *SwiftClient) LoadRsaKeyFile(priveteKey, publicKey string) {
	readPrivateByte, readErr := ioutil.ReadFile(priveteKey)
	if readErr != nil || len(readPrivateByte) == 0 {
		panic("private key file read fail")
	}
	swiftClient.PrivateKey = readPrivateByte
	readPublicByte, readErr := ioutil.ReadFile(publicKey)
	if readErr != nil || len(readPublicByte) == 0 {
		panic("private key file read fail")
	}
	swiftClient.PublicKey = readPublicByte
}

//加载rsa秘钥文件
func (swiftClient *SwiftClient) LoadRsaKeyByte(priveteKey, publicKey []byte) {
	swiftClient.PrivateKey = priveteKey
	swiftClient.PublicKey = publicKey
}

//加载秘钥文件
func loadPemByte(fileContent []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(fileContent)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	return key
}

//生成md5签名
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

//生成ras签名
func (swiftClient *SwiftClient) GenerateSignRsa(req request.Request) string {
	if len(swiftClient.PrivateKey) == 0 || len(swiftClient.PublicKey) == 0 {
		panic("please load key file!")
	}
	data := []byte(toUrl(req))
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	privateKey := loadPemByte(swiftClient.PrivateKey)
	newData, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(newData)
}

//提交支付
func (swiftClient *SwiftClient) Execute(swpReq request.Request, swpResp response.Response) (bool, error) {
	swiftClient.GetWay = SWIFTPASS_GETWAY
	sign := ""
	if swiftClient.SignType == "" || swiftClient.SignType == SignTypeMd5 {
		sign = swiftClient.GenerateSign(swpReq)
	} else {
		sign = swiftClient.GenerateSignRsa(swpReq)
	}
	data := swpReq.DecodeToXml(sign)
	swiftClient.RequestXml = string(data)
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
	swiftClient.ResponseXml = string(content)
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

//生成url串
func toUrl(req request.Request) string {
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
	paramStr = strings.TrimRight(paramStr, "&")
	return paramStr
}
