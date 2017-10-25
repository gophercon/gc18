package actions

import (
	"errors"
	"time"

	"github.com/gobuffalo/buffalo"
	mware "github.com/gophercon/gc18/gophercon/middleware"
	"github.com/gophercon/gc18/services/user-srv/proto/account"
	proto "github.com/gophercon/gc18/services/user-srv/proto/account"
	"github.com/micro/go-micro/client"

	mot "github.com/micro/go-plugins/wrapper/trace/opentracing"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	slow(c)
	//cl := client.NewClient(client.Wrap(mot.NewClientWrapper(Tracer)))
	client.DefaultClient = client.NewClient(
		client.Wrap(
			mot.NewClientWrapper(Tracer)),
	)
	ctx := mware.MetadataContext(c)
	user := proto.NewAccountClient("account", client.DefaultClient)
	user.Search(ctx, &account.SearchRequest{})
	return c.Render(200, r.HTML("index.html"))
}

//BadHandler returns an error
func BadHandler(c buffalo.Context) error {
	return c.Error(401, errors.New("Unauthorized!"))
}
func slow(c buffalo.Context) {
	sp := mware.ChildSpan("slow", c)
	defer sp.Finish()
	time.Sleep(10 * time.Millisecond)
}
