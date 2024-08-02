package ds

type AppContext struct {
	LogFilePath    string `yaml:"logFilePath"`
	ApplicationId  string `yaml:"applicationId"`
	GuildId        string `yaml:"guildId"`
	BotName        string `yaml:"botName"`
	BotToken       string `yaml:"botToken"`
	DifyToken      string `yaml:"difyToken"`
	ClearCmd       string `yaml:"clearCmd"`
	ClearCmdDesc   string `yaml:"clearCmdDesc"`
	ClearDelimiter string `yaml:"clearDelimiter"`
	FreeChatLimit  int    `yaml:"creeChatLimit"`
	OnAt           bool   `yaml:"onAt"`
	BotId          string
	ConfigFilePath string
}
