package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"go-disc-golf/internal/models"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type application struct {
	logger   *slog.Logger
	snippets *models.SnippetModel
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
		// AddSource: true,
	}))

	err := godotenv.Load()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	addr := os.Getenv("ADDRESS")
	mongodb := os.Getenv("MONGODB")

	db, err := openDB(mongodb)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		if err = db.Disconnect(context.TODO()); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	app := &application{
		logger:   logger,
		snippets: &models.SnippetModel{DB: db},
	}

	logger.Info("Pinged your deployment. You successfully connected to MongoDB!")
	logger.Info("starting server", slog.String("addr", addr))
	err = http.ListenAndServe(addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(mongodb string) (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongodb).SetServerAPIOptions(serverAPI)

	db, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	if err := db.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, err
	}

	return db, err
}
