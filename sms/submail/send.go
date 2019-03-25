package submail

import (
	"fmt"
	"net/url"
)

type SmsSend struct {
	appId    string
	appKey   string
	signType string

	to      string
	content string
	tag     string
}

func NewSmsSend(appid int, appkey, signtype string) *SmsSend {
	return &SmsSend{
		appId:    fmt.Sprintf("%d", appid),
		appKey:   appkey,
		signType: signtype,
	}
}

func (this *SmsSend) SetTo(to string) {
	this.to = to
}

func (this *SmsSend) SetContent(content string) {
	this.content = content
}

func (this *SmsSend) SetTag(tag string) {
	this.tag = tag
}

func (this *SmsSend) Send() (string, error) {
	config := make(map[string]string)

	config["appid"] = this.appId
	config["appkey"] = this.appKey
	config["signType"] = this.signType

	values := url.Values{}

	values.Set("appid", this.appId)
	values.Set("to", this.to)

	if this.signType != "normal" {
		timestamp, err := GetTimestamp()

		if err != nil {
			return "", err
		}

		values.Set("sign_type", this.signType)
		values.Set("timestamp", fmt.Sprintf("%d", timestamp))
		values.Set("sign_version", "2")
	}

	if this.tag != "" {
		values.Set("tag", this.tag)
	}

	signature := caculSign(values, config)

	values.Set("signature", signature)

	//v2 数字签名 content 不参与计算
	values.Set("content", this.content)

	return HttpPost(SUBMAIL_SMS_SEND_URL, values.Encode())
}
