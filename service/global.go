package service

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type GAccessConfig struct {
	AccessToken AccessToken
	WssURL      string
}

type QQServer struct {
	SandBoxURL   string `xml:"sandbox"`
	AuthUrl      string `xml:"auth_url"`
	AppId        string `xml:"appid"`
	ClientSecret string `xml:"client_secret"`
}

var QQService QQServer

var GAccessCfg GAccessConfig

func ReadConfig(path string) error {
	xmlFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer xmlFile.Close()

	// 读取文件内容
	data, err := io.ReadAll(xmlFile)
	if err != nil {
		return err
	}

	// 解析XML到Server结构体
	err = xml.Unmarshal(data, &QQService)
	if err != nil {
		return err
	}
	return nil
}

func InitAccessToken() error {
	var err error
	GAccessCfg.AccessToken, err = GetAccessToken(QQService.AuthUrl, QQService.AppId, QQService.ClientSecret)
	if err != nil {
		return err
	}
	return nil
}

func RefreshAccessToken() {
	// var err error
	for {
		duration, err := strconv.Atoi(GAccessCfg.AccessToken.ExpiresIn)
		if err != nil {
			return
		}
		ticker := time.NewTicker(time.Second * time.Duration(duration-50))

		<-ticker.C
		GAccessCfg.AccessToken, err = GetAccessToken(QQService.AuthUrl, QQService.AppId, QQService.ClientSecret)
		if err != nil {
			fmt.Println(err.Error())
			// return
		}

	}
}
