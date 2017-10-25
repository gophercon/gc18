package grifts

import (
	"github.com/bketelsen/becty/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
