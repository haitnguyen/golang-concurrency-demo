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
	//db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=gorm password=123456 sslmode=disable search_path=myschema")
	db, err := gorm.Open("postgres", "host=localhost port=5432 user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatal("Server execute error: " + err.Error())
	}
	return db
}

func getAllItems(db *gorm.DB) []model.Item {
	items := make([]model.Item, 0)
	db.Table("items").Where("is_available=true AND amount > 0").Find(&items)
	return items
}

func pickItem(db *gorm.DB, itemId int64) model.PickingResult {
	item := model.Item{}
	db.Table("items").Where("id=?", itemId).Find(&item)
	remainingItem := item.Amount
	item.Amount = remainingItem - 1
	go db.Table("items").Save(&item)
	return model.PickingResult{
		ItemId:        item.Id,
		ItemName:      item.Name,
		PickedSuccess: item.IsAvailable && remainingItem >= 0,
	}
}

var db *gorm.DB

func main() {
	db = connectDB()
	db.AutoMigrate(&model.Item{})
	router := gin.Default()
	router.LoadHTMLFiles("template/running-square.html")
	router.Static("/static", "./template/static")
	router.GET("/", func(context *gin.Context) {
		context.HTML(200, "running-square.html", nil)
	})
	router.GET("/ws", wsHandler)

	router.GET("/items", func(context *gin.Context) {
		context.JSONP(http.StatusOK, getAllItems(db))
	})

	router.GET("/picking/:itemId", func(context *gin.Context) {
		itemId, _ := strconv.ParseInt(context.Param("itemId"), 10, 32)

		context.JSONP(http.StatusOK, pickItem(db, itemId))
	})

	err := router.Run()
	if err != nil {
		log.Fatal("Server execute error: " + err.Error())
	}
}
