package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/go-discord-dify/cmd/g"
	"github.com/mglslg/go-discord-dify/cmd/g/ds"
	"regexp"
	"time"
)

var (
	ctx = g.AppContext
)

func initDiscordSession() (*discordgo.Session, error) {
	session, err := discordgo.New("Bot " + g.SecToken.Discord)
	if err != nil {
		logger.Fatal("Error creating Discord session:", err)
		return nil, err
	}

	//设置机器人权限
	//intents := discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent
	intents := discordgo.IntentsAllWithoutPrivileged
	session.Identify.Intents = intents

	logger.Println("开始创建cmd监听")

	createCmd(session, ctx.ClearCmd, ctx.ClearCmdDesc)
	//createCmd(session, "准备魔杖", "自定义Prompt")
	session.AddHandler(onCommand)

	//监听消息
	session.AddHandler(onMsgCreate)

	//删除命令方法
	//session.ApplicationCommandDelete(g.Role.ApplicationId, g.Conf.GuildID, "1103997867103899679")

	return session, nil
}

func createCmd(session *discordgo.Session, cmdName string, cmdDesc string) {
	_, cmdErr := session.ApplicationCommandCreate(ctx.ApplicationId, "", &discordgo.ApplicationCommand{
		Name:        cmdName,
		Description: cmdDesc,
	})
	if cmdErr != nil {
		logger.Fatal("create "+cmdName+" cmd error", cmdErr)
	}
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//如果是机器人发的消息则不予理睬
	if m.Author.ID == g.AppContext.BotId {
		return
	}

	//为当前用户创建session(机器人本身也会有用户session)
	us := g.GetUserSession(m.Author.ID, m.ChannelID, m.Author.Username)

	g.Logger.Println("******************************************************OnMessage******************************************************")
	g.Logger.Println("OnAt:", ctx.OnAt)

	if ctx.OnAt {
		reply(s, m, us)
	} else {
		simpleReply(s, m, us)
	}
}

// Contextual reply without AT
func simpleReply(s *discordgo.Session, m *discordgo.MessageCreate, us *ds.UserSession) {
	//todo
}

// Contextual reply that requires AT
func reply(s *discordgo.Session, m *discordgo.MessageCreate, us *ds.UserSession) {
	if m.Mentions != nil {
		for _, mentioned := range m.Mentions {
			if mentioned.ID == ctx.BotId {
				//todo 使用dify接口实现

				replayContent := apisdk.Chat()

				//allMsg, e := fetchMessagesByCount(s, us.ChannelID, ctx.MaxUserRecord)
				//if e != nil {
				//	logger.Fatal("抓取聊天记录失败", e)
				//}
				//
				////获取聊天上下文
				//conversation := geMentionContext(allMsg, us)
				//
				////异步获取响应结果并提示[正在输入],go关键字后是生产端,asyncResponse中的select是消费端
				//respChannel := make(chan string)
				//go callOpenAIChat(conversation, us, respChannel)
				//asyncResponse(s, m, us, respChannel)

				return
			}
		}
	}
}

