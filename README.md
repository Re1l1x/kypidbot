# kypidbot

## Quick Start

> You can build image locally: `docker build -t ghcr.io/jus1d/kypidbot:latest .` or pull it from `ghcr.io`: `docker pull ghcr.io/jus1d/kypidbot:latest`

Fill up the config in `./config` and set environment variables: `CONFIG_PATH` -- where config for the bot is located, `POSTGRES_PASSWORD` -- password, with which postgres will start

Run compose

```bash
$ docker compose up -d
```

Pull model for  vectorizing abouts

```bash
$ docker exec ollama ollama pull mxbai-embed-large
```
