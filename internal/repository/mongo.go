package repository

import (
	"canonicalAuditlog/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	DBName         = "auditlogdb"
	CollectionName = "events"
	timeout        = 10 * time.Second // set a timeout of 5 seconds
)

type mongoEventRepository struct {
	client *mongo.Client
}

func NewMongoRepository(client *mongo.Client) EventRepository {
	return &mongoEventRepository{
		client: client,
	}
}

func (m *mongoEventRepository) Insert(event models.Event) error {
	collection := m.client.Database(DBName).Collection(CollectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := collection.InsertOne(ctx, event)
	return err
}

func (m *mongoEventRepository) Find(filter map[string]string) ([]models.Event, error) {
	collection := m.client.Database(DBName).Collection(CollectionName)

	query := bson.M{}
	for key, value := range filter {
		query[key] = value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cursor, err := collection.Find(ctx, query)
	if err != nil {
		return nil, err
	}

	var events []models.Event
	if err := cursor.All(ctx, &events); err != nil {
		return nil, err
	}

	return events, nil
}
