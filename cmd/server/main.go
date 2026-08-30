package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Heiji57/ETB-BE/config"
	"github.com/Heiji57/ETB-BE/internal/handler"
	"github.com/Heiji57/ETB-BE/internal/infra"
	"github.com/Heiji57/ETB-BE/internal/middleware"
	repository "github.com/Heiji57/ETB-BE/internal/repository/sqlc"
	"github.com/Heiji57/ETB-BE/internal/service"
)

func main() {
	cfg := config.Load()

	pool := infra.NewPostgres(cfg)
	defer pool.Close()

	rdb := infra.NewRedis(cfg)
	defer rdb.Close()

	mailer := infra.NewMailer(cfg)

	q := repository.New(pool)

	authSvc := service.NewAuth(q, rdb, cfg, mailer)
	deviceSvc := service.NewDevice(q)
	guardianSvc := service.NewGuardian(q)
	destinationSvc := service.NewDestination(q)
	sosSvc := service.NewSos(q)
	ocrSvc := service.NewOcr(q)

	embeddedUserSvc := service.NewEmbeddedUser(q)
	locationSvc := service.NewLocation(q)
	deviceStatusSvc := service.NewDeviceStatus(q)
	navigationSvc := service.NewNavigation(q)

	authH := handler.NewAuth(authSvc)
	deviceH := handler.NewDevice(deviceSvc)
	guardianH := handler.NewGuardian(guardianSvc)
	destinationH := handler.NewDestination(destinationSvc)
	sosH := handler.NewSos(sosSvc)
	ocrH := handler.NewOcr(ocrSvc)
	embeddedUserH := handler.NewEmbeddedUser(embeddedUserSvc)
	locationH := handler.NewLocation(locationSvc)
	deviceStatusH := handler.NewDeviceStatus(deviceStatusSvc)
	navigationH := handler.NewNavigation(navigationSvc)
	hub := handler.NewHub()
	wsH := handler.NewWs(hub, q)

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-Device-Key"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	v1 := r.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.POST("/signup", authH.Signup)
		auth.POST("/verify-email", authH.VerifyEmail)
		auth.POST("/login", authH.Login)
		auth.POST("/login/google", authH.LoginGoogle)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", middleware.Auth(cfg), authH.Logout)
	}

	authRequired := v1.Group("")
	authRequired.Use(middleware.Auth(cfg))
	{
		authRequired.POST("/devices", deviceH.Create)
		authRequired.GET("/devices", deviceH.List)
		authRequired.DELETE("/devices/:deviceId", deviceH.Delete)
		authRequired.GET("/devices/:deviceId/status", deviceStatusH.GetForUser)

		authRequired.POST("/guardians", guardianH.Create)
		authRequired.GET("/guardians", guardianH.List)
		authRequired.DELETE("/guardians/:guardianId", guardianH.Delete)

		authRequired.POST("/destinations", destinationH.Create)
		authRequired.GET("/destinations", destinationH.List)
		authRequired.DELETE("/destinations/:destinationId", destinationH.Delete)

		authRequired.POST("/sos", sosH.Create)
		authRequired.GET("/sos", sosH.List)

		authRequired.POST("/ocr-logs", ocrH.Create)

		authRequired.POST("/navigation/sessions", navigationH.CreateSession)
		authRequired.POST("/navigation/sessions/:sessionId/instructions", navigationH.UpdateInstruction)
		authRequired.PATCH("/navigation/sessions/:sessionId/status", navigationH.UpdateSessionStatus)
	}

	// Embedded API — /api/... (no JWT, userId in body)
	embedded := r.Group("/api")
	{
		embedded.POST("/users", embeddedUserH.Register)
		embedded.GET("/users/:userId", embeddedUserH.Get)
		embedded.POST("/location", locationH.Save)
		embedded.POST("/sos", middleware.DeviceKeyAuth(q), sosH.CreateEmbedded)
		embedded.POST("/device/status", deviceStatusH.Update)
		embedded.GET("/device/status/:deviceId", deviceStatusH.Get)
		embedded.POST("/ocr/results", ocrH.Create)
		embedded.POST("/destinations", destinationH.CreateEmbedded)
		embedded.GET("/destinations", destinationH.ListEmbedded)
		embedded.GET("/navigation/sessions/:sessionId/instruction", navigationH.GetLatestInstruction)
	}

	r.GET("/ws/device", middleware.Auth(cfg), wsH.Handle)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Printf("server started on :%s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
