package opentelemetry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/viper"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/log/global"
	otm "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"

	//"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type OtelMetrics struct {
	TotalReqCounter otm.Int64UpDownCounter
}

func Setup(ctx context.Context) (shutdownFn []func(context.Context) error, otelMetrics *OtelMetrics, err error) {

	var otelServiceName = semconv.ServiceNameKey.String(viper.GetString("app.name"))
	var otelServiceVersion = semconv.ServiceVersionKey.String(viper.GetString("app.version"))
	var grpcEndpoint = viper.GetString("otel.endpoint")
	var otmMetrics = OtelMetrics{}

	otelResource, errResource := resource.New(context.Background(),
		resource.WithAttributes(
			// The service name used to display traces in backends
			otelServiceName,
			// The service version used to display traces in backends
			otelServiceVersion,
		),
	)
	if errResource != nil {
		panic(fmt.Errorf("error setting up opentelemetry resource [%w]", errResource))
	}

	var shutdownFuncs []func(context.Context) error
	var shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	grpcCnx, errGrpc := initGrpcConn(grpcEndpoint)
	if errGrpc != nil {
		return nil, &otmMetrics, errGrpc
	}

	// Logger provider
	loggerProvider, err := initLoggerProvider(ctx, otelResource, grpcCnx)
	if err != nil {
		handleErr(err)
		return shutdownFuncs, &otmMetrics, err
	}
	global.SetLoggerProvider(loggerProvider)
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)

	// Tracer provider
	tracerProvider, err := initTracerProvider(ctx, otelResource, grpcCnx)
	if err != nil {
		handleErr(err)
		return shutdownFuncs, &otmMetrics, err
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)

	// Meter provider
	meterProvider, errMeter := initMeterProvider(ctx, otelResource, grpcCnx)
	if errMeter != nil {
		handleErr(errMeter)
		return shutdownFuncs, &otmMetrics, err
	}
	otel.SetMeterProvider(meterProvider)
	meter := meterProvider.Meter(viper.GetString("app.name") + "/metrics")
	totalCounter, errCount := meter.Int64UpDownCounter("http.server.requests.total", otm.WithUnit("1"), otm.WithDescription("total number of HTTP requests"))

	if errCount != nil {
		handleErr(errMeter)
		return shutdownFuncs, &otmMetrics, errMeter
	}

	otmMetrics.TotalReqCounter = totalCounter

	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	return shutdownFuncs, &otmMetrics, nil
}

func initLoggerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*log.LoggerProvider, error) {
	exporter, err := otlploggrpc.New(ctx, otlploggrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	processor := log.NewBatchProcessor(exporter)
	provider := log.NewLoggerProvider(
		log.WithResource(res),
		log.WithProcessor(processor),
	)
	return provider, nil
}

func initTracerProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*sdktrace.TracerProvider, error) {
	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// Set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider, nil
}

func initMeterProvider(ctx context.Context, res *resource.Resource, conn *grpc.ClientConn) (*metric.MeterProvider, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}
	meterProvider := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			// Default is 1m. Set to 3s for demonstrative purposes.
			metric.WithInterval(3*time.Second))),
	)
	return meterProvider, nil
}

// Connect the OpenTelemetry Collector through local gRPC connection.
func initGrpcConn(grpcEndpoint string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(grpcEndpoint,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}
