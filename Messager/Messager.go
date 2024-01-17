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
var failedMessages []FailedMessage

func init() {
	config = Config.GetAppConfig()
	log = Log.NewLogger("PUSH")
}

func GetFailedMessages() []FailedMessage {
	return failedMessages
}

type FailedMessage struct {
	Chatid  int64
	Message string
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
		AddFailed(FailedMessage{
			Chatid:  chatid,
			Message: message,
		})
		log.Error(err)
	}
	ResendMessage()
}
func AddFailed(msg FailedMessage) {
	failedMessages = append(failedMessages, msg)
}
func Remove(msg FailedMessage) {
	failedMessages = nil
	for _, message := range failedMessages {
		if message.Message == msg.Message && message.Chatid == msg.Chatid {
			continue
		}
		failedMessages = append(failedMessages, message)
	}
}
func ResendMessage() {
	if len(failedMessages) > 200 {
		log.Error("未发送的消息达", len(failedMessages))
		return
	}
	url := fmt.Sprint("https://", config.Telegram.EndPoint, "/bot", config.Telegram.BotToken, "/sendMessage")
	for _, msg := range failedMessages {
		bytes, _ := json.Marshal(struct {
			ChatId int64  `json:"chat_id"`
			Text   string `json:"text"`
		}{msg.Chatid, msg.Message})
		log.Info("RESending ", msg.Chatid, ": ", msg.Message)
		_, err := http.Post(url, "application/json", strings.NewReader(string(bytes)))
		if err != nil {
			log.Error("Failed to resent")
		}
		Remove(msg)
	}

}
