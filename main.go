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

type Calender struct{
	Month int `json:month`
	Day int `json:day`
	Color int `json:color`
}
func main(){
	port := os.Getenv("PORT")
	ctx := context.Background()

	// Creates a client.
	_, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	router := gin.Default()
	router.GET("/",getIndex)
	router.POST("/sentiment",getSentiment)
	router.GET("/getCalendertest",getCalenderTest)
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
		c.String(http.StatusPreconditionFailed,"failed")

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

func getCalenderTest(c *gin.Context){
	var calender1 Calender = Calender{
		Month:11,
		Day:1,
		Color:12648430, //#c0ffee
	}
	var calender2 Calender = Calender{
		Month: 11,
		Day:   2,
		Color: 15789568, //f0ee00
	}// }
	var calender3 Calender = Calender{
		Month:11,
		Day:3,
		Color:11239568, //f0ee00
	}

	var calender []Calender
	calender = append(calender,calender1)
	calender = append(calender,calender2)
	calender = append(calender,calender3)
	//b, _ := json.Marshal(calender)
	//fmt.Printf("%s\n",b)
	c.JSON(200,calender)
}