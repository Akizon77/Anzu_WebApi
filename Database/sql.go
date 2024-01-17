package Database

import (
	"Anzu_WebApi/Config"
	"Anzu_WebApi/Log"
	"Anzu_WebApi/Messager"
	"Anzu_WebApi/Types"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	cfg *Config.Config
	db  *sql.DB
	log *Log.Logger
)

func init() {
	var err error
	log = Log.NewLogger("SQL")
	cfg = Config.GetAppConfig()
	log.Debug("Connecting to database ", fmt.Sprint(cfg.SQLServer.User, ":", cfg.SQLServer.Passwd, "@tcp("+cfg.SQLServer.Host, ":", cfg.SQLServer.Port, ")/", cfg.SQLServer.Name))
	db, _ = sql.Open("mysql", fmt.Sprint(cfg.SQLServer.User, ":", cfg.SQLServer.Passwd, "@tcp("+cfg.SQLServer.Host, ":", cfg.SQLServer.Port, ")/", cfg.SQLServer.Name))
	//验证数据库链接
	err = db.Ping()
	if err != nil {
		if strings.Contains(err.Error(), "Error 1045") {
			log.Error("Can not connect to database,incorrect user or password(1045)")
		} else {
			log.Error("Can not connect to database,", err)
		}
		return
	} else {
		log.Info("Successfully connected to database")
	}
	//初始化表头
	log.Debug("Initialize database header")

	query := fmt.Sprint("CREATE TABLE IF NOT EXISTS `", cfg.SQLServer.Name, "`.`RssSubs` (`user` BIGINT NOT NULL , `title` TEXT NOT NULL , `link` TEXT NOT NULL , `caches` JSON NOT NULL DEFAULT '' ) ENGINE = InnoDB;")
	_, err = db.Exec(query)
	if err != nil {
		log.Error("Can not create table,", err)
		log.Error("Can not initialize table")
		return
	}
	if err == nil {
		log.Debug("Initialize database header done")
	}

}
func GetAllSubs(userID int64) (*Types.RssSubs, error) {
	//准备返回数据
	result := &Types.RssSubs{
		User: userID,
		Rss:  []Types.RSS{},
	}
	//查询
	query, err := db.Query("SELECT * FROM `RssSubs` WHERE user = ?", userID)
	defer func() {
		if query != nil {
			query.Close()
		}
	}()
	//处理查询错误
	if err != nil {
		log.Error("Can not query table RssSubs,", err)
		return nil, err
	}
	//遍历查询
	for query.Next() {
		var (
			user   int64
			title  string
			link   string
			caches string
		)
		if err = query.Scan(&user, &title, &link, &caches); err != nil {
			log.Error("An error occurred while walking row,", err)
			//可能查询到一半遇到错误，先返回查到的数据，后面交给其他模块处理
			return result, err
		}
		//添加查询到的数据
		result.Rss = append(result.Rss, Types.RSS{
			Title:   title,
			SubLink: link,
		})
	}
	return result, nil
}
func GetUpdates(userID int64) (*Types.RssUpdates, error) {
	//准备数据
	result := &Types.RssUpdates{
		User:    userID,
		Updates: []Types.Update{},
	}
	//查询
	rows, err := db.Query("SELECT * FROM RssSubs WHERE user = ?", userID)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	//处理查询错误
	if err != nil {
		log.Error("Can not query table RssSubs,", err)
		return nil, err
	}
	//遍历查询
	for rows.Next() {
		var (
			user   int64
			title  string
			link   string
			caches string
		)
		if err = rows.Scan(&user, &title, &link, &caches); err != nil {
			log.Error("An error occurred while walking row,", err)
			//可能查询到一半遇到错误，先返回查到的数据，后面交给其他模块处理
			return result, err
		}
		var u []Types.Update
		err := json.Unmarshal([]byte(caches), &u)
		if err != nil {
			log.Error("An error occurred while walking row,", err)
			//可能查询到一半遇到错误，先返回查到的数据，后面交给其他模块处理
			return result, err
		}
		//添加查询到的数据
		for _, update := range u {
			result.Updates = append(result.Updates, update)
		}

	}
	return result, nil
}
func AddNew(userID int64, title string, link string) error {
	//检测数据是否已存在
	all, _ := GetAllSubs(userID)
	for _, rss := range all.Rss {
		if link == rss.SubLink {
			return errors.New("already exists")
		}
	}
	//插入数据
	q := fmt.Sprintf("INSERT INTO `RssSubs` (`user`, `title`, `link`, `caches`) VALUES ('%d', '%s', '%s', '{}')", userID, title, link)
	_, err := db.Exec(q)
	if err != nil {
		log.Error("An error occurred while insert row,", err)
		return err
	}
	return nil
}
func Del(userID int64, link string) error {
	_, err := db.Exec(fmt.Sprintf("DELETE FROM `RssSubs` WHERE user = '%d' AND link = '%s'", userID, link))
	if err != nil {
		log.Error("An error occurred while delete item,", err)
		return err
	}
	return nil
}
func GetCache(link string) ([]Types.Update, error) {
	result := []Types.Update{}
	rows, err := db.Query("SELECT * FROM RssSubs WHERE link = ?", link)
	if err != nil {
		return nil, err
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	for rows.Next() {
		var (
			user  int64
			title string
			d     string
			cache string
		)
		err := rows.Scan(&user, &title, &d, &cache)
		if err != nil {
			return result, err
		}
		err = json.Unmarshal([]byte(cache), &result)
		if err != nil {
			return result, err
		}
		break //只遍历第一个
	}

	return result, nil
}
func UpdateCache(link string, update []Types.Update) error {
	bytes, err := json.Marshal(update)
	if err != nil {
		return err
	}
	// 和旧的缓存比较，将更新的推送给用户
	cache, _ := GetCache(link)
	updates := findExtraElements(update, cache)
	//获取订阅了link的用户列表
	var users []int64
	rows, _ := db.Query("SELECT user,link FROM `RssSubs`")
	for rows.Next() {
		var (
			user int64
			li   string
		)
		rows.Scan(&li, &user)
		if li == link {
			users = append(users, user)
		}
	}
	go func(users []int64, updates []Types.Update) {
		for _, user := range users {
			for _, t := range updates {
				ticker := time.NewTicker(2 * time.Second)
				<-ticker.C
				Messager.TelegramPush(user, fmt.Sprint("#RSS ", t.Title, "\n", t.Link))
			}
		}
	}(users, updates)

	//更新缓存
	_, err = db.Exec("UPDATE `RssSubs` SET caches = ? WHERE link = ?", string(bytes), link)
	if err != nil {
		return err
	}
	return nil
}
func GetAllLinks() []string {
	var result []string
	rows, err := db.Query("SELECT DISTINCT link FROM `RssSubs`")
	if err != nil {
		log.Debug(err)
		return nil
	}
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	for rows.Next() {
		var link string

		err := rows.Scan(&link)
		if err != nil {
			log.Debug(err)
			return nil
		}
		result = append(result, link)
	}
	return result
}
func findExtraElements(a, b []Types.Update) []Types.Update {
	var extra []Types.Update

	for _, num := range a {
		found := false
		for _, n := range b {
			if num == n {
				found = true
				break
			}
		}
		if !found {
			extra = append(extra, num)
		}
	}
	return extra
}
