package main

import (
	"context"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	vk "github.com/crossworth/vk-topic-to-json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var folder string
	var mongoDSN string
	var database string
	var collection string

	flag.StringVar(&folder, "folder", "backup", "Folder of JSON files")
	flag.StringVar(&mongoDSN, "mongo", "", "MongoDB DSN")
	flag.StringVar(&database, "database", "", "MongoDB Database")
	flag.StringVar(&collection, "collection", "", "MongoDB Collection")
	flag.Parse()

	if mongoDSN == "" {
		log.Fatalln("you must provide the MongoDSN")
	}

	if database == "" {
		log.Fatalln("you must provide the database")
	}

	if collection == "" {
		log.Fatalln("you must provide the collection")
	}

	files, err := filepath.Glob(folder + "/*.json")
	if err != nil {
		log.Fatalf("could not read the JSON files, %v", err)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoDSN))
	if err != nil {
		log.Fatalf("could not create mongodb client, %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("could not connect to mongodb, %v", err)
	}

	topicsCollection := client.Database(database).Collection(collection)

	for _, file := range files {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			log.Printf("could not read the file %s, %v\n", file, err)
			continue
		}

		var topic vk.Topic
		err = json.Unmarshal(content, &topic)
		if err != nil {
			log.Printf("could not decode the file %s, %v\n", file, err)
			continue
		}

		topicFilter := bson.D{{
			"id", topic.ID,
		}}

		var topicFromDB vk.Topic
		err = topicsCollection.FindOne(context.Background(), topicFilter).Decode(&topicFromDB)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Printf("error reading mongodb data, %v\n", err)
			continue
		}

		if topicFromDB.ID != 0 && topic.UpdatedAt == topicFromDB.UpdatedAt {
			log.Printf("topic %d already updated\n", topic.ID)
			continue
		}

		// update
		if topicFromDB.ID != 0 {
			_, err = topicsCollection.UpdateOne(context.Background(), topicFilter, topic)
			log.Printf("topic %d updated\n", topic.ID)
		} else {
			_, err = topicsCollection.InsertOne(context.Background(), topic)
			log.Printf("topic %d created\n", topic.ID)
		}
	}

	log.Printf("done")
}
