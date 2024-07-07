package service

import (
	"fmt"
	"log"
)

// 1.登录 获取accsess_token 并存储，增加定时事件，过期前60s刷新

// 2.创建wss连接

func InitConfig() {
	if err := ReadConfig("./config/config.xml"); err != nil {
		log.Fatal("InitConfig failed")
	}
	fmt.Println("===========read config success ===========")
	if err := InitAccessToken(); err != nil {
		log.Fatal("init access token err")
	}
	fmt.Println("==========get access token success============")
	var err error
	GAccessCfg.WssURL, err = GetWssUrl(QQService.SandBoxURL+"/gateway", QQService.AppId, GAccessCfg.AccessToken.Token)
	if err != nil {
		log.Fatal("init access token err")
	}
	fmt.Println("====================get wss url success")
	fmt.Println(QQService)
	fmt.Println(GAccessCfg)
}

func StartEngine() {
	InitConfig()
	go RefreshAccessToken()

	fmt.Println("GAccessCfg.WssURL=============")
	fmt.Println(GAccessCfg.WssURL)
	WssConn(GAccessCfg.WssURL)
	// time.Sleep(time.Second * 60)

}