// Asynchronous reception of interface response and prompt [typing]
func asyncResponse(s *discordgo.Session, m *discordgo.MessageCreate, us *ds.UserSession, respChannel chan string) {
	for {
		select {
		case gptResp := <-respChannel:
			// Mention the user who asked the question
			msgContent := fmt.Sprintf("%s %s", m.Author.Mention(), gptResp)

			//当消息超长时拆分成两段回复用户,并且不会宕机
			var err error
			if len(msgContent) > 2000 {
				half := len(msgContent) / 2
				firstHalf := msgContent[:half]
				secondHalf := msgContent[half:]
				_, err = s.ChannelMessageSend(us.ChannelID, firstHalf)
				_, err = s.ChannelMessageSend(us.ChannelID, fmt.Sprintf("%s %s", m.Author.Mention(), secondHalf))
			} else {
				_, err = s.ChannelMessageSend(us.ChannelID, msgContent)
			}
			if err != nil {
				logger.Println("发送discord消息失败,当前消息长度:", len(msgContent), err)
				_, err = s.ChannelMessageSend(us.ChannelID, fmt.Sprint("[发送discord消息失败,当前消息长度:", len(msgContent), "]"))
			}
			return

		default:
			err := s.ChannelTyping(us.ChannelID)
			if err != nil {
				return
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func getCleanMsg(content string) string {
	// 创建一个正则表达式，用于匹配尖括号及其内容，格式为：<@数字>
	re := regexp.MustCompile(`<@(\d+)>`)

	// 使用正则表达式替换匹配的内容为空字符串
	cleanedMsg := re.ReplaceAllString(content, "")

	return cleanedMsg
}

func getLatestMessage(messages []*discordgo.Message) *ds.Stack {
	msgStack := ds.NewStack()
	length := len(messages)
	msgStack.Push(messages[length-1])
	return msgStack
}

//func getLatestMentionContext(messages []*discordgo.Message, us *ds.UserSession) *ds.Stack {
//	g.Logger.Println("getLatestMentionContext:delimiter:", us.ClearDelimiter)
//
//	msgStack := ds.NewStack()
//	for _, msg := range messages {
//		for _, mention := range msg.Mentions {
//			//找出当前用户艾特机器人的最后一条记录
//			if msg.Author.ID == us.UserId && mention.ID == ctx.BotId {
//				msgStack.Push(msg)
//				return msgStack
//			}
//		}
//	}
//	return msgStack
//}

//func geMentionContext(messages []*discordgo.Message, us *ds.UserSession) *ds.Stack {
//	g.Logger.Println("geMentionContext:delimiter:", us.ClearDelimiter)
//
//	msgStack := ds.NewStack()
//	for _, msg := range messages {
//		for _, mention := range msg.Mentions {
//			//找出当前用户艾特GPT以及GPT艾特当前用户的聊天记录
//			if (msg.Author.ID == ctx.BotId && mention.ID == us.UserId) || (msg.Author.ID == us.UserId && mention.ID == ctx.BotId) {
//				//一旦发现clear命令的分隔符则直接终止向消息栈push,直接返回
//				if strings.Contains(msg.Content, us.ClearDelimiter) {
//					g.Logger.Println("geMentionContext:delimiter:", us.ClearDelimiter, "context:", msgStack)
//					return msgStack
//				}
//				msgStack.Push(msg)
//			}
//		}
//	}
//	return msgStack
//}

//func fetchLatestMessages(s *discordgo.Session, channelID string, count int) (*discordgo.Message, error) {
//	var messages *discordgo.Message
//
//	msgs, err := s.ChannelMessages(channelID, ctx.MaxFetchRecord, "", "", "")
//
//	if err != nil {
//		logger.Fatal("Error fetching channel messages:", err)
//		return messages, err
//	}
//	for _, msg := range msgs {
//		if msg.Content != "" {
//			messages = msg
//			break
//		}
//	}
//	return messages, nil
//}

//func fetchMessagesByCount(s *discordgo.Session, channelID string, count int) ([]*discordgo.Message, error) {
//	var messages []*discordgo.Message
//
//	msgs, err := s.ChannelMessages(channelID, ctx.MaxFetchRecord, "", "", "")
//
//	if err != nil {
//		logger.Fatal("Error fetching channel messages:", err)
//		return messages, err
//	}
//	for index, msg := range msgs {
//		if index < count {
//			messages = append(messages, msg)
//
//			// 打印附件
//			for _, attachment := range msg.Attachments {
//				fmt.Printf("  [Attachment] %s: %s\n", attachment.Filename, attachment.URL)
//			}
//
//			// 打印嵌入内容
//			for _, embed := range msg.Embeds {
//				fmt.Printf("  [Embed] Title: %s, Description: %s, URL: %s\n", embed.Title, embed.Description, embed.URL)
//			}
//
//			// 打印自定义表情
//			for _, reaction := range msg.Reactions {
//				fmt.Printf("  [Reaction] Emoji: %s, Count: %d\n", reaction.Emoji.Name, reaction.Count)
//			}
//		}
//	}
//	return messages, nil
//}

//func callOpenAIChat(msgStack *ds.Stack, us *ds.UserSession, resultChannel chan string) {
//	if msgStack.IsEmpty() {
//		resultChannel <- "[没有获取到任何聊天记录,无法对话]"
//		return
//	}
//
//	//打包消息列表
//	messages := make([]ds.ChatMessage, 0)
//
//	//人设
//	makeSystemRole(&messages, us.Prompt)
//
//	for !msgStack.IsEmpty() {
//		msg, _ := msgStack.Pop()
//
//		role := "user"
//		if msg.Author.ID == ctx.BotId {
//			role = "assistant"
//		}
//
//		messages = append(messages, ds.ChatMessage{
//			Role:    role,
//			Content: getCleanMsg(msg.Content),
//		})
//	}
//
//	//消息数大于20时使用概括策略,否则使用完整策略
//	if len(messages) > 20 {
//		resultChannel <- abstractChatStrategy(messages, us)
//	} else {
//		resultChannel <- fullChatStrategy(messages, us)
//	}
//}

func fullChatStrategy(messages []ds.ChatMessage, us *ds.UserSession) (resp string) {
	logger.Println("================", us.UserName, ":", us.UserChannelID, "================")
	for _, m := range messages {
		logger.Println(m.Role, ":", getCleanMsg(m.Content))
	}
	logger.Println("================================")

	result, _ := openaisdk.Chat(messages)

	return result
}

//func abstractChatStrategy(messages []ds.ChatMessage, us *ds.UserSession) (resp string) {
//	//处理数组越界问题
//	defer func() {
//		if r := recover(); r != nil {
//			logger.Println("Panic occurred:", r)
//		}
//	}()
//
//	lastIdx := len(messages) - 1
//	lastQuestion := messages[lastIdx]
//
//	messages[lastIdx] = ds.ChatMessage{
//		Role:    "user",
//		Content: "尽量详细的概括上述聊天内容",
//	}
//
//	abstract, _ := openaisdk.Chat(messages)
//	abstractMsg := make([]ds.ChatMessage, 0)
//
//	//人设
//	makeSystemRole(&abstractMsg, us.Prompt)
//
//	//上下文的概括
//	abstractMsg = append(abstractMsg, ds.ChatMessage{
//		Role:    "assistant",
//		Content: abstract,
//	})
//
//	//用户问题
//	abstractMsg = append(abstractMsg, ds.ChatMessage{
//		Role:    "user",
//		Content: lastQuestion.Content,
//	})
//
//	logger.Println("================", us.UserName, ":", us.UserChannelID, "================")
//
//	for _, m := range abstractMsg {
//		logger.Println(m.Role, ":", getCleanMsg(m.Content))
//	}
//	logger.Println("================================")
//
//	result, _ := openaisdk.Chat(abstractMsg)
//	return result
//}

func makeSystemRole(msg *[]ds.ChatMessage, prompt string) {
	*msg = append(*msg, ds.ChatMessage{
		Role:    "system",
		Content: prompt,
	})
}
