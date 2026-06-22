package llm

type ChatRequest struct {
	SystemPrompt string
	Prompt       string
	Messages     []Message
	Temperature  float32
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Content string `json:"content"`
	Model   string `json:"model,omitempty"`
}
