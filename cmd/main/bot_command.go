package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mglslg/go-discord-gpt/cmd/g"
	"github.com/mglslg/go-discord-gpt/cmd/g/ds"
)

func onCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//清除聊天上下文(实际上就是打印一句话,后面取聊天记录时按照它作分隔)
	us := g.GetUserSession(i.Interaction.Member.User.ID, i.Interaction.ChannelID, i.Interaction.Member.User.Username)

	if i.ApplicationCommandData().Name == "一忘皆空" {
		doForgetAllCmd(s, i, us)
	}
	if i.ApplicationCommandData().Name == "准备魔杖" {
		//doForgetAllCmd(s, i, us)
	}
}

func doForgetAllCmd(s *discordgo.Session, i *discordgo.InteractionCreate, us *ds.UserSession) {
	userMention := i.Member.User.Mention()
	replyContent := fmt.Sprintf("%s %s", userMention, ctx.ClearDelimiter)

	//todo 这边调用dify的清除聊天上下文接口，清除当前会话中的内容

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
