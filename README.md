# TakakuraAnzu 后端API
## 配置文件
```json
{
    "token": "自定义，用于验证的字符串",
    "debug": false,
    "listen": "0.0.0.0:8080",
    "interval": 600,
    "sql": {
        "host": "YOUR_SQL_HOST",
        "port": 3306,
        "name": "YOUR_TABLE",
        "user": "YOUR_USERNAME",
        "password": "YOUR_PASSWD"
    },
    "telegram": {
        "end_point": "api.telegram.org",
        "bot_token": "YOUR_TELEGRAM_BOT_TOKEN"
    }
}
```
## API接口
### GET `/all` 获取订阅列表

参数：`token` `user`

返回示例

GET /all?token=ac48ed64124a?user=1
```json
{
  "code": 200,
  "message": "",
  "data": {
    "user": 1,
    "rss": [
      {
        "title": "Rss",
        "sub_link": "https://mikanani.me/RSS/MyBangumi?token=xxxxx"
      }
    ]
  }
}
```
### GET `/updates` 获取更新订阅链接

参数：`token` `user`

返回示例

GET /updates?token=ac48ed64124a?user=1
```json
{
  "code": 200,
  "message": "",
  "data": {
    "user": 1,
    "updates": [
      {
        "title": "【喵萌奶茶屋】★10月新番★[葬送的芙莉莲 / Sousou no Frieren][17][1080p][繁日双语][招募翻译]",
        "link": "https://mikanani.me/Home/Episode/4d12b34215085acebbcd021162ff12ad81b921a7"
      },
      {
        "title": "【喵萌奶茶屋】★10月新番★[葬送的芙莉莲 / Sousou no Frieren][17][1080p][简日双语][招募翻译]",
        "link": "https://mikanani.me/Home/Episode/cdbcf945e243817cbddc30d6d657d3f0968d08de"
      }
    ]
  }
}
```
### POST `/add` 添加订阅

参数：`token`

Body:
```json
{
  "user":132456,
  "title":"Title",
  "link":"https://example.com/"
}
```

返回示例

POST /add?token=ac48ed64124a
```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```
### POST `/del` 删除订阅

参数：`token`

Body:
```json
{
  "user":132456,
  "link":"https://example.com/"
}
```

返回示例

POST /del?token=ac48ed64124a
```json
{
  "code": 200,
  "message": "success",
  "data": null
}
```
