package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minpeter/hijackJwtAdmin/data"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type SignupRequest struct {
	Name     string
	Email    string
	Password string
}

type LoginRequest struct {
	Email    string
	Password string
}

func main() {
	app := fiber.New()

	// Default config
	app.Use(cors.New())

	// Or extend your config for customization
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
	}))

	engine, err := data.CreateDBEngine()
	if err != nil {
		panic(err)
	}

	if err := godotenv.Load("local.env"); err != nil {
		log.Print("local.env file miss, but it's ok")
	}
	if err := godotenv.Load("production.env"); err != nil {
		log.Print("production.env file miss, but it's ok")
	}
	if temp := os.Getenv("JWT_SECRET"); temp == "" {
		log.Fatal("JWT_SECRET is empty\nplease set JWT_SECRET in local.env or production.env")
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg": "This server is a JWT authentication server.",
			"auth/signup": fiber.Map{
				"method":  "POST",
				"field":   []string{"name", "email", "password"},
				"Explain": "Creating new users and requesting tokens",
			},
			"auth/login": fiber.Map{
				"method":  "POST",
				"field":   []string{"email", "password"},
				"Explain": "Request for user token that has already been created",
			},
			"/private": fiber.Map{
				"method":  "GET",
				"Explain": "Private router that only logged in users can request",
			},
			"/public": fiber.Map{
				"method":  "GET",
				"Explain": "A public router that everyone can request",
			},
			"/admin/flag": fiber.Map{
				"method":  "GET",
				"Explain": "Catch Me If You Can",
			},
		})
	})

	auth := app.Group("/auth")
	auth.Post("/signup", func(c *fiber.Ctx) error {
		req := new(SignupRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		if req.Name == "" || req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid signup credentials")
		}

		// db에 정보 저장
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := &data.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: string(hash),
		}

		_, err = engine.Insert(user)
		if err != nil {
			return err
		}
		token, exp, err := createJWTToken(*user)
		if err != err {
			return err
		}

		// jwt 토큰 생성
		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	auth.Post("/login", func(c *fiber.Ctx) error {
		req := new(LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		if req.Email == "" || req.Password == "" {
			return fiber.NewError(fiber.StatusBadRequest, "invalid login credentials")
		}

		user := new(data.User)
		has, err := engine.Where("email = ?", req.Email).Desc("id").Get(user)
		if err != nil {
			return err
		}
		if !has {
			return fiber.NewError(fiber.StatusBadRequest, "invalid login credentials")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			return err
		}

		token, exp, err := createJWTToken(*user)
		if err != err {
			return err
		}

		// jwt 토큰 생성
		return c.JSON(fiber.Map{"token": token, "exp": exp, "user": user})
	})

	admin := app.Group("/admin")
	admin.Use(jwtWare.New(jwtWare.Config{
		// SigningKey: []byte(os.Getenv("JWT_SECRET")),
		KeyFunc: isAdmin(),
	}))
	admin.Get("/flag", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "/flag", "flag": os.Getenv("FLAG")})
	})

	private := app.Group("/private")
	private.Use(jwtWare.New(jwtWare.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	private.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "/private", "user": c.Locals("user")})
	})

	public := app.Group("/public")
	public.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "/public"})
	})

	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("BACKEND_PORT"))); err != nil {
		panic(err)
	}
}

func createJWTToken(user data.User) (string, int64, error) {
	exp := time.Now().Add(time.Minute * 10).Unix()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user id"] = user.Id
	claims["exp"] = exp
	claims["isAdmin"] = false
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return t, exp, nil
}

func isAdmin() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if !token.Claims.(jwt.MapClaims)["isAdmin"].(bool) {
			return nil, fmt.Errorf("not admin")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	}
}
