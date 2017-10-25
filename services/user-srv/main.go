package main

import (
	"io"
	"log"

	"github.com/gophercon/gc18/services/user-srv/db"
	"github.com/gophercon/gc18/services/user-srv/handler"
	proto "github.com/gophercon/gc18/services/user-srv/proto/account"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/server"
	mot "github.com/micro/go-plugins/wrapper/trace/opentracing"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

const ServiceName = "gophercon.srv.user"

func main() {

	tracer, closer, err := initTracer()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	service := micro.NewService(
		micro.Name(ServiceName),
		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/user",
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
	client.DefaultClient = client.NewClient(
		client.Wrap(
			mot.NewClientWrapper(tracer)),
	)
	server.DefaultServer = server.NewServer(
		server.WrapHandler(mot.NewHandlerWrapper(tracer)),
	)
	service.Init()
	db.Init()

	proto.RegisterAccountHandler(service.Server(), new(handler.Account))

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func initTracer() (opentracing.Tracer, io.Closer, error) {
	cfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	closer, err := cfg.InitGlobalTracer(
		"serviceName",
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
		return nil, closer, err
	}
	defer closer.Close()
	tracer, closer, err := cfg.New(ServiceName, jaegercfg.Logger(jLogger), jaegercfg.Metrics(jMetricsFactory))
	return tracer, closer, err
}
