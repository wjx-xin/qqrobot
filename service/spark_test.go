package service

import (
	"fmt"
	"testing"
)

func TestSpark(t *testing.T) {
	client := NewSparkClient("", "", "", "")
	res, err := client.Infer("你好，今晚夜色很美")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(res)
}
