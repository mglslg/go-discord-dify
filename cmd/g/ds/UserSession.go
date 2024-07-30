package ds

// UserSession Each user holds a session in each channel to save its state
type UserSession struct {
	UserChannelID  string //UserSession's unique key
	UserId         string
	ChannelID      string
	UserName       string
	ConversationID string
	ChatCount      int //Number of messages sent by the user
}
