package utils

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type WebServer interface {
	GetListenPort() int
	GetRouter() *gin.Engine
}

func GetGinEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	return engine
}

func Start(s WebServer) {
	httpServer := &http.Server{
		Addr:                         fmt.Sprintf("0.0.0.0%d", s.GetListenPort()),
		Handler:                      s.GetRouter(),
		DisableGeneralOptionsHandler: false,
		ReadHeaderTimeout:            5 * time.Second,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic("Server Start Error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Println("Shutdown failed")
	}
	<-shutdownCtx.Done()
	log.Println("Server Existing")
}
