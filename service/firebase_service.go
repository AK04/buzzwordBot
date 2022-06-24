package service

import (
	"buzzwordBot/models"
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"

	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func InitFireStoreClient(ctx context.Context) (*firestore.Client, error) {
	opt := option.WithCredentialsFile("PathToFirestoreKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	return client, nil
}

func AddItemsToFirestore(client *firestore.Client, ctx context.Context, items []models.Item) {
	for _, item := range items {
		_, err := client.Collection("buzzwordItems").Doc(item.Id).Set(ctx, item)
		if err != nil {
			log.Printf("Failed adding aturing: %v\n", err)
		}
	}
}

func GetItemIdsFromFirestore(client *firestore.Client, ctx context.Context) []string {
	var itemIds []string
	iter := client.Collection("buzzwordItems").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		itemId := doc.Data()["Id"].(string)
		itemIds = append(itemIds, itemId)
	}
	return itemIds
}
