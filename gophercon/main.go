package main

import (
	"fmt"
	"io"
	"log"

	"github.com/micro/go-plugins/registry/kubernetes"
	opentracing "github.com/opentracing/opentracing-go"

	"github.com/gophercon/gc18/gophercon/actions"
	"github.com/micro/go-web"

	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/udp"
)

// ServiceName is the string name of the service, since
// it is used in multiple places, it's an exported Constant
const ServiceName = "gophercon.web"

func main() {

	registry := kubernetes.NewRegistry() //a default to using env vars for master API
	// create new web service
	service := web.NewService(
		web.Name(ServiceName),
		web.Version("latest"),
		web.Registry(registry),
		web.Address("0.0.0.0:3000"),
	)
	tracer, closer := initTracer()
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	fmt.Println(tracer)
	app := actions.App(tracer)

	// register handler
	service.Handle("/", app)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// run service
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
