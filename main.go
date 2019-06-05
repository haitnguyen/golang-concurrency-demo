package main

import (
	"./model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strconv"
)

var wsUpgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool)

func wsHandler(c *gin.Context) {
	conn, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Failed to set websocket upgrade: ", err.Error())
		return
	}
	clients[conn] = true
	for {
		t, msg, err := conn.ReadMessage()
		fmt.Println(string(msg))
		if err != nil {
			break
		}
		go broadcastMsg(conn, msg, t)
	}
}

func broadcastMsg(conn *websocket.Conn, msg []byte, msgType int) {
	for client := range clients {
		if conn != client {
			go writeMsg(client, msg, msgType)
		}
	}
}

func writeMsg(client *websocket.Conn, msg []byte, msgType int) {
	err := client.WriteMessage(msgType, msg)
	if err != nil {
		client.Close()
		delete(clients, client)
	}
}

func connectDB() *gorm.DB {
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gorm password=123456 sslmode=disable search_path=myschema")
	if err != nil {
		log.Fatal("Server execute error: " + err.Error())
	}
	return db
}

func getAllItems(db *gorm.DB) []model.Item {
	items := make([]model.Item, 0)
	db.Table("item").Find(&items)
	return items
}

func pickItem(db *gorm.DB) []int {
	ids := make([]int, 0)
	db.Table("item").Pluck("id", &ids)
	return ids
}

func getItems() []model.Item {
	item1 := model.Item{
		Id:   1,
		Name: "abc",
	}

	item2 := model.Item{
		Id:   2,
		Name: "cdf",
	}

	return []model.Item{item1, item2}
}

var db *gorm.DB

func main() {
	db = connectDB()
	db.AutoMigrate(&model.Item{})
	//getAllItems(db)
	//pickItem(db)
	//homepageViewPath := "home.html"
	router := gin.Default()
	router.LoadHTMLFiles("template/running-square.html")
	router.Static("/static", "./template/static")
	router.GET("/", func(context *gin.Context) {
		context.HTML(200, "running-square.html", nil)
	})
	router.GET("/ws", wsHandler)

	router.GET("/items", func(context *gin.Context) {
		context.JSONP(http.StatusOK, getItems())
	})

	router.GET("/picking/:itemId", func(context *gin.Context) {
		itemId, _ := strconv.ParseInt(context.Param("itemId"), 10, 32)

		context.JSONP(http.StatusOK, checkingResult(itemId))
	})

	err := router.Run()
	if err != nil {
		log.Fatal("Server execute error: " + err.Error())
	}
}

func checkingResult(itemId int64) bool {
	return itemId%2 == 0
}
