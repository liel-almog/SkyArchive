package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lielalmog/file-uploader/backend/database"
	"github.com/lielalmog/file-uploader/backend/kafka"
)

const addr = ":8080"

var app *fiber.App

func Serve() {
	app = fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})
	app.Use(recover.New())

	setupRouter(app)

	fmt.Println("Server strating on port", addr)

	if err := app.Listen(addr); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Shutdown(ctx context.Context) error {
	database.GetDB().Close()
	kafka.GetKafkaProducer().Close()

	err := app.Shutdown()

	if err != nil {
		return err
	}

	return nil
}
