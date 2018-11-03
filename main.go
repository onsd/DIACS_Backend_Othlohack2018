package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"cloud.google.com/go/language/apiv1"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"os"
	"net/http"
)

type Emotion struct{
	Article string `json:article`
	Emotion float32 `json:emotion`
}

func main(){
	port := os.Getenv("PORT")

	router := gin.Default()
	router.GET("/",getIndex)
	router.POST("/sentiment",getSentiment)
	router.Run(":"+ port)
}

func getSentiment(c *gin.Context){
	c.Request.ParseForm()
	//title := c.Request.Form["title"]
	article := c.Request.Form["article"]

	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Detects the sentiment of the text.
	sentiment, err := client.AnalyzeSentiment(ctx, &languagepb.AnalyzeSentimentRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: article[0],
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	if err != nil {
		log.Fatalf("Failed to analyze text: %v", err)
	}

	//fmt.Printf("Text: %v\n", article)
	//fmt.Printf("Sentiment: %v" ,sentiment.DocumentSentiment.Score)

	//jsonStr := `{"channel":"` + channel + `","username":"` + name + `","text":"` + str(lending) + `"}`
	//var emotion Emotion
	emotion := &Emotion {
		Article : article[0],
		Emotion : sentiment.DocumentSentiment.Score,

	}
//	b ,_ := json.Marshal(emotion)

	c.JSON(200,emotion)
}


func getIndex(c *gin.Context){
	c.String(http.StatusAccepted,"hello")
}