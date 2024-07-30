package difysdk

import (
	"bytes"
	"encoding/json"
	ds "github.com/mglslg/go-discord-gpt/cmd/difysdk/ds"
	"github.com/mglslg/go-discord-gpt/cmd/g"
	"io/ioutil"
	"net/http"
)

func Chat(msg string, userName string, conversationId string) (string, string, error) {
	api := "http://dify.hogwartscoder.com/v1/chat-messages"

	chatRequestBody := ds.ChatRequestBody{
		Query:            msg,
		ResponseMode:     "blocking",
		User:             userName,
		ConversationID:   conversationId,
		AutoGenerateName: true,
	}

	body, err := json.Marshal(chatRequestBody)

	req, err := http.NewRequest("POST", api, bytes.NewBuffer(body))
	if err != nil {
		g.Logger.Println("Error creating request:", err)
		return "[Error creating request:" + err.Error() + "]", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+g.SecToken.Dify)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		g.Logger.Println("Error sending request", err)
		return "[Error sending request:" + err.Error() + "]", "", err
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		g.Logger.Println("Error reading response", err)
		return "[Error reading response:" + err.Error() + "]", "", err
	}

	chatResponse := ds.ChatCompletionResponse{}
	err = json.Unmarshal(body, &chatResponse)
	if err != nil {
		g.Logger.Println("[Error unmarshalling response]", err)
		return "[Error unmarshalling response:" + err.Error() + "]", "", err
	}
	g.Logger.Println(">>>>>dify response:", chatResponse.Answer)

	return chatResponse.Answer, chatResponse.ConversationID, nil
}
