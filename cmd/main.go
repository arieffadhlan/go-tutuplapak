package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"tutuplapak-user/internal/config"
	"tutuplapak-user/internal/handlers"
	"tutuplapak-user/internal/middleware"
	"tutuplapak-user/internal/repository"
	"tutuplapak-user/internal/route"
	"tutuplapak-user/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jmoiron/sqlx"
)

func main() {
	cfg, err := config.LoadsAllAppConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbs, err := config.InitsDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	defer config.CloseDBConnection(dbs)

	app := fiber.New(fiber.Config{
		IdleTimeout: 600 * time.Second,
		ReadTimeout: 600 * time.Second,
	})

	RegRoutes(app, cfg, dbs)
	RunServer(app, cfg)
}

func RunServer(app *fiber.App, cfg *config.Config) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := app.Listen(cfg.App.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Server shutdown gracefully...")
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(timeoutCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exited")
}

func RegRoutes(app *fiber.App, cfg *config.Config, db *sqlx.DB) {
	v1 := app.Group("/api/v1", logger.New())

	v1.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "Ok"})
	})

	v1.Get("/protected-route", middleware.Protected(), func(c *fiber.Ctx) error {
		userId := c.Locals("userId").(string)

		return c.JSON(fiber.Map{"status": "Ok", "message": "Protected route", "userId": userId})
	})

	userRepo := repository.NewUserRepository(db)
	fileRepo := repository.NewFileRepository(db)
	productsRepo := repository.NewProductsRepository(db)
	purchaseRepo := repository.NewPurchaseRepository(db)

	fileService := services.NewFileService(*cfg, fileRepo)
	authService := services.NewAuthService(userRepo)
	userService := services.NewUserService(userRepo)
	productsService := services.NewProductsService(productsRepo)
	purchaseService := services.NewPurchaseService(purchaseRepo, db)

	fileHandler := handlers.NewFileHandler(fileService)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService, fileService)
	productsHandler := handlers.NewProductsHandler(productsService, fileService)
	purchaseHandler := handlers.NewPurchaseHandler(purchaseService)

	route.RegisterAuthRoutes(v1, authHandler)
	route.RegisterUserRoutes(v1, userHandler)
	route.RegisterFileRoutes(v1, fileHandler)
	route.RegisterPurchaseRoutes(v1, purchaseHandler)
	route.RegisterProductsRoutes(v1, productsHandler)
}
