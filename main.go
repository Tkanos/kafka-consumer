package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Shopify/sarama"
	// import mssql driver used by sqlx
	_ "github.com/denisenkom/go-mssqldb"
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	"github.com/tkanos/kafka-consumer/config"
	"github.com/tkanos/kafka-consumer/consumers"
	"github.com/tkanos/kafka-consumer/myService"
)

var (
	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func init() {
	// read config file and log an error if it's not present
	err := config.InitConfig()
	if err != nil {
		panic(err)
	}
}

func main() {

	// Tracing domain.
	var tracer opentracing.Tracer
	{
		zipkinAddr := config.Config.ZipkinURI
		if zipkinAddr != "" {
			//create a collector
			collector, err := zipkin.NewHTTPCollector(zipkinAddr)
			if err != nil {
				logger.Fatal("Unable to create a zipkin collector : ", err)
			}
			defer collector.Close()

			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, true, "kafka-consumer:"+strconv.Itoa(config.Config.Port), "kafka-consumer"),
				zipkin.TraceID128Bit(true),
			)
			if err != nil {
				logger.Fatal("Unbale to create a Zipkin Tracer : ", err)
			}

			// explicitly set our tracer to be the default tracer.
			opentracing.SetGlobalTracer(tracer)

		} else {
			tracer = opentracing.GlobalTracer() // no-op
		}
	}

	// create web server just for the /healthz
	httpAddr := ":" + strconv.Itoa(config.Config.Port)
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	go func() {
		http.ListenAndServe(httpAddr, mux)
	}()

	// check kafka config
	if config.Config.KafkaBrokers == "" {
		logger.Fatal("Kafka broker is empty.")
	}

	// Create myService service
	var trackingService myService.MyServiceTracker
	{
		trackingService = myService.NewService()                       // Create the service
		trackingService = myService.NewLoggingService(trackingService) // add logging ability to our service
		trackingService = myService.NewTracingService(trackingService) // add global tracing

		//subscribe to topics
		loginConsumer, err := consumers.Subscribe(context.Background(),
			strings.Split(config.Config.KafkaBrokers, ";"),
			"my-Service",
			config.Config.KafkaTopic,
			sarama.OffsetNewest,
			myService.MakeMyServiceTrackerEndpoint(trackingService))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start login-consumer: %s", err)
			os.Exit(-3)
		}
		defer loginConsumer.Close()
	}

	//quit application
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan
	fmt.Fprintf(os.Stdout, "The user choose to interrupt the program.")
}
