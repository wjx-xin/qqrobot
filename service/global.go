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

type Server struct {
	SandBoxURL   string `xml:"sandbox"`
	AuthUrl      string `xml:"auth_url"`
	AppId        string `xml:"appid"`
	ClientSecret string `xml:"client_secret"`
	Spark        Spark  `xml:"spark"`
}

type Spark struct {
	SparkAIAppID     string `xml:"SPARKAI_APP_ID"`
	SparkAIAPISecret string `xml:"SPARKAI_API_SECRET"`
	SparkAIAPIKey    string `xml:"SPARKAI_API_KEY"`
}

var RemoteSrv Server

var GAccessCfg GAccessConfig

var GSparkClient SparkClient

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
	err = xml.Unmarshal(data, &RemoteSrv)
	if err != nil {
		return err
	}
	return nil
}

func InitAccessToken() error {
	var err error
	GAccessCfg.AccessToken, err = GetAccessToken(RemoteSrv.AuthUrl, RemoteSrv.AppId, RemoteSrv.ClientSecret)
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
		GAccessCfg.AccessToken, err = GetAccessToken(RemoteSrv.AuthUrl, RemoteSrv.AppId, RemoteSrv.ClientSecret)
		if err != nil {
			fmt.Println(err.Error())
			// return
		}

	}
}
