package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/1219796395/myProject2/internal/conf"
	"github.com/1219796395/myProject2/internal/server"
	"github.com/1219796395/myProject2/internal/task"

	kratosLogrus "github.com/go-kratos/kratos/contrib/log/logrus/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/sirupsen/logrus"
	_ "go.uber.org/automaxprocs"
	"gopkg.in/natefinch/lumberjack.v2"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string

	id, _ = os.Hostname()

	logger log.Logger
)

func newApp(logger log.Logger, hs *server.HttpServer, ms *server.MetricsServer, gs *grpc.Server, task *task.Task) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			ms.Server,
			hs.Server,
			gs,
		),
	)
}

func main() {
	bc, err := conf.GetConf()
	if err != nil || bc == nil {
	}

	fmt.Printf("testest")

	// init logger
	logger = initLogger(bc.Server.Env)

	// init app
	app, cleanup, err := initApp(bc, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	fmt.Printf("testest")

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func initLogger(env string) log.Logger {
	logDir := os.Getenv("LOG_DIR")
	if len(logDir) == 0 {
		if env == "dev" {
			logDir = "."
		} else {
			panic("LOG_DIR is empty!")
		}
	}
	// assamble log output path
	if index := strings.LastIndex(logDir, "/"); index > 0 && index == len(logDir)-1 {
		// 去除末尾"/"符号
		logDir = logDir[:index]
	}

	// get SERVICE_NAME
	//appName := os.Getenv("APP_NAME")
	appName := "biz-platform-game-config.core-api"
	if len(appName) == 0 {
		panic("APP_NAME is empty!")
	}

	// get APP_ENV
	appEnv := os.Getenv("APP_ENV")
	if len(appEnv) == 0 {
		panic("APP_ENV is empty!")
	}

	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})
	l.SetLevel(logrus.InfoLevel)
	l.SetReportCaller(true)
	logger := log.With(kratosLogrus.NewLogger(l),
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace_id", tracing.TraceID(),
		"span_id", tracing.SpanID(),
		"client_ip", clientIp(),
	)
	// logurus set file output
	if env == "pre" || env == "prod" {
		// set logger output
		l.SetOutput(&lumberjack.Logger{
			Filename:   logDir + fmt.Sprintf("/%s.log", appName),
			MaxSize:    512, // megabytes
			Compress:   false,
			MaxAge:     14,
			MaxBackups: 4,
		})
	} else {
		l.SetOutput(os.Stdout)
	}
	// make this logger global
	log.SetLogger(logger)

	return logger
}

// TraceID returns a traceid valuer.
func clientIp() log.Valuer {
	return func(ctx context.Context) interface{} {
		common, ok := ctx.Value("clientIP").(string)
		if !ok {
			return ""
		}
		return common
	}
}
