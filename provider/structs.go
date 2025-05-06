package provider

// Define the Role type
type Role string

// Define constants for the Role type
const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// GET - /v1/models
type GetModelsResponse struct {
	Object string `json:"object"`
	Data   []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

// POST - /check-balance
type PostCheckBalanceResponse struct {
	UsdBalance         string `json:"usd_balance"`
	NanoBalance        string `json:"nano_balance"`
	NanoDepositAddress string `json:"nanoDepositAddress"`
}

// POST - /v1/chat/completions
type PostChatCompletionsRequest struct {
	Model            string                              `json:"model"`
	Messages         []PostChatCompletionsRequestMessage `json:"messages"`
	Stream           bool                                `json:"stream"`
	Temperature      float64                             `json:"temperature"`
	MaxTokens        int                                 `json:"max_tokens"`
	TopP             int                                 `json:"top_p"`
	FrequencyPenalty int                                 `json:"frequency_penalty"`
	PresencePenalty  int                                 `json:"presence_penalty"`
}

type PostChatCompletionsRequestMessage struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

// POST - /v1/chat/completions
type PostChatCompletionsResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
