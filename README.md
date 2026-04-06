# **NotifyClass**

大学の LMS から講義・資料・課題情報を取得して、  
DB と Notion に同期しつつ、Discord から未完了タスクを確認できるアプリです。

## **Overview**

- 学校の LMS は課題や資料の確認が散らばりやすく、締切管理もしづらかったため作成しました。
- 取得した情報を MySQL に保存し、Notion に同期し、Discord からすぐ見られるようにしています。
- Notion 側で完了チェックを付けた内容を DB に戻すこともできます。

## **Architecture**

### **Tech Stack**

- Backend: Go
- Scraping: Colly / goquery
- Database: MySQL / GORM
- Notion integration: Notion API
- Bot: DiscordGo
- Container: Docker / Docker Compose
- Scheduler: robfig/cron

### **Container**

- `db`
  MySQL
- `backend`
  定期同期と `/sync` の手動実行を受け持つ常駐 worker
- `discord`
  Discord Bot

## **Project Structure**

```text
.
├── cmd
│   ├── app
│   ├── discord
│   ├── fetch
│   ├── notionpull
│   ├── notionpush
│   └── repository
├── db
├── pkg
│   ├── appflow
│   ├── discord
│   ├── logger
│   ├── notion
│   ├── parser
│   ├── repository
│   ├── scraping
│   └── service
├── Dockerfile
├── compose.yaml
├── LOG.md
└── README.md
```

## **Main Flow**

`backend` が行う基本の一連処理は以下です。

1. Notion から `完了` 状態を DB に反映
2. 学校サイトから講義情報を取得
3. JSON に export
4. DB に保存
5. DB の内容を Notion に push

この処理は、

- 毎日決まった時間の自動実行
- Discord の `/sync`

のどちらでも起動できます。

## **Commands**

### **各コマンドの役割**

- `cmd/app`
  定期実行と手動 sync を受け持つ backend worker
- `cmd/discord`
  Discord Bot
- `cmd/fetch`
  学校サイトから fetch して `course.json` に export
- `cmd/repository`
  学校サイトから fetch して DB 保存まで
- `cmd/notionpull`
  Notion から DB に `完了` 状態を反映
- `cmd/notionpush`
  DB から Notion に push

### **ローカル実行**

1. `.env` を作成
   `.env.template` を元に設定してください

2. backend を起動

```bash
go run ./cmd/app
```

3. Discord Bot を起動

```bash
go run ./cmd/discord
```

### **単体実行**

_fetch only_

```bash
go run ./cmd/fetch
```

_fetch + DB save_

```bash
go run ./cmd/repository
```

_Notion pull_

```bash
go run ./cmd/notionpull
```

_Notion push_

```bash
go run ./cmd/notionpush
```

## **Docker**

### **Run**

```bash
docker compose up -d
```

起動するもの:

- MySQL
- backend worker
- Discord Bot

## **Discord Bot**

### **Usage**

- Bot をメンションすると、未完了タスク一覧を表示
- `/2025010102` のように送ると、指定した年度・学期・曜日・時限の未完了タスクを表示
- `/sync` で backend に同期処理を依頼

`/sync` 実行中に再度コマンドやメンションを送った場合は、現在処理中であることを返します。

## **Notion**

### **Sync**

- DB -> Notion
  `cmd/notionpush`
- Notion -> DB
  `cmd/notionpull`

Event DB の `完了` checkbox を DB の `is_done` に反映します。

## **Log**

ログについては [LOG.md](/Users/med/Documents/App_dev/product/NotifyClass/LOG.md) にまとめています。

- 日次ローテーション
- 5MB 超過時ローテーション
- 世代保持数

を記載しています。

## **Environment Variables**

主に必要なものは以下です。

- MySQL 接続情報
- LMS のログイン情報
- Discord Bot Token
- Discord 通知先 Channel ID
- Notion API Key / Database ID
- `APP_SYNC_SCHEDULE`

例:

```env
APP_SYNC_SCHEDULE=0 2,13 * * *
```

毎日 `02:00` と `13:00` に同期を行います。

## **Technical Highlights**

- 学校サイトの情報取得を scraping / parser / repository / notion / discord で責務分離しています
- Notion は push / pull の両方向同期に対応しています
- Discord から未完了タスク確認と手動同期ができます
- backend は定期実行と手動実行を同じ flow で扱えるようにしています
- ログはアプリ内で自動ローテーションするようにしています

## **License**

MIT
