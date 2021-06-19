package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/services"
)

func main() {
	db, err := gorm.Open(sqlite.Open("../db/database.sqlite"), &gorm.Config{})

	if err != nil {
		log.Fatalf("error connecting to SQLITE3 Database: %v", err)
	}

 	err = db.AutoMigrate(&models.User{}, &models.Message{})

	if err != nil {
		log.Fatalf("Failed to migrate DATABASE: %v", err)
	}

	listener, err := net.Listen("tcp", ":3000")

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
