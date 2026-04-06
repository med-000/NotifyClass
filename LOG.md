# LOG

## 概要

このプロジェクトでは、`pkg/logger` 内でログローテーションを自動的に行っています。

主な出力先は以下です。

- `logs/app.log`
- `logs/Service/info.log`
- `logs/Service/warn.log`
- `logs/Service/error.log`
- `logs/repository/info.log`
- `logs/repository/warn.log`
- `logs/repository/error.log`
- `logs/notion/info.log`
- `logs/notion/warn.log`
- `logs/notion/error.log`
- `logs/discord/info.log`
- `logs/discord/warn.log`
- `logs/discord/error.log`
- `logs/parser/info.log`
- `logs/parser/warn.log`
- `logs/parser/error.log`
- `logs/scraper/info.log`
- `logs/scraper/warn.log`
- `logs/scraper/error.log`

## ローテーションルール

ログは以下の 2 つの条件でローテーションされます。

### 1. 日次ローテーション

- 日付が変わった後、最初の書き込み時にローテーションします
- 日付単位で最大 `7日分` 保持します

ローテーション後のファイル名は以下の形式です。

```text
logs/app.log.2026-04-06.01
logs/app.log.2026-04-06.02
```

ファイル名の日付は、ローテーション直前まで書き込まれていたログの日付です。

### 2. サイズローテーション

- 現在のログファイルが `5MB` を超える場合、その時点でローテーションします
- 同じ日付の分割ログは最大 `5世代` まで保持します

これにより、特定の日に大量ログが出てもログサイズが無制限に増えないようにしています。

## 実際の挙動

- 通常日は、1日につき 1 つのローテーションファイルが作られます
- ログ出力量が多い日は、同じ日付で複数ファイルに分割されます
- 古い日付のログは、新しい日付を `7日分` 残して削除されます
- 同一日内で分割されたログは、新しいものを `5個` 残して削除されます

## 補足

- ローテーションは Docker や OS の機能ではなく、Go の実装で行っています
- 現在書き込み中のファイル名は固定です
  例: `logs/app.log`
- ローテーション済みファイルは退避用であり、直接書き込み対象にはしません
- 日付が変わってもログ出力がなければ、その時点ではローテーションされません
  次の書き込み時にまとめて判定されます
