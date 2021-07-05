package main

import (
	"context"
	"flag"
	"github.com/BrosSquad/ts-1-chat-app/backend/handlers"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/fatih/color"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/BrosSquad/ts-1-chat-app/backend/di"
	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/middleware"
	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
)

var (
	Version = "dev"
	Author  = "BrosSquad Dev Team"

	configPath string
	logsPath   string
	logLevel   string
	addr       string
	httpAddr   string

	logJson      bool
	logToFile    bool
	logToConsole bool
)

func getServerAddr(flag string) string {
	if flag != "" {
		return os.ExpandEnv(flag)
	}

	env := os.Getenv("SERVER_ADDR")

	if env != "" {
		return os.ExpandEnv(env)
	}

	return ":4000"
}

func getHTTPServerAddr(flag string) string {
	if flag != "" {
		return os.ExpandEnv(flag)
	}

	env := os.Getenv("HTTP_SERVER_ADDR")

	if env != "" {
		return os.ExpandEnv(env)
	}

	return ":4001"
}

func main() {
	flag.StringVar(&configPath, "config", ".", "Path to the configuration directory")
	flag.StringVar(&logsPath, "logs", "./logs", "Path to the root logs directory")
	flag.StringVar(&logLevel, "level", "trace", "Console logger default logging level")
	flag.StringVar(&addr, "addr", "", "Address of the HTTP2 gRPC Server")
	flag.StringVar(&httpAddr, "http", "", "Address of the HTTP Server")

	flag.BoolVar(&logJson, "json", false, "Log global logs as json")
	flag.BoolVar(&logToFile, "file", false, "Log global logs output to file")
	flag.BoolVar(&logToConsole, "console", true, "All logs output to file and console")

	flag.Parse()

	color.Blue("Author: %s\t", Author)
	color.Green("Version: %s", Version)

	var err error
	ctx, cancel := context.WithCancel(context.Background())

	exit := make(chan os.Signal, 1)

	signal.Notify(exit, os.Interrupt)

	if logToFile {
		p := path.Join(logsPath, "global.jsonl")
		output, err := utils.CreateLogFile(p, 0644)

		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Cannot open %s file for logging", p)
		}

		defer output.Close()

		logging.ConfigureDefaultLogger(ctx, logLevel, output, logToConsole, logJson)
	} else {
		logging.ConfigureDefaultLogger(ctx, logLevel, nil, logToConsole, logJson)
	}

	log.Trace().Str("logsPath", logsPath).Msg("Path to Logs")

	log.Trace().
		Str("configPath", configPath).
		Msg("Path to config directory")

	container := di.New(
		di.Config{
			LogsRoot:     logsPath,
			LogToConsole: logToConsole,
		},
		configPath,
	)

	addr = getServerAddr(addr)
	httpAddr = getHTTPServerAddr(httpAddr)

	log.Trace().
		Str("addr", addr).
		Msg("Starting the gRPC server")

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error while binding net.Listen")
	}

	log.Trace().
		Str("addr", addr).
		Msg("Starting the server")

	grpcServer := grpc.NewServer(
		middleware.Register(container)...,
	)

	handlers.Register(grpcServer, container)
	handlers.RegisterHTTP(container)

	go func() {
		err = grpcServer.Serve(listener)

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error while starting grpc server")
		}
	}()

	log.Trace().
		Str("addr", addr).
		Msg("Starting the HTTP server")

	srv := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	go func() {
		err := srv.ListenAndServe()

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error while starting http server")
		}
	}()

	log.Info().
		Str("addr", addr).
		Msg("Server started")
	log.Info().
		Str("addr", httpAddr).
		Msg("HTTP Server started")

	<-exit
	log.Trace().Msg("Signal Interrupt detected")

	cancel()

	log.Debug().Msg("Stopping gRPC server")

	grpcServer.GracefulStop()

	log.Debug().Msg("Stopping HTTP server")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Error while HTTP server Shutdown")
	}

	log.Trace().Msg("Exiting")
}
