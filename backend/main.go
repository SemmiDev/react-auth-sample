package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"
)

type Server struct {
	config     Config
	tokenMaker Maker
	dataStore  *DataStore
	router     *fiber.App
}

func main() {
	appConfig, err := LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	tokenMaker, err := NewJWTMaker(appConfig.TokenSymmetricKey)
	if err != nil {
		log.Fatal("cannot create token maker:", err)
	}

	conn, err := sql.Open(appConfig.DBDriver, appConfig.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	if err := conn.Ping(); err != nil {
		log.Fatal("cannot ping db:", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(10)
	conn.SetConnMaxLifetime(time.Hour)

	fiberConfig := fiber.Config{
		AppName:       "backend-for-website",
		StrictRouting: true,
		CaseSensitive: true,
	}

	s := &Server{
		config:     appConfig,
		tokenMaker: tokenMaker,
		dataStore:  NewDataStore(conn),
		router:     fiber.New(fiberConfig),
	}

	s.router.Use(cors.New())
	s.router.Use(logger.New())

	s.router.Post("api/auth/register", s.registerUserHandler)
	s.router.Post("api/auth/login", s.loginUserHandler)

	s.RunHttpServer()
}

func (s *Server) RunHttpServer() {
	log.Fatal(s.router.Listen(s.config.HTTPServerAddress))
}
