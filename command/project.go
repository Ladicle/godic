package command

import (
	"strings"
)

type ProjectCommand struct {
	Meta
}

func (c *ProjectCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *ProjectCommand) Synopsis() string {
	return "not implemented"
}

func (c *ProjectCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
