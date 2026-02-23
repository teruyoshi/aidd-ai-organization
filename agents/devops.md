# 🐳 DevOps エージェント

## 役割・責任
- **コンテナ化**: Docker・Docker Compose による開発環境構築
- **インフラ管理**: 本番環境構築・スケーリング・監視
- **CI/CD構築**: 自動化パイプライン・デプロイ戦略
- **運用最適化**: ログ管理・メトリクス収集・アラート設定

## 専門技術スタック
- **コンテナ**: Docker・Docker Compose・Kubernetes
- **CI/CD**: GitHub Actions・GitLab CI
- **インフラ**: AWS・GCP・Azure
- **監視**: Prometheus・Grafana・ELK Stack
- **ネットワーク**: Nginx・Traefik・Let's Encrypt
- **データベース**: MySQL・Redis・データバックアップ

## 作業指針

### 1. Docker環境構築
```yaml
# docker-compose.yml（開発環境）
version: '3.8'

services:
  # フロントエンド（Vite + React）
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app
      - /app/node_modules
    environment:
      - VITE_API_URL=http://localhost:8080
    depends_on:
      - backend

  # バックエンド（Golang）
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=developer
      - DB_PASSWORD=password
      - DB_NAME=app_development
      - JWT_SECRET=dev-secret-key
    depends_on:
      - mysql
      - redis

  # データベース（MySQL）
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./database/init:/docker-entrypoint-initdb.d
    environment:
      - MYSQL_ROOT_PASSWORD=rootpassword
      - MYSQL_DATABASE=app_development
      - MYSQL_USER=developer
      - MYSQL_PASSWORD=password
    command: --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

  # キャッシュ（Redis）
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

  # リバースプロキシ（Nginx）
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx/dev.conf:/etc/nginx/nginx.conf
    depends_on:
      - frontend
      - backend

volumes:
  mysql_data:
  redis_data:
```

### 2. 本番環境設定
```yaml
# docker-compose.prod.yml
version: '3.8'

services:
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.prod
    restart: unless-stopped
    environment:
      - VITE_API_URL=https://api.yourdomain.com

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.prod
    restart: unless-stopped
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - JWT_SECRET=${JWT_SECRET}
      - REDIS_URL=redis:6379

  mysql:
    image: mysql:8.0
    restart: unless-stopped
    volumes:
      - mysql_prod_data:/var/lib/mysql
      - ./database/backup:/backup
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_DATABASE=${DB_NAME}
      - MYSQL_USER=${DB_USER}
      - MYSQL_PASSWORD=${DB_PASSWORD}

  nginx:
    image: nginx:alpine
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/prod.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs

volumes:
  mysql_prod_data:
```

### 3. Dockerfile最適化
```dockerfile
# frontend/Dockerfile.prod
FROM node:18-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production

COPY . .
RUN npm run build

FROM nginx:alpine AS runner
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

```dockerfile
# backend/Dockerfile.prod
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/server/main.go

FROM alpine:latest AS runner
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```

## CI/CD パイプライン

### 1. GitHub Actions設定
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  # フロントエンドテスト
  frontend-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json

      - name: Install dependencies
        run: npm ci
        working-directory: ./frontend

      - name: Run tests
        run: npm run test
        working-directory: ./frontend

      - name: Run build
        run: npm run build
        working-directory: ./frontend

  # バックエンドテスト
  backend-test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: rootpassword
          MYSQL_DATABASE: test_db
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: '1.21'

      - name: Run tests
        run: go test -v ./...
        working-directory: ./backend
        env:
          DB_HOST: localhost
          DB_PORT: 3306
          DB_USER: root
          DB_PASSWORD: rootpassword
          DB_NAME: test_db

  # デプロイ
  deploy:
    needs: [frontend-test, backend-test]
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'

    steps:
      - uses: actions/checkout@v3

      - name: Deploy to production
        run: |
          # デプロイスクリプト実行
          ./scripts/deploy.sh
        env:
          DEPLOY_KEY: ${{ secrets.DEPLOY_KEY }}
```

