package service

import (
	"fmt"
	"testing"
)

func TestSpark(t *testing.T) {
	resp, err := GetSparkResp("说一个成语")
	if err != nil {
		fmt.Println("============ err in test spark ===============")
	}
	fmt.Println(resp)

}
