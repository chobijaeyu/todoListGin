package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	dbname         = "dev"
	collectionname = "todo"
	ctxB           = context.Background()
)

type ToDo struct {
	ID       primitive.ObjectID "_id"
	Desc     string             `json:"Desc",bson:"Desc"`
	Imgs     string             `json:"Img",bson:"Img"`
	Done     bool               `json:"Done",bson:"Done"`
	Deadline string             `json:"Deadline",bson:"Deadline"`
}

func GetClient() (*mongo.Client, context.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Panicln("Couldn't connect to the database", err.Error())
	} else {
		log.Println("Connected!")
	}
	return client, ctx
}

func (t ToDo) AddRecord(db *mongo.Client, td ToDo) (r *mongo.InsertOneResult, err error) {
	colRef := db.Database(dbname).Collection(collectionname)
	r, err = colRef.InsertOne(ctxB, td)
	if err != nil {
		log.Printf("db add record err: %s\n", err.Error())
		return
	}
	return
}

func (t ToDo) LoadRecord(db *mongo.Client) (tlist []*ToDo, err error) {
	cur, err := db.Database(dbname).Collection(collectionname).Find(ctxB, bson.M{})
	if err != nil {
		log.Printf("Error on Finding all the documents:%s\n", err.Error())
		return
	}
	for cur.Next(context.TODO()) {
		var _t ToDo
		err = cur.Decode(&_t)
		if err != nil {
			log.Printf("Error on Decoding the document:%s\n", err.Error())
			return
		}
		tlist = append(tlist, &_t)
	}
	return
}

func (t ToDo) UpateRecord(db *mongo.Client, td ToDo) (r *mongo.UpdateResult, err error) {
	r, err = db.Database(dbname).Collection(collectionname).UpdateOne(ctxB, bson.M{"_id": t.ID}, td)
	if err != nil {
		log.Printf("db update record err: %s\n", err.Error())
		return
	}
	return
}

func (t ToDo) DeleteRecord(db *mongo.Client, id string) (r *mongo.DeleteResult, err error) {
	_id, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	r, err = db.Database(dbname).Collection(collectionname).DeleteOne(ctxB, bson.M{"_id": _id})
	if err != nil {
		log.Printf("db delete record err: %s\n", err.Error())
		return
	}
	return
}
