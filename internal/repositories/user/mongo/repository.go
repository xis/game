package mongo

import (
	"context"
	"fmt"

	"game/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrInvalidID = fmt.Errorf("%w, invalid record ID", domain.ErrInternal)
)

type MongoUserRepositoryDependencies struct {
	UsersCollection *mongo.Collection
}

type MongoUserRepository struct {
	usersCollection *mongo.Collection
}

func NewMongoUserRepository(deps MongoUserRepositoryDependencies) *MongoUserRepository {
	return &MongoUserRepository{
		usersCollection: deps.UsersCollection,
	}
}

func (repo *MongoUserRepository) Create(ctx context.Context, username, passwordHash string) (domain.User, error) {
	result, err := repo.usersCollection.InsertOne(ctx, bson.M{
		"username":     username,
		"passwordHash": passwordHash,
	})
	if err != nil {
		return domain.User{}, err
	}

	id, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return domain.User{}, ErrInvalidID
	}

	return domain.User{
		ID:           id.Hex(),
		Name:         username,
		PasswordHash: passwordHash,
	}, nil
}

func (repo *MongoUserRepository) GetByName(ctx context.Context, username string) (domain.User, error) {
	result := repo.usersCollection.FindOne(ctx, bson.M{
		"username": username,
	})
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return domain.User{}, domain.ErrResourceNotFound
		}

		return domain.User{}, result.Err()
	}

	var user userRecord

	err := result.Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:           user.ID.Hex(),
		Name:         user.Username,
		PasswordHash: user.PasswordHash,
	}, nil
}

func (repo *MongoUserRepository) CheckExistsByName(ctx context.Context, username string) (bool, error) {
	count, err := repo.usersCollection.CountDocuments(ctx, bson.M{
		"username": username,
	})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *MongoUserRepository) GetUsersByIDs(ctx context.Context, ids []string) ([]domain.User, error) {
	objectIDs := make([]primitive.ObjectID, len(ids))

	for i, id := range ids {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}

		objectIDs[i] = objectID
	}

	cursor, err := repo.usersCollection.Find(ctx, bson.M{
		"_id": bson.M{
			"$in": objectIDs,
		},
	})
	if err != nil {
		return nil, err
	}

	var users []domain.User

	for cursor.Next(ctx) {
		var userRecord userRecord

		err := cursor.Decode(&userRecord)
		if err != nil {
			return nil, err
		}

		users = append(users, domain.User{
			ID:           userRecord.ID.Hex(),
			Name:         userRecord.Username,
			PasswordHash: userRecord.PasswordHash,
		})
	}

	return users, nil
}
