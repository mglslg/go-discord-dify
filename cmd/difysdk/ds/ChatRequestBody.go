package ds

type ChatRequestBody struct {
	Inputs           map[string]interface{} `json:"inputs"`
	Query            string                 `json:"query"`
	ResponseMode     string                 `json:"response_mode"`
	User             string                 `json:"user"`
	ConversationID   string                 `json:"conversation_id"`
	AutoGenerateName bool                   `json:"auto_generate_name"`
	Files            []struct {
		Type           string `json:"type"`
		TransferMethod string `json:"transfer_method"`
		Url            string `json:"url"`
	} `json:"files"`
}
