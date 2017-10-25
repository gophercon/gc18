package main

import (
	"io"
	"log"

	"github.com/gophercon/gc18/services/profile-srv/db"
	"github.com/gophercon/gc18/services/profile-srv/handler"
	"github.com/gophercon/gc18/services/profile-srv/proto/record"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	mot "github.com/micro/go-plugins/wrapper/trace/opentracing"

	opentracing "github.com/opentracing/opentracing-go"

	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/udp"
)

func main() {
	tracer, closer := initTracer()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	service := micro.NewService(
		micro.Name("gophercon.srv.profile"),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/profile",
			},
		),

		micro.WrapClient(mot.NewClientWrapper(tracer)),
		micro.WrapHandler(mot.NewHandlerWrapper(tracer)),

		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				db.Url = c.String("database_url")
			}
		}),
	)

	service.Init()

	db.Init()

	record.RegisterRecordHandler(service.Server(), new(handler.Record))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func initTracer() (opentracing.Tracer, io.Closer) {
	sampler := jaeger.NewConstSampler(true)
	transport, err := udp.NewUDPTransport("", 0)
	if err != nil {
		log.Fatal(err)
	}
	reporter := jaeger.NewRemoteReporter(transport)

	tracer, closer := jaeger.NewTracer(ServiceName, sampler, reporter)
	return tracer, closer
}
