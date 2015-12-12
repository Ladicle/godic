package main

import (
	"github.com/ladicle/godic/command"
	"github.com/mitchellh/cli"
)

func Commands(meta *command.Meta) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"search": func() (cli.Command, error) {
			return &command.SearchCommand{
				Meta: *meta,
			}, nil
		},
		"project": func() (cli.Command, error) {
			return &command.ProjectCommand{
				Meta: *meta,
			}, nil
		},
		"lookup": func() (cli.Command, error) {
			return &command.LookupCommand{
				Meta: *meta,
			}, nil
		},
	}
}
