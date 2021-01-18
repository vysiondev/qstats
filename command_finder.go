package main

import "strings"

func (b *BaseHandler) FindCommandIndex(name string) *int {
	for i, cmd := range b.CommandList {
		if cmd.Name == strings.ToLower(name) {
			return &i
		} else {
			if len(cmd.Shorthands) > 0 {
				for _, short := range cmd.Shorthands {
					if short == strings.ToLower(name) {
						return &i
					}
				}
			}
		}
	}
	return nil
}
