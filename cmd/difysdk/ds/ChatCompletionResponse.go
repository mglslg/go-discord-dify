package ds

type ChatCompletionResponse struct {
	Event          string `json:"event"`
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	Mode           string `json:"mode"`
	Answer         string `json:"answer"`
	Metadata       struct {
		Usage struct {
			PromptTokens        int     `json:"prompt_tokens"`
			PromptUnitPrice     string  `json:"prompt_unit_price"`
			PromptPriceUnit     string  `json:"prompt_price_unit"`
			PromptPrice         string  `json:"prompt_price"`
			CompletionTokens    int     `json:"completion_tokens"`
			CompletionUnitPrice string  `json:"completion_unit_price"`
			CompletionPriceUnit string  `json:"completion_price_unit"`
			CompletionPrice     string  `json:"completion_price"`
			TotalTokens         int     `json:"total_tokens"`
			TotalPrice          string  `json:"total_price"`
			Currency            string  `json:"currency"`
			Latency             float64 `json:"latency"`
		} `json:"usage"`
		RetrieverResources []struct {
			Position     int     `json:"position"`
			DatasetID    string  `json:"dataset_id"`
			DatasetName  string  `json:"dataset_name"`
			DocumentID   string  `json:"document_id"`
			DocumentName string  `json:"document_name"`
			SegmentID    string  `json:"segment_id"`
			Score        float64 `json:"score"`
			Content      string  `json:"content"`
		} `json:"retriever_resources"`
	} `json:"metadata"`
	CreatedAt int `json:"created_at"`
}
