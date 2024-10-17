package difysdk

import (
	"bytes"
	"encoding/json"
	ds "github.com/mglslg/go-discord-dify/cmd/difysdk/ds"
	"github.com/mglslg/go-discord-dify/cmd/g"
	"io"
	"net/http"
)

func Chat(msg string, userName string, conversationId string) (string, string, error) {
	url := "https://dify.hogwartscoder.com/v1/chat-messages"

	chatRequestBody := ds.ChatRequestBody{
		Query:            msg,
		ResponseMode:     "blocking",
		User:             userName,
		ConversationID:   conversationId,
		AutoGenerateName: true,
		Inputs: map[string]interface{}{
			"none": "none",
		},
	}

	body, err := json.Marshal(chatRequestBody)

	g.Logger.Println("request dify :", body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.Println("Error creating request:", err)
		return "[Error creating request:" + err.Error() + "]", "", nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.SecToken.Dify)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.Logger.Println("Error sending request", err)
		return "[Error sending request:" + err.Error() + "]", "", nil
	}

	if resp.StatusCode != 200 {
		return resp.Status, "", nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			g.Logger.Println("Error closing response body", err)
		}
	}(resp.Body)

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		g.Logger.Println("Error reading response", err)
		return "[Error reading response:" + err.Error() + "]", "", nil
	}

	chatResponse := ds.ChatCompletionResponse{}
	err = json.Unmarshal(responseBody, &chatResponse)
	if err != nil {
		g.Logger.Println("[Error unmarshalling response]", err)
		return "[Error unmarshalling response:" + err.Error() + "]", "", nil
	}
	g.Logger.Println(">>>>>chat response:", chatResponse.Answer)

	return chatResponse.Answer, chatResponse.ConversationID, nil
}

func DeleteConversation(conversationId string, userName string) (string, error) {
	url := "https://dify.hogwartscoder.com/v1/conversations/" + conversationId
	params := map[string]interface{}{
		"user": userName,
	}
	body, err := json.Marshal(params)

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.Println("Error creating request:", err)
		return "[Error creating request:" + err.Error() + "]", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.SecToken.Dify)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		g.Logger.Println("Error sending request", err)
		return "[Error sending request:" + err.Error() + "]", err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			g.Logger.Println("Error closing response body", err)
		}
	}(resp.Body)

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		g.Logger.Println("Error reading response", err)
		return "[Error reading response:" + err.Error() + "]", err
	}
	response := ds.CommonResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		g.Logger.Println("[Error unmarshalling response]", err)
		return "[Error unmarshalling response:" + err.Error() + "]", err
	}
	g.Logger.Println(">>>>>delete response:", response)

	return response.Result, nil
}
