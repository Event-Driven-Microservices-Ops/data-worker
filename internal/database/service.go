package database

import (
	"context"
	"log"
	"time"

	"github.com/urbaniakmichal/data-worker/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DataBaseService struct {
	cfg         *config.Config
	mongoClient *mongo.Client
}

func NewDataBaseService(cfg *config.Config, mC *mongo.Client) *DataBaseService {
	return &DataBaseService{
		cfg:         cfg,
		mongoClient: mC,
	}
}

func ConnectToMongo(cfg *config.Config) *mongo.Client {
	opts := options.Client().ApplyURI(cfg.DatabaseURL)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Pinged database. Failed connecting to MongoDB: %v", err)
	}

	return client
}

func (ds *DataBaseService) InsertBatch(ctx context.Context, documents []any) error {
	if len(documents) == 0 {
		return nil
	}
	db := ds.mongoClient.Database(ds.cfg.DatabaseName)
	coll := db.Collection(ds.cfg.BatchCollection)

	_, err := coll.InsertMany(ctx, documents)
	return err
}

func (ds *DataBaseService) InsertStream(ctx context.Context, document any) error {
	db := ds.mongoClient.Database(ds.cfg.DatabaseName)
	coll := db.Collection(ds.cfg.StreamCollection)

	_, err := coll.InsertOne(ctx, document)
	return err
}
