package main

import (
	"flag"
	"log"
	"net"
	"path/filepath"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services"
)

var (
	dbPath string
	addr   string
)

func main() {
	var err error

	flag.StringVar(&dbPath, "db", "./database.sqlite", "Path to the SQLite Database file")
	flag.StringVar(&addr, "addr", ":3000", "Addres of the HTTP2 Server")

	flag.Parse()

	if !filepath.IsAbs(dbPath) {
		dbPath, err = filepath.Abs(dbPath)

		if err != nil {
			log.Fatalf("Cannot get absolute path of %s: %v", dbPath, err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		log.Fatalf("error connecting to SQLITE3 Database: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Message{})

	if err != nil {
		log.Fatalf("Failed to migrate DATABASE: %v", err)
	}

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("error while binding net.Listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	services.Register(grpcServer, db)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatalf("error while starting grpc server: %v", err)
	}
}
