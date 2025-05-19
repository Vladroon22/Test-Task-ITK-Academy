package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Vladroon22/TestTask-ITK-Academy/internal/database"
	"github.com/Vladroon22/TestTask-ITK-Academy/internal/handlers"
	"github.com/Vladroon22/TestTask-ITK-Academy/internal/repository"
	"github.com/Vladroon22/TestTask-ITK-Academy/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func validateUUIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if _, err := uuid.Parse(id); err != nil {
			c.JSON(400, gin.H{"error": "invalid UUID format"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err)
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	repo := repository.NewRepo(db)
	srv := service.NewService(repo)
	h := handlers.NewHandler(srv)

	addr := os.Getenv("addr")
	r := gin.Default()
	r.POST("/api/v1/wallet", h.WalletOperation)
	r.GET("/api/v1/wallet/:id", validateUUIDMiddleware(), h.GetBalance)

	go func() {
		if err := r.Run(addr); err != nil {
			log.Println(err)
			return
		}
	}()

	exitSig := make(chan os.Signal, 1)
	signal.Notify(exitSig, syscall.SIGINT, syscall.SIGTERM)
	<-exitSig

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		db.Close()
	}()
	wg.Wait()

	log.Println("Graceful shutdown")
}
