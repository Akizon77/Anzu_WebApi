package Timer

import (
	"Anzu_WebApi/Database"
	"Anzu_WebApi/Log"
	"Anzu_WebApi/RSS"
	"fmt"
	"time"
)

func StartRssAutoRefresh(seconds int) {
	log := Log.NewLogger("Timer")
	duration, _ := time.ParseDuration(fmt.Sprint(seconds, "s"))
	ticker := time.NewTicker(duration)
	for {
		log.Debug("Updating rss")
		links := Database.GetAllLinks()
		for _, link := range links {
			update, err := RSS.GetUpdate(link)
			if err != nil {
				log.Error(err)
				continue
			}
			err = Database.UpdateCache(link, update)
			if err != nil {
				log.Error(err)
			}
		}
		log.Debug("Updating rss done")
		<-ticker.C
	}
}
func UpdateRssNow() {
	log := Log.NewLogger("Update")
	links := Database.GetAllLinks()
	log.Debug("Updating rss")
	for _, link := range links {
		go func(link string) {
			update, err := RSS.GetUpdate(link)
			if err != nil {
				log.Error(err)
				return
			}
			err = Database.UpdateCache(link, update)
			if err != nil {
				log.Error(err)
			}
		}(link)
	}
	log.Debug("Updating rss done")
}
