package app

import (
	"buzzwordBot/common"
	"buzzwordBot/models"
	"buzzwordBot/service"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"cloud.google.com/go/firestore"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
)

func Start() {
	godotenv.Load(".env")

	evilUser := os.Getenv("EVIL_USER")
	ownerUser := os.Getenv("OWNER_USER")

	botState := false

	slackClient := service.InitSlackClient()
	socketClient := service.InitSocketClient(slackClient)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	firestoreClient, err := service.InitFireStoreClient(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore client: %v\n", err)
	}

	go func(ctx context.Context, client *slack.Client, socketClient *socketmode.Client) {
		for {
			select {
			case <-ctx.Done():
				log.Println("Shutting down socketmode listener")
				return
			case event := <-socketClient.Events:
				switch event.Type {
				case socketmode.EventTypeEventsAPI:
					eventsAPIEvent, ok := event.Data.(slackevents.EventsAPIEvent)
					if !ok {
						log.Printf("Could not type cast the event to the EventsAPIEvent: %v\n", event)
						continue
					}
					socketClient.Ack(*event.Request)
					data := eventsAPIEvent.InnerEvent.Data.(*slackevents.MessageEvent)
					if data.User == ownerUser {
						if data.Text == models.COMMAND_START {
							botState = botStart(slackClient)
						} else if data.Text == models.COMMAND_STOP {
							botState = botStop(slackClient)
						}
					}
					if data.User == evilUser && botState {
						botState = sendItemToChannel(firestoreClient, slackClient, ctx)
					}
				}
			}
		}
	}(ctx, slackClient, socketClient)

	socketClient.Run()

	fmt.Printf("\ndone am I\n")
}

func botStart(slackClient *slack.Client) bool {
	service.UploadFile(slackClient, "awake", "assets/wakeup"+strconv.Itoa(rand.Intn(3))+".gif")
	return true
}

func botStop(slackClient *slack.Client) bool {
	service.UploadFile(slackClient, "awake", "assets/sleep1.gif")
	return false
}

func sendItemToChannel(firestoreClient *firestore.Client, slackClient *slack.Client, ctx context.Context) bool {
	items := service.ScrapeHackerNewsLinks()

	itemIds := service.GetItemIdsFromFirestore(firestoreClient, ctx)
	newItems := common.RemoveDuplicateItems(items, itemIds)
	fmt.Printf("Removed %d duplicate items\n", len(items)-len(newItems))

	itemCountCutoff := 2
	if len(newItems) < 3 {
		itemCountCutoff = len(newItems)
	}

	for _, item := range newItems[0:itemCountCutoff] {
		service.SendSlackMessage(slackClient, item.Link)
	}

	if len(newItems) == 0 {
		service.SendSlackMessage(slackClient, "<@OwneruserID> I am out of links ")
		fmt.Printf("No new items\n")
		return botStop(slackClient)
	} else {
		service.AddItemsToFirestore(firestoreClient, ctx, newItems[0:itemCountCutoff])
		fmt.Printf("Firestore Upload complete\n")
		return true
	}
}
