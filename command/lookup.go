package command

import (
	"strings"
)

type LookupCommand struct {
	Meta
}

func (c *LookupCommand) Run(args []string) int {
	// Write your code here

	return 0
}

func (c *LookupCommand) Synopsis() string {
	return "not implemented"
}

func (c *LookupCommand) Help() string {
	helpText := `

`
	return strings.TrimSpace(helpText)
}
