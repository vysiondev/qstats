package utils

import "github.com/bwmarrin/discordgo"

type Embed struct {
	*discordgo.MessageEmbed
}

func CreateEmbed() Embed {
	return Embed{&discordgo.MessageEmbed{
		Color: 300543,
	}}
}

func (embed *Embed) AddTitle(t string) *Embed {
	embed.Title = t
	return embed
}

func (embed *Embed) AddImage(url string) *Embed {
	embed.Image = &discordgo.MessageEmbedImage{URL: url}
	return embed
}

func (embed *Embed) AddDescription(d string) *Embed {
	embed.Description = d
	return embed
}

func (embed *Embed) ColorOverride(c int) *Embed {
	embed.Color = c
	return embed
}

func (embed *Embed) AddThumbnail(t string) *Embed {
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: t}
	return embed
}

func (embed *Embed) AddFooter(f string) *Embed {
	embed.Footer = &discordgo.MessageEmbedFooter{Text: f}
	return embed
}

func (embed *Embed) AddTitleURL(t string) *Embed {
	embed.URL = t
	return embed
}

func (embed *Embed) AddField(name, value string) *Embed {
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: false,
	})
	return embed
}
