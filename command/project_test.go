package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestProjectCommand_implement(t *testing.T) {
	var _ cli.Command = &ProjectCommand{}
}
