package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"robot/model"
	"strings"
)

func ExtractContent(msg *model.MessageData) string {
	strings.Index(msg.Data.Content, "")
	return ""
}

func Reply(msg *model.MessageData) {

	// 替换为你的子频道ID
	channelId := msg.Data.ChannelID
	fmt.Println("=============check==============")
	fmt.Println(msg.Data.Mentions)
	content := strings.ReplaceAll(msg.Data.Content, "\u003c@!"+msg.Data.Mentions[0].ID+"\u003e", "@leaf")
	// 替换为你的机器人Token
	// token := QQService.AppId + "." + GAccessCfg.AccessToken.Token
	fmt.Println(content)
	token := GAccessCfg.AccessToken.Token

	// 构造请求体
	// var inferData = ""
	inferData, err := GetSparkResp(content)

	if err != nil {
		fmt.Println("============ err in spark==========")
	}
	requestBody := map[string]interface{}{
		"content": fmt.Sprintf("<@!%s>%s", msg.Data.Author.ID, inferData), // 替换USER_ID为实际用户ID
		"msg_id":  msg.ID,
	}
	requestBodyJson, _ := json.Marshal(requestBody)

	// 创建HTTP请求
	url := QQService.SandBoxURL + "/channels/" + channelId + "/messages"
	// url := https://api.q.qq.com/v2/channels/
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	fmt.Println("===================reply=====================")
	fmt.Println(url)
	fmt.Println(requestBody)
	// 设置请求头
	req.Header.Add("Authorization", "QQBot "+token)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println(req.Header)
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// 处理响应
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("======================reply err============================")
	} else {
		fmt.Println("=====================reply success========")
		fmt.Println("Response:", string(responseBody))
	}

}