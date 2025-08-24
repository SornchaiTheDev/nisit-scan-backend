package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/adapters/rest"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/auth"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/repositories/pgx"
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
	adminRepo := pgx.NewAdminRepo(ctx, q)
	eventRepo := pgx.NewEventRepo(ctx, q)
	staffRepo := pgx.NewStaffRepository(ctx, q)
	participantRepo := pgx.NewParticipantRepo(ctx, q)
	tokenRepo := pgx.NewTokenRepository(ctx, q)
	userRepo := pgx.NewUserRepository(q)

	// Init Service
	adminService := services.NewAdminService(adminRepo)
	eventService := services.NewEventService(eventRepo)
	staffService := services.NewStaffService(staffRepo)
	participantService := services.NewParticipantService(participantRepo, userRepo)
	tokenService := services.NewTokenService(tokenRepo)
	userService := services.NewUserService(userRepo)

	// Init Auth
	authService := auth.NewGoogleOAuth(adminService, staffService)

	port := os.Getenv("PORT")

	app := fiber.New()

	origin := func() string {
		url := os.Getenv("WEB_URL")
		regex, err := regexp.Compile(`^(https?://[^/]+)`)
		if err != nil {
			log.Fatalln("WEB_URL is not valid")
		}

		origin := regex.FindString(url)

		return origin
	}()

	// Middlewares
	app.Use(cors.New(cors.Config{
		AllowOrigins:     origin,
		AllowCredentials: true,
	}))

	rest.NewAdminHandler(app, adminService)
	rest.NewEventHandler(app, adminService, eventService, staffService, participantService)
	rest.NewAuthHandler(app, authService, tokenService)
	rest.NewUserHandler(app, userService)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
