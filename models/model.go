package models

import (
	"context"
	"log"
	"strconv"
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
	// ID       string `json:"ID",bson:"ID"`
	Desc     string             `json:"Desc" bson:"Desc"`
	Imgs     string             `json:"Img" bson:"Img"`
	Done     bool               `json:"Done" bson:"Done"`
	Deadline string             `json:"Deadline" bson:"Deadline"`
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
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
	update := bson.D{{"$set", td}}
	r, err = db.Database(dbname).Collection(collectionname).UpdateOne(ctxB, bson.M{"_id": t.ID}, update)
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
	// r, err = db.Database(dbname).Collection(collectionname).DeleteOne(ctxB, bson.M{"ID": id})
	r, err = db.Database(dbname).Collection(collectionname).DeleteOne(ctxB, bson.M{"_id": _id})
	if err != nil {
		log.Printf("db delete record err: %s\n", err.Error())
		return
	}
	return
}

func (t ToDo) QueryRecord(db *mongo.Client, q, p string) (tlist []*ToDo, err error) {
	log.Println(q, p)
	_p, err := strconv.ParseBool(p)
	if err != nil {
		log.Println("strconv.ParseBool(p) failure")
		return nil, err
	}
	cur, err := db.Database(dbname).Collection(collectionname).Find(ctxB, bson.M{q: _p})
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

func (t ToDo) WsRecord(db *mongo.Client, docChan chan []byte) {
	var (
		ctx = context.TODO()
	)

	// matchStage := bson.D{{"$match", bson.D{{"operationType", "insert"}}}}
	opts := options.ChangeStream().SetMaxAwaitTime(2 * time.Second).SetFullDocument(options.UpdateLookup)
	cs, err := db.Database(dbname).Collection(collectionname).Watch(ctx, mongo.Pipeline{}, opts)
	if err != nil {
		log.Println("cs err:", err.Error())
		return
	}
	defer cs.Close(ctx)
	defer close(docChan)
	for cs.Next(context.TODO()) {
		// var _t ToDo
		// err := cs.Current.Lookup("fullDocument").Unmarshal(&_t)
		// if err != nil {
		// 	log.Println("cs json err:", err.Error())
		// 	log.Println(cs.Current.String())
		// 	continue
		// }
		// inrec, _ := json.Marshal(&_t)
		docChan <- []byte(cs.Current.String())
	}
}
