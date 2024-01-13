package Log

import (
	"Anzu_WebApi/Config"
	"fmt"
	"log"
	"os"
	"time"
)

var cfg *Config.Config

func init() {
	cfg = Config.GetAppConfig()
}

var instance = &Logger{
	"Main",
	log.New(os.Stdout, "", 0),
}
var web = &Logger{
	"HTTP",
	log.New(os.Stdout, "", 0),
}

type Logger struct {
	Module string
	Logger *log.Logger
}

func NewLogger(module string) *Logger {
	return &Logger{
		module,
		log.New(os.Stdout, "", 0),
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if !cfg.Debug {
		return
	}
	l.log(l.Module, "Debug", v, l.Logger)
}

func (l *Logger) Warn(v ...interface{}) {
	l.log(l.Module, "Warn", v, l.Logger)
}

func (l *Logger) Info(v ...interface{}) {
	l.log(l.Module, "Info", v, l.Logger)
}

func (l *Logger) Error(v ...interface{}) {
	l.log(l.Module, "Error", v, l.Logger)
}
func (l *Logger) Panic(v ...interface{}) {
	l.log(l.Module, "Panic", v, l.Logger)
}

func (l *Logger) log(module string, level string, v []interface{}, logger *log.Logger) {
	timestamp := time.Now().Format(time.DateTime)
	logMessage := fmt.Sprintf("[%s][%s/%s]%s", timestamp, module, level, fmt.Sprint(v...))
	if level == "Panic" {
		logger.Panic(logMessage)
	} else {
		logger.Println(logMessage)
	}
}
func Debug(v ...any) {
	instance.Debug(fmt.Sprint(v...))
}
func Info(v ...any) {
	instance.Info(fmt.Sprint(v...))
}
func Warn(v ...any) {
	instance.Warn(fmt.Sprint(v...))
}
func Error(v ...any) {
	instance.Error(fmt.Sprint(v...))
}
func Panic(v ...any) {
	instance.Panic(fmt.Sprint(v...))
}
func GetWebLogger() *Logger {
	return web
}
