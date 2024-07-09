package service

import (
	"fmt"
	"log"
	"log/slog"
)

// 1.登录 获取accsess_token 并存储，增加定时事件，过期前60s刷新

// 2.创建wss连接

func InitConfig() {
	slog.SetLogLoggerLevel(slog.LevelInfo)
	if err := ReadConfig("./config/config.xml"); err != nil {
		log.Fatal("InitConfig failed")
	}

	slog.Info("read config success ")
	if err := InitAccessToken(); err != nil {
		log.Fatal("init access token err")
	}
	slog.Info("get access token success")
	var err error
	GAccessCfg.WssURL, err = GetWssUrl(RemoteSrv.SandBoxURL+"/gateway", RemoteSrv.AppId, GAccessCfg.AccessToken.Token)
	if err != nil {
		log.Fatal("init access token err")
	}
	slog.Info("get wss url success")

	GSparkClient = *NewSparkClient(
		RemoteSrv.Spark.SparkAIAppID,
		RemoteSrv.Spark.SparkAIAPIKey,
		RemoteSrv.Spark.SparkAIAPISecret,
		"https://spark-api-open.xf-yun.com/v1/chat/completions",
	)
}

func StartEngine() {
	InitConfig()
	go RefreshAccessToken()

	fmt.Println("GAccessCfg.WssURL=============")
	fmt.Println(GAccessCfg.WssURL)
	WssConn(GAccessCfg.WssURL)
	// time.Sleep(time.Second * 60)

}
