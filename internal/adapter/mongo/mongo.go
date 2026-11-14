package mongo

import (
	"context"
	"time"

	"stream-rec-core/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	usersCol        *mongo.Collection
	itemsCol        *mongo.Collection
	interactionsCol *mongo.Collection
}

func NewMongoRepository(ctx context.Context, db *mongo.Database) (*MongoRepository, error) {
	repo := &MongoRepository{
		usersCol:        db.Collection("users"),
		itemsCol:        db.Collection("items"),
		interactionsCol: db.Collection("interactions"),
	}
	if err := repo.ensureIndexes(ctx); err != nil {
		return nil, err
	}

	_ = repo.seedItems(ctx)

	return repo, nil
}

func (r *MongoRepository) seedItems(ctx context.Context) error {
	count, err := r.itemsCol.CountDocuments(ctx, bson.M{})
	if err != nil || count > 0 {
		return err
	}
	seed := []interface{}{
		domain.Item{ID: primitive.NewObjectID(), Name: "Inception", Description: "A thief steals corporate secrets through dream-sharing technology.", Category: "Sci-Fi", ReleaseYear: 2010, DurationMin: 148, Rating: 8.8},
		domain.Item{ID: primitive.NewObjectID(), Name: "Interstellar", Description: "Explorers travel through a wormhole in space to save humanity.", Category: "Sci-Fi", ReleaseYear: 2014, DurationMin: 169, Rating: 8.7},
		domain.Item{ID: primitive.NewObjectID(), Name: "The Dark Knight", Description: "Batman raises the stakes in his war on crime in Gotham.", Category: "Action", ReleaseYear: 2008, DurationMin: 152, Rating: 9.0},
		domain.Item{ID: primitive.NewObjectID(), Name: "Whiplash", Description: "A promising young drummer enrolls at a cut-throat music conservatory.", Category: "Drama", ReleaseYear: 2014, DurationMin: 106, Rating: 8.5},
	}
	_, err = r.itemsCol.InsertMany(ctx, seed)
	return err
}

func (r *MongoRepository) ensureIndexes(ctx context.Context) error {
	_, err := r.usersCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = r.itemsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "category", Value: 1}, {Key: "release_year", Value: -1}},
	})
	if err != nil {
		return err
	}

	_, err = r.interactionsCol.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "timestamp", Value: -1}},
	})
	return err
}

func (r *MongoRepository) CreateUser(ctx context.Context, user *domain.User) error {
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().UTC()
	_, err := r.usersCol.InsertOne(ctx, user)
	return err
}

func (r *MongoRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.usersCol.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	return &u, err
}

func (r *MongoRepository) GetUserByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var u domain.User
	err := r.usersCol.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	return &u, err
}

func (r *MongoRepository) CreateItem(ctx context.Context, item *domain.Item) error {
	item.ID = primitive.NewObjectID()
	_, err := r.itemsCol.InsertOne(ctx, item)
	return err
}

func (r *MongoRepository) GetItems(ctx context.Context) ([]domain.Item, error) {
	cursor, err := r.itemsCol.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var items []domain.Item
	err = cursor.All(ctx, &items)
	return items, err
}

func (r *MongoRepository) GetItemByID(ctx context.Context, id primitive.ObjectID) (*domain.Item, error) {
	var item domain.Item
	err := r.itemsCol.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	return &item, err
}

func (r *MongoRepository) UpdateItem(ctx context.Context, id primitive.ObjectID, item *domain.Item) error {
	_, err := r.itemsCol.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{
		"name":         item.Name,
		"description":  item.Description,
		"category":     item.Category,
		"release_year": item.ReleaseYear,
		"duration_min": item.DurationMin,
		"rating":       item.Rating,
	}})
	item.ID = id
	return err
}

func (r *MongoRepository) DeleteItem(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.itemsCol.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *MongoRepository) SearchItems(ctx context.Context, f domain.ItemSearchFilter) ([]domain.Item, error) {
	filter := bson.M{}
	if f.Query != "" {
		filter["name"] = bson.M{"$regex": f.Query, "$options": "i"}
	}
	if f.Category != "" {
		filter["category"] = bson.M{"$regex": "^" + f.Category + "$", "$options": "i"}
	}
	if f.Year > 0 {
		filter["release_year"] = f.Year
	}
	cursor, err := r.itemsCol.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var items []domain.Item
	err = cursor.All(ctx, &items)
	return items, err
}

func (r *MongoRepository) RecordInteraction(ctx context.Context, inter *domain.Interaction) error {
	inter.ID = primitive.NewObjectID()
	inter.Timestamp = time.Now().UTC()
	_, err := r.interactionsCol.InsertOne(ctx, inter)
	return err
}

func (r *MongoRepository) GetUserInteractions(ctx context.Context, userID primitive.ObjectID) ([]domain.Interaction, error) {
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}})
	cursor, err := r.interactionsCol.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}
	var history []domain.Interaction
	err = cursor.All(ctx, &history)
	return history, err
}

func (r *MongoRepository) GetRecommendations(ctx context.Context, userID primitive.ObjectID, limit int) ([]domain.Item, error) {
	return nil, nil
}
