package worker

import (
	"context"

	"github.com/urbaniakmichal/data-worker/internal/broker/consumer"
	"github.com/urbaniakmichal/data-worker/internal/config"
	"github.com/urbaniakmichal/data-worker/internal/database"
)

func main() {
	cfg := config.LoadConfig("internal/config/config.yaml")
	ctx := context.Background()

	reader := consumer.SetUpConsumer()
	defer reader.Close()

	mongoClient := database.ConnectToMongo(cfg)
	defer mongoClient.Disconnect(ctx)

	dbService := database.NewDataBaseService(cfg, mongoClient)
	conService := consumer.NewConsumerService(reader)

	conService.StartListening(ctx, dbService)
}
