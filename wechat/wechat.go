package wechat

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bitly/go-simplejson"
)

const (
	AccTokenUrl       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
	SendmsgUrl        = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"
	TokeExpSec  int64 = 7200
)

type extend struct {
	AccToken string
	TokenTS  int64
}

type WeChat struct {
	CorpID  string
	AgentId int
	Secret  string
	*extend
}

var TLSClient *http.Client

func init() {
	/* file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)

	   if err != nil {
	       log.Fatalln("LogFile ioError:", err)
	   } */
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM([]byte(wxRootPEM))
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}

	TLSClient = &http.Client{Transport: tr}

}

func New(CorpID string, AgentId int, Secret string) *WeChat {
	return &WeChat{
		CorpID:  CorpID,
		AgentId: AgentId,
		Secret:  Secret,
		extend:  &extend{},
	}
}

type Content struct {
	Content string `json:"content"`
}

type JsonMsg struct {
	ToUser  string  `json:"touser,omitempty"`
	ToParty string  `json:"toparty,omitempty"`
	MsgType string  `json:"msgtype"`
	AgentID int     `json:"agentid"`
	Text    Content `json:"text"`
}

/*func checkErr(err error){

}*/

func (wx *WeChat) GetAccToken() error {
	getAccTokenUrl := fmt.Sprintf(AccTokenUrl, wx.CorpID, wx.Secret)

	rsp, err := TLSClient.Get(getAccTokenUrl)
	if err != nil {
		return err
	}
	json, err := simplejson.NewFromReader(rsp.Body)
	if err != nil {
		return err
	}
	errcode := json.Get("errcode").MustInt(1)
	if errcode != 0 {
		return fmt.Errorf("get WeChat Access Token error:", json.Get("errmsg").MustString(""))

	}
	wx.AccToken = json.Get("access_token").MustString("")
	wx.TokenTS = time.Now().Unix()
	return fmt.Errorf("getAccToken done: %s", wx.AccToken)

}

func (wx WeChat) SendMsg(touser, toparty, content string) (string, error) {

	msg := JsonMsg{
		ToUser:  touser,
		ToParty: toparty,
		MsgType: "text",
		AgentID: wx.AgentId,
		Text: Content{
			Content: content,
		},
	}
	for i := 0; i < 3; i++ {
		if wx.AccToken == "" || wx.TokenTS-time.Now().Unix() <= -TokeExpSec {
			wx.GetAccToken()
		} else {
			jmsg, err := json.Marshal(msg)
			if err != nil {
				return "", err
			}

			postSendmsgUrl := fmt.Sprintf(SendmsgUrl, wx.AccToken)
			rsp, err := TLSClient.Post(postSendmsgUrl, "application/json;charset=utf-8", bytes.NewReader(jmsg))
			if err != nil {
				return "", err
			}
			byteData, err := ioutil.ReadAll(rsp.Body)
			return string(byteData[:]), err
		}
	}

	return "", fmt.Errorf("getAccToken failed: %s", wx.GetAccToken().Error())
}
