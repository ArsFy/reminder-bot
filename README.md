# Reminder Bot
![](https://img.shields.io/badge/license-MIT-blue)
![](https://img.shields.io/badge/GO-1.20-blue)
![](https://img.shields.io/badge/PRs-welcome-green)

透過 Telegram Bot 創建提醒事項

## How to Use


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

### 3. Run

```
./reminder-bot
```

-----

### 使用 Linux Service

https://github.com/ArsFy/add_service