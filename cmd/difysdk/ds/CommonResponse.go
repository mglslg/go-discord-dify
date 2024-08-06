package ds

type CommonResponse struct {
	Result string      `json:"result"`
	Data   interface{} `json:"data"`
}
