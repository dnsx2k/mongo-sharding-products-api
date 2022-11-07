package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	productshttphandler "github.com/dnsx2k/mongo-sharding-products-api/cmd/httphandlers"
	"github.com/dnsx2k/mongo-sharding-products-api/pkg/lookupclient"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	mongoClientPrimary, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoPrimaryShardConnectionString))
	if err != nil {
		panic(err)
	}

	mongoClientHot, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDBHotShardConnectionString))
	if err != nil {
		panic(err)
	}

	lookupClient := lookupclient.New(cfg.LookupServiceBaseURL)

	productsHTTPHandler := productshttphandler.New(mongoClientPrimary, mongoClientHot, lookupClient)

	router := gin.Default()
	apiV1 := router.Group("v1")
	productsHTTPHandler.Setup(apiV1)

	server := http.Server{
		Addr:         "localhost:8085",
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Graceful HTTP shutdown
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	signal.Stop(signalChan)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
