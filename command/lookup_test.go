package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestLookupCommand_implement(t *testing.T) {
	var _ cli.Command = &LookupCommand{}
}
