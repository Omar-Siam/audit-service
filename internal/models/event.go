package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	EventID      primitive.ObjectID     `json:"event_id" bson:"_id,omitempty"`
	EventType    string                 `json:"event_type" bson:"event_type"`
	Timestamp    time.Time              `json:"timestamp" bson:"timestamp"`
	ServiceName  string                 `json:"service_name" bson:"service_name"`
	CustomerID   string                 `json:"customer_id" bson:"customer_id"`
	EventDetails map[string]interface{} `json:"event_details" bson:"event_details"`
}
