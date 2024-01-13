package Messager

import (
	"Anzu_WebApi/Config"
	"Anzu_WebApi/Log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var config *Config.Config
var log *Log.Logger

func init() {
	config = Config.GetAppConfig()
	log = Log.NewLogger("PUSH")
}
func TelegramPush(chatid int64, message string) {

	url := fmt.Sprint("https://", config.Telegram.EndPoint, "/bot", config.Telegram.BotToken, "/sendMessage")
	bytes, _ := json.Marshal(struct {
		ChatId int64  `json:"chat_id"`
		Text   string `json:"text"`
	}{chatid, message})
	log.Info("Sending ", chatid, ": ", message)
	_, err := http.Post(url, "application/json", strings.NewReader(string(bytes)))
	if err != nil {
		log.Error(err)
	}
}
