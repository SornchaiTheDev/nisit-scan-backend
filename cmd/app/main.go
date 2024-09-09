package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/adapters/rest"
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
	conn, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Init sqlc
	q := sqlc.New(conn)

	// Init repositories
	adminRepo := repositories.NewAdminRepo(q)
	eventRepo := repositories.NewEventRepo(q)
	staffRepo := repositories.NewStaffRepository(q)
	participantRepo := repositories.NewParticipantRepo(q)

	// Init Service
	adminService := services.NewAdminService(adminRepo)
	eventService := services.NewEventService(eventRepo)
	staffService := services.NewStaffService(staffRepo)
	participantService := services.NewParticipantService(participantRepo)

	port := os.Getenv("PORT")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))

	rest.NewAdminHandler(app, adminService)
	rest.NewEventHandler(app, eventService, staffService, participantService)

	err = app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
