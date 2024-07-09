package model

type SparkResp struct {
	Code      int      `json:"code"`
	Message   string   `json:"message"`
	SessionID string   `json:"sid"` // 通常 sid 指的是 session ID
	Choices   []Choice `json:"choices"`
	Usage     Usage    `json:"usage"`
}

// 定义用于解析choices数组中每个元素的结构体
type Choice struct {
	Message Message `json:"message"`
	Index   int     `json:"index"`
}

// 定义用于解析message内嵌对象的结构体
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// 定义用于解析usage对象的结构体
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
