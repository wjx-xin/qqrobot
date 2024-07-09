package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"robot/model"
)

// #星火认知大模型调用秘钥信息，请前往讯飞开放平台控制台（https://console.xfyun.cn/services/bm35）查看

type SparkClient struct {
	appId     string
	apiKey    string
	apiSecret string
	url       string
	client    http.Client
}

func NewSparkClient(appId, apiKey, apiSecret, url string) *SparkClient {
	return &SparkClient{
		appId:     appId,
		apiKey:    apiKey,
		apiSecret: apiSecret,
		url:       url,
		client:    http.Client{},
	}
}

func (spark *SparkClient) Infer(msg string) (string, error) {
	data := map[string]interface{}{
		"model":    "generalv3.5",
		"messages": []map[string]string{{"role": "user", "content": "假装你是一个叫作leaf的智者，需要以简短的话语解答疑问 。问：\n" + msg}},
	}
	header := map[string]string{
		"Authorization": "Bearer " + spark.apiKey + ":" + spark.apiSecret, // 注意此处替换自己的key和secret
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Error marshaling data:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", spark.url, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating request:", err)
		return "", err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}
	resp, err := spark.client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body:", err)
		return "", err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(body, &m); err != nil {
		return "", err
	}

	// slog.Error(m)
	var respMsg model.SparkResp

	// var m map[string]interface{}
	if err := json.Unmarshal(body, &respMsg); err != nil {
		return "", err
	}
	res := respMsg.Choices[0].Message.Content
	return res, nil
}
