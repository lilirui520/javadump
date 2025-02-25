package notify

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/braumye/grobot"
	"jvmdump4k8s/config"
	"net/http"
	"net/url"
	"time"
)

// https://github.com/braumye/grobot
func SendDingtalk(fileurl string) {
	token := config.GlobalConfig.NotifyDingToken
	podName := config.GlobalConfig.PodName
	fmt.Printf("开始推送钉钉 token %s %s \n", token, fileurl)
	robot, _ := grobot.New("dingtalk", token)
	// 发送文本消息
	err := robot.SendTextMessage(fmt.Sprintf("报警 %s应用发生OOM , dump文件%s ", podName, fileurl))
	fmt.Println("推送钉钉完成 err=", err)
}

type Md struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

type DingTalkTextMessage struct {
	MsgType  string `json:"msgtype"`
	Markdown Md     `json:"markdown"`
}

// 计算签名
func calculateSign(secret string, timestamp int64) string {
	// 拼接字符串
	strToSign := fmt.Sprintf("%d\n%s", timestamp, secret)

	// 使用 HMAC-SHA256 加密
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(strToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func SendDingTalkWithSign(fileurl string, podName string) {
	accessToken := config.GlobalConfig.NotifyDingToken
	secret := config.GlobalConfig.DingTokenSign

	// 获取当前时间戳（单位：毫秒）
	timestamp := time.Now().UnixNano() / 1e6

	// 计算签名
	sign := calculateSign(secret, timestamp)

	// 构建完整的 Webhook URL
	webhookURL := fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s&timestamp=%d&sign=%s",
		accessToken, timestamp, url.QueryEscape(sign))

	fmt.Println(fileurl, webhookURL)

	// 创建消息
	message := DingTalkTextMessage{
		MsgType: "markdown",
		Markdown: Md{
			Title: "JAVA OOM告警",
			Text:  fmt.Sprintf("### ⚠️ 风险预警 ⚠️ \n\n **线上环境**: %s应用发生OOM \n\n  **下载地址**：[点击下载文件](http://%s)", podName, fileurl),
		},
	}

	// 将消息转换为 JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		fmt.Println("JSON 编码失败:", err)
		return
	}

	// 发送 HTTP POST 请求
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode == http.StatusOK {
		fmt.Println("消息发送成功！")
	} else {
		fmt.Println("消息发送失败，状态码:", resp.StatusCode)
	}
}
