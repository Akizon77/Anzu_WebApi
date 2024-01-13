package Types

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
type RSS struct {
	Title   string `json:"title"`
	SubLink string `json:"sub_link"`
}
type Update struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}
type RssSubs struct {
	User int64 `json:"user"`
	Rss  []RSS `json:"rss"`
}
type RssUpdates struct {
	User    int64    `json:"user"`
	Updates []Update `json:"updates"`
}
type POST_Add struct {
	User  int64  `json:"user"`
	Title string `json:"title"`
	Link  string `json:"link"`
}
type RssCache struct {
	Updates []Update `json:"updates"`
}
