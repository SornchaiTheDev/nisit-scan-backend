package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/SornchaiTheDev/nisit-scan-backend/internal/adapters/rest"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	repositories "github.com/SornchaiTheDev/nisit-scan-backend/internal/repositories/pgx"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/services"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func main() {
	libs.InitEnv()

	dbUrl := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close(context.Background())

	// Init sqlc
	q := sqlc.New(conn)

	// Init repositories
	adminRepo := repositories.NewAdminRepo(q)

	// Init Service
	adminService := services.NewAdminService(adminRepo)

	port := os.Getenv("PORT")

	app := fiber.New()

	rest.NewAdminHandler(app, adminService)

	err = app.Listen(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
