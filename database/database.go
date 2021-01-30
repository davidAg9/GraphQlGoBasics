package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/davidAg9/go-graphql/env"
	"github.com/davidAg9/go-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//DB is database client structure
type DB struct {
	Client *mongo.Client
}

//Connect connects user the a local database server and exposes the graphql server resource
func Connect() *DB {
	env.Set()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	username := os.Getenv("dbUsername")
	pass := os.Getenv("dbPass")
	dbname := os.Getenv("dbname")
	uri := fmt.Sprintf("mongodb+srv://%v:%v@cluster0.1rozt.mongodb.net/%v?retryWrites=true&w=majority", username, pass, dbname)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("line 26:%+v", err)
	}

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("line 31:%+v", err)
	}
	return &DB{
		Client: client,
	}
}

//METHODS  FOR DATABASE OPERATIONS

//Save method accepts a Dog object and saves it into the database
func (db *DB) Save(input *model.NewDog) *model.Dog {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	collection := db.Client.Database("Animalia").Collection("dogs")
	res, err := collection.InsertOne(ctx, input)
	if err != nil {
		log.Fatalf("line 47:%+v", err)
	}

	return &model.Dog{
		ID:        res.InsertedID.(primitive.ObjectID).Hex(),
		Name:      input.Name,
		IsGoodBoy: input.IsGoodBoy,
	}

}

//GetByID accept a dogID and returns  a dog object with matching ID
func (db *DB) GetByID(ID string) *model.Dog {

	ObjID, err := primitive.ObjectIDFromHex(ID)

	if err != nil {
		log.Fatalf("%+v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	collection := db.Client.Database("Animalia").Collection("dogs")

	res := collection.FindOne(ctx, bson.M{"_id": ObjID})
	dog := model.Dog{}
	res.Decode(&dog)
	return &dog

}

//GetAllDogs return a slice containing all Dogs objects in database
func (db *DB) GetAllDogs() []*model.Dog {
	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)
	defer cancel()
	collection := db.Client.Database("Animalia").Collection("dogs")
	cur, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatalf("line 85:%+v", err)
	}
	defer cur.Close(ctx)
	var dogs []*model.Dog
	for cur.Next(ctx) {
		var dog *model.Dog
		err := cur.Decode(&dog)
		if err != nil {
			log.Fatalf("line 93:%+v", err)
		}
		dogs = append(dogs, dog)

	}
	return dogs
}
