package actions

import (
	"errors"
	"time"

	"github.com/gobuffalo/buffalo"

	mware "github.com/gophercon/gc18/gophercon/middleware"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	slow(c)
	return c.Render(200, r.HTML("index.html"))
}

//BadHandler returns an error
func BadHandler(c buffalo.Context) error {
	return c.Error(401, errors.New("Unauthorized!"))
}
func slow(c buffalo.Context) {
	sp := mware.ChildSpan("slow", c)
	defer sp.Finish()
	time.Sleep(1 * time.Millisecond)
}
