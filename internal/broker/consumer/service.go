package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/urbaniakmichal/data-worker/internal/database"
)

type ConsumerService struct {
	kF *kafka.Reader
}

func NewConsumerService(kF *kafka.Reader) *ConsumerService {
	return &ConsumerService{
		kF: kF,
	}
}

func SetUpConsumer() *kafka.Reader {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"kafka1:9092"},
		Topic:          "data",
		GroupID:        "data-processor",
		MinBytes:       1,
		MaxBytes:       10e6,
		CommitInterval: 0,
	})

	return reader
}

func (cS *ConsumerService) StartListening(ctx context.Context, dbService *database.DataBaseService) {
	log.Println("Worker started listening to Kafka topic 'data'...")

	for {
		msg, err := cS.kF.FetchMessage(ctx)
		if err != nil {
			log.Printf("fetch error: %v", err)
			break
		}

		var payloads []database.EventDocument

		if err := json.Unmarshal(msg.Value, &payloads); err != nil {
			var singlePayload database.EventDocument
			if errSingle := json.Unmarshal(msg.Value, &singlePayload); errSingle != nil {
				log.Printf("failed to unmarshal message: %v", err)
				cS.kF.CommitMessages(ctx, msg)
				continue
			}

			if err := dbService.InsertStream(ctx, singlePayload); err != nil {
				log.Printf("failed to insert stream into mongo: %v", err)
				continue
			}
		} else {
			docs := make([]any, len(payloads))
			for i, v := range payloads {
				docs[i] = v
			}
			if err := dbService.InsertBatch(ctx, docs); err != nil {
				log.Printf("failed to insert batch into mongo: %v", err)
				continue
			}
		}

		if err := cS.kF.CommitMessages(ctx, msg); err != nil {
			log.Printf("failed to commit message: %v", err)
		} else {
			log.Printf("Successfully processed and committed message from partition %d at offset %d", msg.Partition, msg.Offset)
		}
	}
}
