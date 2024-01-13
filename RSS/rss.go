package RSS

import (
	"Anzu_WebApi/Types"
	"github.com/mmcdole/gofeed"
)

func GetUpdate(link string) ([]Types.Update, error) {
	var result []Types.Update
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(link)
	if err != nil {
		return nil, err
	}
	for _, item := range feed.Items {
		temp := &Types.Update{
			Title: item.Title,
			Link:  item.Link,
		}

		result = append(result, *temp)

	}
	return result, nil
}
func contains(s []Types.Update, i Types.Update) bool {
	for _, x := range s {
		if x == i {
			return true
		}
	}
	return false
}