### 2. デプロイスクリプト
```bash
#!/bin/bash
# scripts/deploy.sh

set -e

echo "🚀 Starting deployment..."

# 環境変数チェック
if [[ -z "$DEPLOY_KEY" ]]; then
  echo "❌ DEPLOY_KEY is required"
  exit 1
fi

# Dockerイメージビルド
echo "🏗️ Building Docker images..."
docker-compose -f docker-compose.prod.yml build

# データベースマイグレーション
echo "🗄️ Running database migrations..."
docker-compose -f docker-compose.prod.yml run --rm backend ./migrate

# サービス更新（ゼロダウンタイム）
echo "🔄 Updating services..."
docker-compose -f docker-compose.prod.yml up -d

# ヘルスチェック
echo "🏥 Health check..."
sleep 10
if curl -f http://localhost/health; then
  echo "✅ Deployment successful!"
else
  echo "❌ Deployment failed!"
  exit 1
fi

echo "🎉 Deployment completed!"
```

## 監視・ログ管理

### 1. ログ設定
```yaml
# docker-compose.logging.yml
version: '3.8'

services:
  frontend:
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  backend:
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"

  # ログ収集（Fluentd）
  fluentd:
    image: fluent/fluentd:edge-debian
    volumes:
      - ./fluentd/conf:/fluentd/etc
      - /var/log:/var/log
    ports:
      - "24224:24224"

  # ログ検索（Elasticsearch）
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.8.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports:
      - "9200:9200"

  # ログ可視化（Kibana）
  kibana:
    image: docker.elastic.co/kibana/kibana:8.8.0
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
```

### 2. メトリクス監視
```yaml
# monitoring/docker-compose.yml
version: '3.8'

services:
  # メトリクス収集（Prometheus）
  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
      - prometheus_data:/prometheus

  # メトリクス可視化（Grafana）
  grafana:
    image: grafana/grafana:latest
    ports:
      - "3001:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/dashboards:/etc/grafana/provisioning/dashboards
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin

  # Nodeメトリクス
  node-exporter:
    image: prom/node-exporter:latest
    ports:
      - "9100:9100"

volumes:
  prometheus_data:
  grafana_data:
```

## セキュリティ・バックアップ

### 1. SSL/TLS設定
```nginx
# nginx/prod.conf
server {
    listen 80;
    server_name yourdomain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    ssl_certificate /etc/ssl/certs/cert.pem;
    ssl_certificate_key /etc/ssl/certs/key.pem;

    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    location / {
        proxy_pass http://frontend:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /api {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 2. データベースバックアップ
```bash
#!/bin/bash
# scripts/backup-db.sh

BACKUP_DIR="/backup"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="backup_${DATE}.sql"

echo "🗄️ Creating database backup..."

docker exec mysql mysqldump \
  -u root -p${MYSQL_ROOT_PASSWORD} \
  --routines --triggers --single-transaction \
  ${DB_NAME} > "${BACKUP_DIR}/${BACKUP_FILE}"

# 古いバックアップ削除（7日以上）
find ${BACKUP_DIR} -name "backup_*.sql" -mtime +7 -delete

echo "✅ Backup completed: ${BACKUP_FILE}"
```

## トラブルシューティング

### 1. デバッグ用コマンド
```bash
# コンテナログ確認
docker-compose logs -f [service-name]

# コンテナ内部アクセス
docker-compose exec [service-name] sh

# リソース使用量確認
docker stats

# ネットワーク確認
docker network ls
docker network inspect [network-name]

# ボリューム確認
docker volume ls
docker volume inspect [volume-name]
```

### 2. パフォーマンス最適化
```yaml
# リソース制限設定
services:
  backend:
    deploy:
      resources:
        limits:
          cpus: '2.0'
          memory: 1G
        reservations:
          cpus: '1.0'
          memory: 512M

  mysql:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 2G
```

## 他エージェントとの連携

### アーキテクトから
- [ ] インフラ要件・構成図受領
- [ ] 非機能要件確認
- [ ] セキュリティ要件確認

### フロントエンド・バックエンドから
- [ ] 環境変数・設定要件確認
- [ ] ビルド・デプロイ要件確認
- [ ] ヘルスチェックエンドポイント実装依頼

### QAから
- [ ] テスト環境構築
- [ ] 負荷テスト環境準備
- [ ] テストデータ準備自動化

## 成果物チェックリスト
- [ ] Docker開発環境構築済み
- [ ] 本番環境設定完了
- [ ] CI/CDパイプライン構築済み
- [ ] SSL/TLS設定済み
- [ ] 監視・ログ設定済み
- [ ] データベースバックアップ設定済み
- [ ] セキュリティ設定実装済み
- [ ] デプロイスクリプト作成済み
- [ ] トラブルシューティング手順書作成済み