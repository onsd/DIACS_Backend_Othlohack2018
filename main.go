package main

import (
	"cloud.google.com/go/language/apiv1"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
	"log"
	"net/http"
	"database/sql"
	"strconv"
	"math/rand"
)

var LastEmotion *Emotion
type Emotion struct {
	UserName   string  `json:username`
	Article    string  `json:article`
	Emotion    float32 `json:emotion`
	EmotionNum int     `json:emotionNum`
	ColorCode int `json:colorCode`
	Year int	`json:year`
	Month int	`json:month`
	Day int	`json:day`
}

type Calender struct {
	Month int `json:month`
	Day   int `json:day`
	Color int `json:color`
	EmotionNum int `json:emotionNum`
}

func getColor(f float32) int {
	if -1.0 < f && f <= -0.75 {
		return 0
	}
	if -0.75 < f && f <= -0.5 {
		return 1
	}
	if -0.5 < f && f < -0.25 {
		return 2
	}
	if -0.25 < f && f < 0.0 {
		return 3
	}

	if 0.0 < f && f <= 0.25 {
		return 5
	}
	if 0.25 < f && f <= 0.5 {
		return 6
	}
	if 0.5 < f && f <= 0.75 {
		return 7
	}
	if 0.75 < f && f <= 1.0 {
		return 8
	}
	return 4 //f == 0
}

func main() {
	initDB()

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/", getIndex)
	router.POST("/getSentiment", getSentiment)
	router.GET("/getCalendertest", getCalenderTest)
	router.POST("/postforexample",postforexample)
	//router.GET("/getCalender/:name",getCalender)
	router.GET("/lastState",getLastState)
	router.Run(":8080")
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		c.Next()
	}
}
func getSentiment(c *gin.Context) {
	Color := []int{0x9999ff,0x99ccff,0x99ffff,0x99ffcc,0x99ff99,0xccff99,0xffff99,0xffcc99,0xff9999}

	c.Request.ParseForm()
	//title := c.Request.Form["title"]
	article := c.Request.Form["article"]
	username := c.Request.Form["username"]
	y := c.Request.Form["year"]
	d := c.Request.Form["day"]
	m := c.Request.Form["month"]
	year,_ := strconv.Atoi(y[0])
	day,_ := strconv.Atoi(d[0])
	month,_ := strconv.Atoi(m[0])


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
		c.String(http.StatusPreconditionFailed, "failed")
	}

	emotionNum := getColor(sentiment.DocumentSentiment.Score)
	emotion := &Emotion{
		Article:    article[0],
		Emotion:    sentiment.DocumentSentiment.Score,
		EmotionNum: emotionNum,
		UserName:   username[0],
		ColorCode: Color[emotionNum],
		Year: year,
		Month : month,
		Day : day,
	}

	db,err := sql.Open("sqlite3","./test.db")
	_, err = db.Exec(
		`insert into dairy (user,year,month,day,article,emotionNum,colorCode)
				values (?,?,?,?,?,?,?)`,
			emotion.UserName,emotion.Year,emotion.Month,emotion.Day, emotion.Article,emotion.EmotionNum,emotion.ColorCode)
	if err != nil{
		log.Fatalf("Error : %v",err)
	}

	LastEmotion = emotion
	c.JSON(200, emotion)
}

func getIndex(c *gin.Context) {
	c.String(http.StatusAccepted, "hello")
}

func getCalenderTest(c *gin.Context) {
	var calender1 Calender = Calender{
		Month: 11,
		Day:   1,
		Color: 12648430, //#c0ffee
		EmotionNum: 5,
	}
	var calender2 Calender = Calender{
		Month: 11,
		Day:   2,
		Color: 15789568, //f0ee00
		EmotionNum: 4,
	} // }
	var calender3 Calender = Calender{
		Month: 11,
		Day:   3,
		Color: 11239568, //f0ee00
		EmotionNum: 6,
	}

	var calender []Calender
	calender = append(calender, calender1)
	calender = append(calender, calender2)
	calender = append(calender, calender3)
	//b, _ := json.Marshal(calender)
	//fmt.Printf("%s\n",b)
	c.JSON(200, calender)
}

func initDB() {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Fatalf("Connection Error: %v", err)
	}
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS "Dairy" ("user" string,"year" int,"month" int,"day" int,"article" string,"emotionNum" int,"colorCode" int);
	`)
	if err != nil {
		log.Fatalf("Connection Error: %v", err)
	}
	db.Close()

}

func postforexample(c *gin.Context){
	Color := []int{0x9999ff,0x99ccff,0x99ffff,0x99ffcc,0x99ff99,0xccff99,0xffff99,0xffcc99,0xff9999}

	c.Request.ParseForm()
	//title := c.Request.Form["title"]
	article := c.Request.Form["article"]
	username := c.Request.Form["username"]
	y := c.Request.Form["year"]
	d := c.Request.Form["day"]
	m := c.Request.Form["month"]
	year,_ := strconv.Atoi(y[0])
	day,_ := strconv.Atoi(d[0])
	month,_ := strconv.Atoi(m[0])



	randint := rand.Int()
	if 9 <= randint{
		randint = 1
	}
	randfloat := rand.Float32()

	emotion := &Emotion{
		Article:    article[0],
		Emotion:    randfloat,
		EmotionNum: randint,
		UserName:   username[0],
		ColorCode: Color[randint],
		Year: year,
		Month : month,
		Day : day,
	}

	db,err := sql.Open("sqlite3","./test.db")
	defer db.Close()
	_, err = db.Exec(
		`insert into dairy (user,year,month,day,article,emotionNum,colorCode)
				values (?,?,?,?,?,?,?)`,
		emotion.UserName,emotion.Year,emotion.Month,emotion.Day, emotion.Article,emotion.EmotionNum,emotion.ColorCode)
	if err != nil{
		log.Fatalf("Error : %v",err)
	}


	LastEmotion = emotion
	c.JSON(200, emotion)
}


func getCalender(c *gin.Context){
	username := c.Param("username")
	db, err := sql.Open("sqlite3", "./test.db")
	if err !=  nil{
		log.Fatalf("Open sqlite3 Error: %v",err)
	}
	res, err := db.Exec("SELECT * FROM DAIRY WHERE user = ?",username)


}

func getLastState(c *gin.Context){
	c.JSON(200,LastEmotion)
}