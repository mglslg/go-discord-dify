package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/go-discord-dify/cmd/difysdk"
	"github.com/mglslg/go-discord-dify/cmd/g"
	"github.com/mglslg/go-discord-dify/cmd/g/ds"
	"regexp"
	"strings"
	"time"
)

var (
	ctx = &g.AppContext
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

	ctx.ClearCmdDesc = strings.Replace(ctx.ClearCmdDesc, "${botName}", ctx.BotName, -1)

	createCmd(session, ctx.ClearCmd, ctx.ClearCmdDesc)

	//监听命令
	session.AddHandler(onCommand)

	//监听消息
	session.AddHandler(onMsgCreate)

	//可以使用session.ApplicationCommands()把所有的命令取出来，然后删掉，再重建

	//删除命令方法(从聊天室的integration里面找)
	//session.ApplicationCommandDelete(ctx.ApplicationId, ctx.GuildId, "1268844898971287562")

	return session, nil
}

func createCmd(session *discordgo.Session, cmdName string, cmdDesc string) {
	_, cmdErr := session.ApplicationCommandCreate(ctx.ApplicationId, "", &discordgo.ApplicationCommand{
		Name:        cmdName,
		Description: cmdDesc,
	})
	if cmdErr != nil {
		logger.Fatal("create "+cmdName+" cmd error: ", cmdErr)
	}
}

// todo 这边可以处理一下长消息变成文档上传后会如何
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
				////异步获取响应结果并提示[正在输入],go关键字后是生产端,asyncResponse中的select是消费端
				respChannel := make(chan string)
				go callDifyChat(m.Content, us, respChannel)
				asyncResponse(s, m, us, respChannel)
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

func callDifyChat(msg string, us *ds.UserSession, resultChannel chan string) {
	replayContent, conversationId, e := difysdk.Chat(getCleanMsg(msg), us.UserName, us.ConversationID)

	if us.ConversationID == "" {
		us.ConversationID = conversationId
	}

	if e != nil {
		logger.Fatal("failed to request dify api", e)
	}

	resultChannel <- replayContent
}
