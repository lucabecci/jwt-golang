package main

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	jwtware "github.com/gofiber/jwt"
)

const jwtSecret = "SECRET"

func authRequired() func(ctx *fiber.Ctx) {
	return jwtware.New(jwtware.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) {
			ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
		SigningKey: []byte(jwtSecret),
	})
}

func main() {
	app := fiber.New()

	app.Use(middleware.Logger())

	app.Get("/", func(ctx *fiber.Ctx) {
		ctx.Send("Hello World")
	})

	app.Get("/hello", authRequired(), func(ctx *fiber.Ctx) {
		user := ctx.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		id := claims["user_id"].(string)
		ctx.Send(fmt.Sprintf("Hello : %s", id))
	})

	app.Post("/login", login)

	err := app.Listen(3000)
	if err != nil {
		panic(err)
	}
}

func login(ctx *fiber.Ctx) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var body request
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Please send all camps",
		})
		return
	}

	if body.Email != "test@test.com" || body.Password != "test123" {
		ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Bad credentials",
		})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = "1"
	claims["exp"] = time.Now().Add(time.Hour * 24 * 3) //three days

	tkn, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		ctx.SendStatus(fiber.StatusInternalServerError)
		return
	}

	ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": tkn,
		"user": struct {
			ID    int    `json:"id"`
			Email string `json:"email"`
		}{
			ID:    1,
			Email: "test@test.com",
		},
	})
}
