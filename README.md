# myRecord

## Development

### run locally
```sh
docker run -d \
  --name pg-dev \
  -e POSTGRES_PASSWORD=c7u0lOx7cmKfNx5m \
  -e POSTGRES_DB=record \
  -p 127.0.0.1:5432:5432 \
  postgres:15-bullseye

go run main.go -db postgresql://postgres:c7u0lOx7cmKfNx5m@127.0.0.1/record
```

### gitmoji
[https://gitmoji.dev](https://gitmoji.dev)

| type    | code               | emoji |
| ------- | ------------------ | ----- |
| ADD     | :heavy_plus_sign:  | ➕     |
| BREAK   | :boom:             | 💥     |
| CHORE   | :hammer:           | 🔨     |
| DEL     | :fire:             | 🔥     |
| DEP     | :package:          | 📦     |
| DOC     | :memo:             | 📝     |
| FEAT    | :sparkles:         | ✨     |
| FIX     | :bug:              | 🐛     |
| INIT    | :tada:             | 🎉     |
| OPT     | :zap:              | ⚡     |
| REF     | :recycle:          | ♻️     |
| TEST    | :white_check_mark: | ✅     |
| VERSION | :bookmark:         | 🔖     |
| WIP     | :construction:     | 🚧     |
