package Config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token             string       `json:"token"`
	Debug             bool         `json:"debug"`
	Listen            string       `json:"listen"`
	RssUpdateInterval int          `json:"interval"`
	SQLServer         SQLServer    `json:"sql"`
	Telegram          TelegramPush `json:"telegram"`
}
type TelegramPush struct {
	EndPoint string `json:"end_point"`
	BotToken string `json:"bot_token"`
}
type SQLServer struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Name   string `json:"name"`
	User   string `json:"user"`
	Passwd string `json:"password"`
}

var (
	Instance = &Config{
		Token:             "leave blank to disable verify token",
		Debug:             false,
		Listen:            "0.0.0.0:8080",
		RssUpdateInterval: 600,
		SQLServer: SQLServer{
			Host:   "127.0.0.1",
			Port:   3306,
			Name:   "YOUR_DATABASE",
			User:   "YOUR_USERNAME",
			Passwd: "YOUR_PASSWORD",
		},
		Telegram: TelegramPush{
			EndPoint: "api.telegram.org",
			BotToken: "YOUR_BOT_TOKEN",
		},
	}
)

func init() {
	if _, err := os.Stat("./config.json"); os.IsNotExist(err) {
		b, err := json.Marshal(Instance)
		if err != nil {
			fmt.Println("Can not marshal config.json,using default config", err)
		}
		f, err := os.Create("./config.json")
		if err != nil {
			fmt.Println("Can not write config.json", err)
		}
		f.Write(b)
		defer f.Close()
		return
	}
	file, err := os.ReadFile("./config.json")
	if err != nil {
		fmt.Println("Can not read config.json,using default config", err)
	}
	err = json.Unmarshal(file, &Instance)
	if err != nil {
		fmt.Println("Can not unmarshal config.json,using default config", err)
	}
}
func GetAppConfig() *Config {
	return Instance
}
