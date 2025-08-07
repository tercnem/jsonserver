package config

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

var FiberConfig fiber.Config
var FiberPort string

func InitFiberConfig() {

	FiberConfig = fiber.Config{
		AppName:         "jsonserver",
		CaseSensitive:   true,
		Immutable:       true,
		BodyLimit:       1 * 1024 * 1024,
		Concurrency:     255,
		ReadBufferSize:  10 * 1024,
		WriteBufferSize: 10 * 1024,
		ReadTimeout:     time.Second * 10,
		WriteTimeout:    time.Second * 10,
		IdleTimeout:     time.Second * 10,

		EnablePrintRoutes: false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			c.WriteString(err.Error())
			return nil
		},
	}

}
