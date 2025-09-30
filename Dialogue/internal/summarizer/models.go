package summarizer

// DialogueSummaryResponse 定义对话总结响应结构
type DialogueSummaryResponse struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    string   `json:"data,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// DialogueRequest 定义对话请求结构
type DialogueRequest struct {
	Content string `json:"content"`
}