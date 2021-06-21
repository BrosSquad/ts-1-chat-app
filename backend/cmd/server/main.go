package main

import (
	"context"
	"flag"
	"net"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/BrosSquad/ts-1-chat-app/backend/di"
	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/services"
	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
)

var (
	Version string = "dev"
	Author  string = "BrosSquad Dev Team"

	dbPath   string
	logsPath string
	logLevel string
	addr     string

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

	return ":3000"
}

func getAbsolutePath(path string) string {
	var err error

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)

		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Cannot get absolute path of %s", dbPath)
		}

		return path
	}

	return ""
}

func main() {
	var err error
	ctx := context.Background()

	flag.StringVar(&dbPath, "db", "./database.sqlite", "Path to the SQLite Database file")
	flag.StringVar(&logsPath, "logs", "./logs", "Path to the root logs directory")
	flag.StringVar(&logLevel, "level", "info", "Console loggger default logging level")
	flag.StringVar(&addr, "addr", "", "Address of the HTTP2 gRPC Server")

	flag.BoolVar(&logJson, "json", false, "Log global logs as json")
	flag.BoolVar(&logToFile, "file", false, "Log global logs output to file")
	flag.BoolVar(&logToConsole, "console", true, "All logs output to file and console")

	flag.Parse()

	var output *os.File = nil

	if logToFile {
		path := path.Join(logsPath, "global.jsonl")
		output, err = utils.CreateLogFile(path, 0744)

		if err != nil {
			log.Fatal().
				Err(err).
				Msgf("Cannot open %s file for logging", path)
		}

		defer output.Close()
	}

	logging.ConfigureDefaultLogger(ctx, logLevel, output, logJson)

	dbPath = getAbsolutePath(dbPath)
	logsPath = getAbsolutePath(logsPath)

	container := di.New(di.Config{
		LogsRoot:     logsPath,
		DbPath:       dbPath,
		LogToConsole: logToConsole,
	})

	listener, err := net.Listen("tcp", getServerAddr(addr))

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error while binding net.Listen")
	}

	grpcServer := grpc.NewServer()

	services.Register(grpcServer, container)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().
			Err(err).
			Msg("error while starting grpc server")
	}
}
