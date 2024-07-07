package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn string `json:"expires_in"`
}

func GetAccessToken(url string, appId string, clientSecret string) (AccessToken, error) {
	var token AccessToken

	jsonData := fmt.Sprintf(`{"appId": "%s", "clientSecret": "%s"}`, appId, clientSecret)
	body := bytes.NewBufferString(jsonData)

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return token, err
	}
	req.Header.Set("Content-Type", "application/json")

	// // 如果需要跳过TLS证书验证，可以设置Transport的TLSClientConfig
	tr := &http.Transport{
		// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 使用自定义的Transport创建HTTP客户端
	client := &http.Client{Transport: tr}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return token, err
	}
	defer resp.Body.Close()

	// 读取响应体内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return token, err
	}
	fmt.Println("Response:", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &token); err != nil {
		return token, err
	}
	// 打印响应体内容

	return token, nil
}
