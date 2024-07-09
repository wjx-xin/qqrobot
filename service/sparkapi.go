package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
)

// #星火认知大模型调用秘钥信息，请前往讯飞开放平台控制台（https://console.xfyun.cn/services/bm35）查看
var (
	SPARKAI_APP_ID     = ""
	SPARKAI_API_SECRET = ""
	SPARKAI_API_KEY    = ""
)

func GetSparkResp(msg string) (string, error) {
	url := "https://spark-api-open.xf-yun.com/v1/chat/completions"
	data := map[string]interface{}{
		"model":    "generalv3.5",
		"messages": []map[string]string{{"role": "user", "content": "1.假装你是一个叫作leaf的智者；2.回复控制在50字以内 。问：\n" + msg}},
	}
	header := map[string]string{
		"Authorization": "Bearer " + RemoteSrv.Spark.SparkAIAPIKey + ":" + RemoteSrv.Spark.SparkAIAPISecret, // 注意此处替换自己的key和secret
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		slog.Error("Error marshaling data:", err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		slog.Error("Error creating request:", err)
		return "", err
	}

	for k, v := range header {
		req.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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
	choices, ok := m["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", errors.New("err in spark reply")
	}
	choicesZero := choices[0].(map[string]interface{})
	message := choicesZero["message"].(map[string]interface{})
	res := message["content"].(string)
	return res, nil
}
