package main

import (
	"./model"
	"bytes"
	"encoding/binary"
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
var db *gorm.DB

func wsHandler(c *gin.Context) {
	conn, err := wsUpgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("Failed to set websocket upgrade: ", err.Error())
		return
	}
	clients[conn] = true
	err = conn.WriteJSON(getAvailableItem())
	if err != nil {
		conn.Close()
		delete(clients, conn)
	}
}

func broadcastItemStatistics(items []model.Item) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, items)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}
	for client := range clients {
		err = client.WriteJSON(items)
		if err != nil {
			fmt.Println("WriteJSON failed:", err)
		}
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

func getAllItems() []model.Item {
	items := make([]model.Item, 0)
	db.Table("items").Find(&items)
	return items
}

func getAvailableItem() []model.Item {
	items := make([]model.Item, 0)
	db.Table("items").Where("is_available=true").Find(&items)
	return items
}

func pickItem(itemId int64) model.PickingResult {
	item := model.Item{}
	db.Table("items").Where("id=?", itemId).Find(&item)
	pickedResult := item.IsAvailable
	db.Table("items").Model(&item).Update("is_available", false)
	go broadcastItemStatistics(getAvailableItem())
	return model.PickingResult{
		ItemId:        item.Id,
		ItemName:      item.Name,
		PickedSuccess: pickedResult,
	}
}

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
		context.JSONP(http.StatusOK, getAllItems())
	})

	router.GET("/picking/:itemId", func(context *gin.Context) {
		itemId, _ := strconv.ParseInt(context.Param("itemId"), 10, 32)
		context.JSONP(http.StatusOK, pickItem(itemId))
	})

	err := router.Run()
	if err != nil {
		log.Fatal("Server execute error: " + err.Error())
	}
}
