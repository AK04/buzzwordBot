# Golang Slack bot
The bot scrapes Hackernews front and posts 2 links if a particular user send any message in the channel.

### Requirements
1. Go packages - Colly, Slack, Firestore (automatically installed with `go download`)
2. Firestore API key

### Setup
1. Create .env file with env.example
2. Download firestore service key file for the project and keep it in root folder.
3. Update the path to the service key in service/firebase_service.go file in line 17
4. go download
5. go build
6. ./buzzwordBot