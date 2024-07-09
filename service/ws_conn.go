package service

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"robot/model"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type WssAuthReq struct {
	Op int            `json:"op"`
	D  wssAuthReqData `json:"d"`
}

type WssHeartReq struct {
	Op int  `json:"op"`
	D  *int `json:"d"`
}

var MsgS int64

type wssAuthReqData struct {
	Token      string            `json:"token"`
	Intents    int               `json:"intents"`
	Shard      []int             `json:"shard"`
	Properties map[string]string `json:"properties"`
}

func WssAuth(conn *websocket.Conn) {
	data := WssAuthReq{
		Op: 2,
		D: wssAuthReqData{
			Token:      "QQBot " + GAccessCfg.AccessToken.Token,
			Intents:    0 | (1 << 30) | (1 << 25),
			Shard:      []int{0, 1},
			Properties: nil,
		},
	}
	// slog.Error(json.Marshal(&data))
	err := conn.WriteJSON(&data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func Heartbeat(conn *websocket.Conn) {
	atomic.LoadInt64(&MsgS)
	var d int = int(MsgS)
	data := WssHeartReq{
		Op: 1,
		D:  &d,
	}
	err := conn.WriteJSON(&data)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func WssConn(url string) {
	// url = "wss://sandbox.api.sgroup.qq.com/websocket"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("wss: Error creating request:", err)
		return
	}

	// 设置请求头，例如设置内容类型为JSON
	req.Header.Set("Content-Type", "application/json")

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			// log.Printf("recv1:%d ==== %s", msgType, message)

			var msgData model.MessageData
			if err := json.Unmarshal(message, &msgData); err != nil {
				slog.Error(err.Error())
			}
			atomic.StoreInt64(&MsgS, int64(msgData.Sequence))
			if msgData.Type == "AT_MESSAGE_CREATE" {
				Reply(&msgData)
			}
		}
	}()

	WssAuth(conn)
	Heartbeat(conn)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			Heartbeat(conn)
		case <-done:
			return

		case <-interrupt:
			log.Println("interrupt")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func GetWssUrl(url string, appId string, accessToken string) (string, error) {

	// jsonData := fmt.Sprintf(`{"access_token": "MttEd34YdL1HIIy4NhxKabO-xHaWcnkMHCtHgvBfoB82jAWA4_etAhBpqBYXcYF02nSt2ehdn7Ix-w"}`
	// body := bytes.NewBufferString(jsonData)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		slog.Error("Error creating request:", err)
		return "", err
	}

	// 设置请求头，例如设置内容类型为JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "QQBot "+accessToken)
	req.Header.Set("X-Union-Appid", appId)

	// 使用自定义的Transport创建HTTP客户端
	client := &http.Client{}

	// 发送请求并获取响应
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response body:", err)
		return "", err
	}

	// 打印响应体内容
	var urlResp struct {
		Url string
	}

	if err := json.Unmarshal(bodyBytes, &urlResp); err != nil {
		return "", err
	}
	slog.Info("Response:", string(bodyBytes))
	return urlResp.Url, nil
}
