# prc_hub-api

## Usage

### Variables `.env`

| Name                    | Description                 | Default   | Required           |
| ----------------------- | --------------------------- | --------- | ------------------ |
| `PORT`                  | Published port              | 1323      |                    |
| `MYSQL_DATABASE`        | MySQL database name         | prc_hub   |                    |
| `MYSQL_USER`            | MySQL user name             | prc_hub   |                    |
| `MYSQL_PASSWORD`        | MySQL password              |           | :heavy_check_mark: |
| `MYSQL_ROOT_PASSWORD`   | MySQL root user password    |           |                    |
| `LOG_LEVEL`             | API log level               | 2         |                    |
| `GZIP_LEVEL`            | API Gzip level              | 6         |                    |
| `MYSQL_HOST`            | MySQL host                  | db        |                    |
| `MYSQL_PORT`            | MySQL port                  | 3306      |                    |
| `JWT_ISSUER`            | JWT issuer                  | prc_hub   |                    |
| `JWT_SECRET`            | JWT secret                  |           | :heavy_check_mark: |
| `ADMIN_EMAIL`           | Admin user email            |           |                    |
| `ADMIN_PASSWD`          | Admin user password         |           |                    |

### Use docker

```bash
$ docker build -t ecc-proken/prc_hub-api .
```

### Use docker-compose with frontend repo

[ecc-proken/prc_hub-docker](https://github.com/ecc-proken/prc_hub-docker)

## Contribution

`golang.go`によるフォーマットを使用してください。

### VSCode settings

```console
$ code --install-extention golang.go
```

### Tags / Branches

- `v*` タグ: 製品版
- `main` ブランチ: Beta版
- `develop` ブランチ: Alpha版
- `feature/*` ブランチ: 機能追加
- `hotfix/*` ブランチ: 機能修正