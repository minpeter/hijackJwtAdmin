package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtWare "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/minpeter/localAuth/data"
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

	engine, err := data.CreateDBEngine()
	if err != nil {
		panic(err)
	}

	app.Post("/signup", func(c *fiber.Ctx) error {
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

	app.Post("/login", func(c *fiber.Ctx) error {
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
		// SigningKey: []byte("MuSaSiN34au0"),
		KeyFunc: isAdmin(),
	}))
	admin.Get("/flag", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "/flag", "flag": "dfc{4h1s_fl@g_1s_just_dummy_fl@g"})
	})

	private := app.Group("/private")
	private.Use(jwtWare.New(jwtWare.Config{
		SigningKey: []byte("MuSaSiN34au0"),
	}))
	private.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "/private", "user": c.Locals("user"), "jwt": c.Locals("jwt")})
	})

	public := app.Group("/public")
	public.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true, "path": "/public"})
	})

	if err := app.Listen(":3000"); err != nil {
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
	t, err := token.SignedString([]byte("MuSaSiN34au0"))
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
		return []byte("MuSaSiN34au0"), nil
	}
}
