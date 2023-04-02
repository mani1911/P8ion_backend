package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"p8ion/config"
	"p8ion/server/router"
	"p8ion/utils"
	"strconv"
	"syscall"
	"time"
)

func Run() {
	config := config.GetConfig()

	// Initialize all the routes
	router.Init()

	server := http.Server{
		Addr:    config.Host + ":" + strconv.FormatUint(uint64(config.Port), 10),
		Handler: router.Router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	fmt.Println("Shutdown Server ...")

	// Timeout of 2s
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Print("Server Shutdown:", err)
	}

	<-ctx.Done()
	fmt.Printf("Timeout of %d", 2)

	utils.Logger.Info("Server exiting")
	fmt.Println("Server exiting")

}
