# Reminder Bot
![](https://img.shields.io/badge/license-MIT-blue)
![](https://img.shields.io/badge/GO-1.20-blue)
![](https://img.shields.io/badge/PRs-welcome-green)

透過 Telegram Bot 創建提醒事項

## How to Use

- 列出所有提醒:  `/my_reminders`
- 新增提醒:     `/add_reminder` `2-26(可省略)` `11:20` `提醒文本`
- 新增每月提醒:  `/add_reminder_month` `日 (數字)` `11:20` `提醒文本`


## Deploy

### 1. Build
> 你可以選擇在 [Releases](https://github.com/ArsFy/reminder-bot/releases) 下載對應架構的二進制版本

```bash
git clone https://github.com/ArsFy/reminder-bot.git
cd reminder-bot
go mod tidy
go build .
```
得到二進制可執行檔 `reminder-bot`

### 2. <span id="bottoken">Get BotToken</span>
1. 創建 Telegram Bot: [@BotFather](https://t.me/BotFather)
2. 配置 Token，你可以透過下面 2 種方式之一進行配置
    - **Env (環境變數):** 設定環境變數名稱 `BOT_TOKEN`
    - **Config (配置檔案):** 在運行目錄下創建配置檔 `./token.conf`，內容為 Bot Token
3. *(可選) 使用 Env 設定時區
    - 設定環境變數名稱 `TIMEZONE`:
    - UTC+8: `8`, UTC-10: `-10` 

### 3. Run

```
./reminder-bot
```

-----

### Use Koyeb

Koyeb 支援託管 Golang 應用程式，點擊 Button 快速部署

[![](https://camo.githubusercontent.com/dbd49fd11e4dea39effabf3572eb66edafb50d32aadb31c7458fe7e42ac93790/68747470733a2f2f7777772e6b6f7965622e636f6d2f7374617469632f696d616765732f6465706c6f792f627574746f6e2e737667)](https://app.koyeb.com/deploy?type=git&repository=github.com/ArsFy/reminder-bot&branch=main&name=reminder-bot)

-----

### 使用 Linux Service

[<img src="https://opengraph.githubassets.com/0ce367d2a8cee652c1242cb4a99af11939ad2161e47eac849791a8695027a549/ArsFy/add_service" width="50%" style="border-radius: 5px" />](https://github.com/ArsFy/add_service)