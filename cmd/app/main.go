package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/adapters/rest"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/auth"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	repositories "github.com/SornchaiTheDev/nisit-scan-backend/internal/repositories/pgx"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	libs.InitEnv()

	dbUrl := os.Getenv("DATABASE_URL")
	log.Println("Database URL: ", dbUrl)

	ctx := context.Background()

	conn, err := pgxpool.New(ctx, dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	if err := conn.Ping(context.Background()); err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected")

	// Init sqlc
	q := sqlc.New(conn)

	// Init repositories
	adminRepo := repositories.NewAdminRepo(ctx, q)
	eventRepo := repositories.NewEventRepo(ctx, q)
	staffRepo := repositories.NewStaffRepository(ctx, q)
	participantRepo := repositories.NewParticipantRepo(ctx, q)
	tokenRepo := repositories.NewTokenRepository(ctx, q)

	// Init Service
	adminService := services.NewAdminService(adminRepo)
	eventService := services.NewEventService(eventRepo)
	staffService := services.NewStaffService(staffRepo)
	participantService := services.NewParticipantService(participantRepo)
	tokenService := services.NewTokenService(tokenRepo)

	// Init Auth
	authService := auth.NewGoogleOAuth(adminService, staffService)

	port := os.Getenv("PORT")

	app := fiber.New()

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("WEB_URL"),
		AllowCredentials: true,
	}))

	rest.NewAdminHandler(app, adminService)
	rest.NewEventHandler(app, adminService, eventService, staffService, participantService)
	rest.NewAuthHandler(app, authService, tokenService)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
