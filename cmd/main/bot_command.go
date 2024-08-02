package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/go-discord-dify/cmd/difysdk"
	"github.com/mglslg/go-discord-dify/cmd/g"
	"github.com/mglslg/go-discord-dify/cmd/g/ds"
)

func onCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//清除聊天上下文(实际上就是打印一句话,后面取聊天记录时按照它作分隔)
	us := g.GetUserSession(i.Interaction.Member.User.ID, i.Interaction.ChannelID, i.Interaction.Member.User.Username)

	if i.ApplicationCommandData().Name == g.AppContext.ClearCmd {
		doForgetAllCmd(s, i, us)
	}
}

func doForgetAllCmd(s *discordgo.Session, i *discordgo.InteractionCreate, us *ds.UserSession) {
	userMention := i.Member.User.Mention()
	replyContent := fmt.Sprintf("%s %s", userMention, ctx.ClearDelimiter)

	conversationId := us.ConversationID
	if conversationId != "" {
		result, err := difysdk.DeleteConversation(conversationId, us.UserName)
		if err != nil || result != "success" {
			replyContent = result
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: replyContent,
		},
	})

	if err != nil {
		logger.Println("Error responding to slash command: ", err)
	}
}
