package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"
	"tutuplapak-user/internal/config"

	"github.com/gofiber/fiber/v2"
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
	v1 := app.Group("/v1")

	v1.Get("/health-check", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "Ok"})
	})

	// userRepo := repository.NewUserRepository(db)
	// fileRepo := repository.NewFileRepository(db)
	// productsRepo := repository.NewProductsRepository(db)
	// purchaseRepo := repository.NewPurchaseRepository(db)

	// userService := service.NewUserService(userRepo)
	// authService := service.NewAuthService(userRepo)
	// fileService := service.NewFileService(fileRepo)
	// productsService := service.NewProductsService(productsRepo)
	// purchaseService := service.NewPurchaseService(purchaseRepo)

	// userHandler := handler.NewUserHandler(userUseCase)
	// authHandler := handler.NewAuthHandler(authUseCase)
	// fileHandler := handler.NewFileHandler(fileUseCase)
	// productsHandler := handler.NewProductsHandler(productsUseCase)
	// purchaseHandler := handler.NewPurchaseHandler(purchaseUseCase)

	// route.RegisterUserRoutes(v1, userHandler)
	// route.RegisterAuthRoutes(v1, authHandler)
	// route.RegisterFileRoutes(v1, fileHandler)
	// route.RegisterProductsRoutes(v1, productsHandler)
	// route.RegisterPurchaseRoutes(v1, purchaseHandler)
}
