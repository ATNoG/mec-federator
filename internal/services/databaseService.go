package services

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DatabaseServiceInterface interface {
}

type DatabaseService struct {
	mongoClient *mongo.Client
}

func NewDatabaseService(mongoClient *mongo.Client) *DatabaseService {
	return &DatabaseService{
		mongoClient: mongoClient,
	}
}

// fetches entities list of a specified type
func FetchEntitiesFromDatabase[T any](collection *mongo.Collection, filter interface{}) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error. could not find entities from database: %s", err)
	}
	defer cursor.Close(ctx)

	var results []T
	for cursor.Next(ctx) {
		var entity T
		if err := cursor.Decode(&entity); err != nil {
			return nil, fmt.Errorf("error. could not decode entity from database: %s", err)
		}
		results = append(results, entity)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error. database cursor error: %s", err)
	}

	return results, nil
}

// fetches a single entity of a specified type
func FetchEntityFromDatabase[T any](collection *mongo.Collection, filter interface{}) (T, error) {
	var result T
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("error. could not fetch entity from database: %s", err)
	}

	return result, nil
}

// saves an entity to the database, returning the id of the entity
func SaveEntityToDatabase[T any](collection *mongo.Collection, entity T) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("error. could not save entity to database: %s", err)
	}

	return result.InsertedID.(string), nil
}
