package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arnab333/golang-employee-management/routes"
	"github.com/arnab333/golang-employee-management/services"
	"github.com/joho/godotenv"
)

func main() {

	router := routes.Init()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	closeMongoConn := services.InitMongoConnection()

	closeRedisConn := services.InitRedisConnection()

	go services.CronInit()

	defer closeMongoConn()

	defer closeRedisConn()

	srv := &http.Server{
		Addr:    ":5000",
		Handler: router,
	}

	defer gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("timeout of 5 seconds.")
	log.Println("Server exiting")
}
