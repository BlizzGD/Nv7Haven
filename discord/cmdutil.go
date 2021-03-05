package discord

import (
	"github.com/bwmarrin/discordgo"
)

type normalResp struct {
	msg *discordgo.MessageCreate
	b   *Bot
}

func (n *normalResp) Error(err error) bool {
	if err != nil {
		n.b.dg.ChannelMessageSend(n.msg.ChannelID, "Error: "+err.Error())
		return true
	}
	return false
}

func (n *normalResp) ErrorMessage(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, "Error: "+msg)
}

func (n *normalResp) Resp(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, msg)
}

func (n *normalResp) Message(msg string) {
	n.b.dg.ChannelMessageSend(n.msg.ChannelID, msg)
}

func (n *normalResp) Embed(emb *discordgo.MessageEmbed) {
	n.b.dg.ChannelMessageSendEmbed(n.msg.ChannelID, emb)
}

func (b *Bot) newMsgNormal(m *discordgo.MessageCreate) msg {
	return msg{
		Author:    m.Author,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	}
}

func (b *Bot) newRespNormal(m *discordgo.MessageCreate) rsp {
	return &normalResp{
		msg: m,
		b:   b,
	}
}

type slashResp struct {
	i          *discordgo.InteractionCreate
	b          *Bot
	hasReplied bool
}

func (s *slashResp) Error(err error) bool {
	if err != nil {
		s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Flags:   1 << 6,
				Content: "Error: " + err.Error(),
			},
		})
		return true
	}
	return false
}

func (s *slashResp) ErrorMessage(msg string) {
	s.hasReplied = true
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Flags:   1 << 6,
			Content: "Error: " + msg,
		},
	})
}

func (s *slashResp) Resp(msg string) {
	s.hasReplied = true
	s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
}

func (s *slashResp) Message(msg string) {
	if !s.hasReplied {
		s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Content: msg,
			},
		})
		s.hasReplied = true
		return
	}
	s.b.dg.ChannelMessageSend(s.i.ChannelID, msg)
}

func (s *slashResp) Embed(emb *discordgo.MessageEmbed) {
	if !s.hasReplied {
		s.b.dg.InteractionRespond(s.i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionApplicationCommandResponseData{
				Embeds: []*discordgo.MessageEmbed{emb},
			},
		})
		s.hasReplied = true
		return
	}
	s.b.dg.ChannelMessageSendEmbed(s.i.ChannelID, emb)
}

func (b *Bot) newMsgSlash(i *discordgo.InteractionCreate) msg {
	return msg{
		Author:    i.Member.User,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
	}
}

func (b *Bot) newRespSlash(i *discordgo.InteractionCreate) rsp {
	return &slashResp{
		i: i,
		b: b,
	}
}

func (b *Bot) pageSwitchHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	handler, exists := b.pages[r.MessageID]
	if exists {
		handler.Handler(r)
	}
}
