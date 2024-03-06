package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type SignUpStruct struct {
	Name          string
	TelegramLogin string
	Password      string
}

var SignUpSlice = []SignUpStruct{} // ? empty
func main() {
	r := gin.Default()

	r.Use(Cors)
	r.POST("/signup", SignUp)
	go Recover()

	r.Run(":3434")
}
func Recover() {
	ReadUser()
	bot, _err := tgbotapi.NewBotAPI("6358252405:AAEdhEPf57Zl5pcSTu0nJroHJj22NP7MHgw")
	if _err != nil {
		fmt.Printf("error: %v\n", _err)
	}
	Update := tgbotapi.NewUpdate(0)
	allUpdates, UpdateError := bot.GetUpdatesChan(Update)

	for Update := range allUpdates {
		if Update.Message.IsCommand() {
			if Update.Message.Command() == "reset" {
				for _, item := range SignUpSlice {
					if item.TelegramLogin == Update.Message.Chat.userName {
						msg := tgbotapi.NewPhotoUpload(Update.Message.Chat.ID, tgbotapi.FileBytes)
						bot.Send(msg)
					}
				}

			}
		} else {
			// !  ==================================ikutiyum
			IsExist := false
			for index, item := range SignUpSlice {
				if item.TelegramLogin == Update.Message.Chat.userName {
					IsExist = true
					SignUpSlice[index].Password = Update.Message.Text
					msg := tgbotapi.NewPhotoUpload(Update.Message.Chat.ID, "THe password updated")
					bot.Send(msg)
				}

			}
			if !IsExist {
				msg := tgbotapi.NewPhotoUpload(Update.Message.Chat.ID, "the message is not found")
				bot.Send(msg)
			}

		}
	}
	WriteUser()

}

func SignUp(c *gin.Context) {
	var SignUpTemp SignUpStruct
	c.ShouldBindJSON(&SignUpTemp)

	if SignUpTemp.Name == "" || SignUpTemp.Password == "" || SignUpTemp.TelegramLogin == "" {
		c.JSON(404, "Empty field")
	} else {
		ReadUser()
		SignUpSlice = append(SignUpSlice, SignUpTemp)
		WriteUser()
	}
}
func WriteUser() {
	marshelData, _ := json.Marshal(SignUpSlice)
	ioutil.WriteFile("app.json", marshelData, 0644)
}
func ReadUser() {
	readedbyte, _ := ioutil.ReadFile("app.json")
	json.Unmarshal(readedbyte, &SignUpSlice)
}

func Cors(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://192.168.43.246:5500")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(200)
	}

	c.Next()
}
