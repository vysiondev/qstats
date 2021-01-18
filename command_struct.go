package main

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name         string
	Description  string
	Shorthands   []string
	RequiresUser bool
	OwnerOnly    bool
	Notes        string
	RunFunction  func(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, details RunDetails) error
}
