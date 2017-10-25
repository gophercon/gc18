package main

import (
	"io"
	"log"

	proto "github.com/gophercon/gc18/services/config-srv/proto/config"
	"github.com/micro/cli"
	"github.com/micro/go-micro"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"

	"github.com/gophercon/gc18/services/config-srv/config"
	"github.com/gophercon/gc18/services/config-srv/handler"
	mot "github.com/micro/go-plugins/wrapper/trace/opentracing"

	// db
	"github.com/gophercon/gc18/services/config-srv/db"
	"github.com/gophercon/gc18/services/config-srv/db/mysql"

	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
)

const ServiceName = "gophercon.srv.config"

func main() {

	tracer, closer, err := initTracer()
	if err != nil {
		panic(err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	service := micro.NewService(
		micro.Name(ServiceName),
		micro.Version("latest"),

		micro.Flags(
			cli.StringFlag{
				Name:   "database_url",
				EnvVar: "DATABASE_URL",
				Usage:  "The database URL e.g root@tcp(127.0.0.1:3306)/config",
			},
		),
		micro.WrapClient(mot.NewClientWrapper(tracer)),
		micro.WrapHandler(mot.NewHandlerWrapper(tracer)),
		// Add for MySQL configuration
		micro.Action(func(c *cli.Context) {
			if len(c.String("database_url")) > 0 {
				mysql.Url = c.String("database_url")
			}
		}),
	)

	service.Init()

	proto.RegisterConfigHandler(service.Server(), new(handler.Config))

	// subcriber to watches
	service.Server().Subscribe(service.Server().NewSubscriber(config.WatchTopic, config.Watcher))

	if err := config.Init(); err != nil {
		log.Fatal(err)
	}

	if err := db.Init(); err != nil {
		log.Fatal(err)
	}

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
