package grifts

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gophercon/gc18/actions"
)

func init() {
	buffalo.Grifts(actions.App())
}
