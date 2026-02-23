# DevOps技術コンテキスト

## 技術スタック詳細

### 核心技術
- **Docker 24.0+**: Multi-stage builds・BuildKit・Docker Compose v2
- **Docker Compose v2**: サービス管理・ネットワーク・ボリューム管理
- **Nginx**: リバースプロキシ・ロードバランシング・SSL終端
- **MySQL 8.0**: レプリケーション・パフォーマンススキーマ・JSON型

### インフラアーキテクチャ
```yaml
# 開発環境構成例
services:
  # リバースプロキシ・ロードバランサー
  nginx:
    image: nginx:1.25-alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/templates:/etc/nginx/templates:ro
      - ./nginx/ssl:/etc/ssl/certs:ro
      - nginx_logs:/var/log/nginx
    environment:
      - NGINX_ENVSUBST_OUTPUT_DIR=/etc/nginx/conf.d
      - BACKEND_HOST=backend
      - FRONTEND_HOST=frontend
    depends_on:
      - frontend
      - backend
    networks:
      - app-network

  # フロントエンド（React + Vite）
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile.dev
      target: development
    ports:
      - "3000:3000"
    volumes:
      - ./frontend:/app:cached
      - /app/node_modules
      - frontend_dist:/app/dist
    environment:
      - NODE_ENV=development
      - VITE_API_URL=http://localhost:8080/api
      - VITE_WS_URL=ws://localhost:8080/ws
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000"]
      interval: 30s
      timeout: 10s
      retries: 3

  # バックエンド（Go + chi）
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile.dev
      target: development
    ports:
      - "8080:8080"
    volumes:
      - ./backend:/app:cached
      - go_cache:/go/pkg/mod
    environment:
      - GO_ENV=development
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USER=developer
      - DB_PASSWORD=dev_password
      - DB_NAME=app_development
      - REDIS_URL=redis:6379
      - JWT_SECRET=dev-jwt-secret
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_started
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  # データベース（MySQL）
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./database/init:/docker-entrypoint-initdb.d:ro
      - ./database/config/my.cnf:/etc/mysql/conf.d/custom.cnf:ro
    environment:
      - MYSQL_ROOT_PASSWORD=root_password
      - MYSQL_DATABASE=app_development
      - MYSQL_USER=developer
      - MYSQL_PASSWORD=dev_password
      - MYSQL_CHARSET=utf8mb4
      - MYSQL_COLLATION=utf8mb4_unicode_ci
    command: >
      --character-set-server=utf8mb4
      --collation-server=utf8mb4_unicode_ci
      --default-authentication-plugin=mysql_native_password
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-proot_password"]
      interval: 10s
      timeout: 5s
      retries: 5

  # キャッシュ（Redis）
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./redis/redis.conf:/usr/local/etc/redis/redis.conf:ro
    command: redis-server /usr/local/etc/redis/redis.conf
    networks:
      - app-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 3

volumes:
  mysql_data:
    driver: local
  redis_data:
    driver: local
  go_cache:
    driver: local
  frontend_dist:
    driver: local
  nginx_logs:
    driver: local

networks:
  app-network:
    driver: bridge
    ipam:
      driver: default
      config:
        - subnet: 172.20.0.0/16
```

## Dockerfile最適化

### 1. フロントエンドDockerfile
```dockerfile
# frontend/Dockerfile
ARG NODE_VERSION=18

# 開発環境ステージ
FROM node:${NODE_VERSION}-alpine AS development

WORKDIR /app

# パッケージインストール（キャッシュ効率化）
COPY package*.json ./
RUN npm ci --only=production && npm cache clean --force

# 開発依存関係インストール
COPY package*.json ./
RUN npm ci

# ソースコード追加
COPY . .

# ヘルスチェック用curl追加
RUN apk add --no-cache curl

EXPOSE 3000

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

# ビルドステージ
FROM development AS builder

RUN npm run build

# 本番環境ステージ
FROM nginx:1.25-alpine AS production

# カスタムnginx設定
COPY nginx.conf /etc/nginx/nginx.conf

# ビルド成果物コピー
COPY --from=builder /app/dist /usr/share/nginx/html

# セキュリティ向上
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001 && \
    chown -R nextjs:nodejs /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

### 2. バックエンドDockerfile
```dockerfile
# backend/Dockerfile
ARG GO_VERSION=1.21

# 開発環境ステージ
FROM golang:${GO_VERSION}-alpine AS development

WORKDIR /app

# 開発ツールインストール
RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest

# 依存関係インストール
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# ソースコード追加
COPY . .

# ヘルスチェック用curl追加
RUN apk add --no-cache curl

EXPOSE 8080 2345

CMD ["air", "-c", ".air.toml"]

# ビルドステージ
FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

# ビルド用依存関係インストール
RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# 静的リンクビルド
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main cmd/server/main.go

# 本番環境ステージ
FROM scratch AS production

# タイムゾーンデータ・CA証明書追加
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# アプリケーション追加
COPY --from=builder /app/main /main

EXPOSE 8080

ENTRYPOINT ["/main"]
```

## CI/CDパイプライン

### 1. GitHub Actions設定
```yaml
# .github/workflows/ci-cd.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

env:
  DOCKER_BUILDKIT: 1
  COMPOSE_DOCKER_CLI_BUILD: 1

jobs:
  # 変更検出
  changes:
    runs-on: ubuntu-latest
    outputs:
      frontend: ${{ steps.changes.outputs.frontend }}
      backend: ${{ steps.changes.outputs.backend }}
      docker: ${{ steps.changes.outputs.docker }}
    steps:
      - uses: actions/checkout@v4
      - uses: dorny/paths-filter@v2
        id: changes
        with:
          filters: |
            frontend:
              - 'frontend/**'
              - 'package*.json'
            backend:
              - 'backend/**'
              - 'go.mod'
              - 'go.sum'
            docker:
              - 'docker-compose*.yml'
              - '**/Dockerfile*'
              - 'nginx/**'

  # フロントエンドテスト
  frontend-test:
    needs: changes
    if: needs.changes.outputs.frontend == 'true'
    runs-on: ubuntu-latest

    strategy:
      matrix:
        node-version: [18, 20]

    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js ${{ matrix.node-version }}
        uses: actions/setup-node@v3
        with:
          node-version: ${{ matrix.node-version }}
          cache: 'npm'
          cache-dependency-path: frontend/package-lock.json

      - name: Install dependencies
        run: npm ci
        working-directory: ./frontend

      - name: Lint
        run: npm run lint
        working-directory: ./frontend

      - name: Type check
        run: npm run type-check
        working-directory: ./frontend

      - name: Unit tests
        run: npm run test:coverage
        working-directory: ./frontend

      - name: Build
        run: npm run build
        working-directory: ./frontend

      - name: Upload coverage reports
        uses: codecov/codecov-action@v3
        with:
          file: ./frontend/coverage/lcov.info

  # バックエンドテスト
  backend-test:
    needs: changes
    if: needs.changes.outputs.backend == 'true'
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root_password
          MYSQL_DATABASE: test_db
          MYSQL_USER: test_user
          MYSQL_PASSWORD: test_password
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping -h localhost"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=5

      redis:
        image: redis:7-alpine
        ports:
          - 6379:6379
        options: >-
          --health-cmd="redis-cli ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=5

    strategy:
      matrix:
        go-version: [1.21, 1.22]

    steps:
      - uses: actions/checkout@v4

      - name: Setup Go ${{ matrix.go-version }}
        uses: actions/setup-go@v4
        with:
          go-version: ${{ matrix.go-version }}
          cache-dependency-path: backend/go.sum

      - name: Install dependencies
        run: go mod download && go mod verify
        working-directory: ./backend

      - name: Lint
        uses: golangci/golangci-lint-action@v3
        with:
          version: latest
          working-directory: backend

      - name: Unit tests
        run: go test -v -race -coverprofile=coverage.out ./...
        working-directory: ./backend
        env:
          DB_HOST: localhost
          DB_PORT: 3306
          DB_USER: test_user
          DB_PASSWORD: test_password
          DB_NAME: test_db
          REDIS_URL: localhost:6379

      - name: Integration tests
        run: go test -v -tags=integration ./...
        working-directory: ./backend
        env:
          DB_HOST: localhost
          DB_PORT: 3306
          DB_USER: test_user
          DB_PASSWORD: test_password
          DB_NAME: test_db
          REDIS_URL: localhost:6379

      - name: Upload coverage reports
        uses: codecov/codecov-action@v3
        with:
          file: ./backend/coverage.out

  # Docker イメージビルド
  docker-build:
    needs: [frontend-test, backend-test]
    if: always() && (needs.frontend-test.result == 'success' || needs.frontend-test.result == 'skipped') && (needs.backend-test.result == 'success' || needs.backend-test.result == 'skipped')
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Login to Docker Hub
        if: github.event_name != 'pull_request'
        uses: docker/login-action@v3
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: |
            myapp/frontend
            myapp/backend
          tags: |
            type=ref,event=branch
            type=ref,event=pr
            type=sha,prefix={{branch}}-
            type=raw,value=latest,enable={{is_default_branch}}

      - name: Build and push Frontend
        uses: docker/build-push-action@v5
        with:
          context: ./frontend
          target: production
          push: ${{ github.event_name != 'pull_request' }}
          tags: myapp/frontend:${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

      - name: Build and push Backend
        uses: docker/build-push-action@v5
        with:
          context: ./backend
          target: production
          push: ${{ github.event_name != 'pull_request' }}
          tags: myapp/backend:${{ steps.meta.outputs.tags }}
          cache-from: type=gha
          cache-to: type=gha,mode=max

  # E2Eテスト
  e2e-test:
    needs: docker-build
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Start services
        run: |
          docker-compose -f docker-compose.test.yml up -d
          sleep 30

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install Playwright
        run: |
          npm install @playwright/test
          npx playwright install --with-deps

      - name: Run E2E tests
        run: npx playwright test
        env:
          BASE_URL: http://localhost:3000

      - name: Upload test results
        uses: actions/upload-artifact@v3
        if: always()
        with:
          name: playwright-report
          path: playwright-report/

      - name: Cleanup
        run: docker-compose -f docker-compose.test.yml down -v

  # デプロイ
  deploy:
    needs: [docker-build, e2e-test]
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v4

      - name: Deploy to staging
        run: |
          echo "Deploying to staging environment"
          # デプロイスクリプト実行
          ./scripts/deploy-staging.sh

      - name: Run smoke tests
        run: |
          echo "Running smoke tests"
          ./scripts/smoke-test.sh

      - name: Deploy to production
        if: success()
        run: |
          echo "Deploying to production environment"
          ./scripts/deploy-production.sh
        env:
          DEPLOY_KEY: ${{ secrets.PRODUCTION_DEPLOY_KEY }}
```

## 監視・ログ管理

### 1. Prometheus + Grafana設定
```yaml
# monitoring/docker-compose.yml
version: '3.8'

services:
  # メトリクス収集
  prometheus:
    image: prom/prometheus:v2.45.0
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - ./prometheus/rules:/etc/prometheus/rules:ro
      - prometheus_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--web.enable-lifecycle'
      - '--web.enable-admin-api'
    networks:
      - monitoring

  # メトリクス可視化
  grafana:
    image: grafana/grafana:10.0.0
    ports:
      - "3001:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning:ro
      - ./grafana/dashboards:/var/lib/grafana/dashboards:ro
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
      - GF_SERVER_DOMAIN=localhost
      - GF_SMTP_ENABLED=true
      - GF_SMTP_HOST=smtp.gmail.com:587
      - GF_SMTP_USER=${SMTP_USER}
      - GF_SMTP_PASSWORD=${SMTP_PASSWORD}
    networks:
      - monitoring

  # システムメトリクス
  node-exporter:
    image: prom/node-exporter:v1.6.0
    ports:
      - "9100:9100"
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - '--path.procfs=/host/proc'
      - '--path.rootfs=/rootfs'
      - '--path.sysfs=/host/sys'
      - '--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)'
    networks:
      - monitoring

  # アラートマネージャー
  alertmanager:
    image: prom/alertmanager:v0.25.0
    ports:
      - "9093:9093"
    volumes:
      - ./alertmanager/alertmanager.yml:/etc/alertmanager/alertmanager.yml:ro
      - alertmanager_data:/alertmanager
    command:
      - '--config.file=/etc/alertmanager/alertmanager.yml'
      - '--storage.path=/alertmanager'
      - '--web.external-url=http://localhost:9093'
    networks:
      - monitoring

volumes:
  prometheus_data:
  grafana_data:
  alertmanager_data:

networks:
  monitoring:
    driver: bridge
```

### 2. ログ管理（ELK Stack）
```yaml
# logging/docker-compose.yml
version: '3.8'

services:
  # ログ検索・保存
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.8.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
      - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
    ports:
      - "9200:9200"
    volumes:
      - elasticsearch_data:/usr/share/elasticsearch/data
    networks:
      - logging

  # ログ可視化
  kibana:
    image: docker.elastic.co/kibana/kibana:8.8.0
    environment:
      - ELASTICSEARCH_HOSTS=http://elasticsearch:9200
      - KIBANA_SYSTEM_PASSWORD=changeme
    ports:
      - "5601:5601"
    depends_on:
      - elasticsearch
    networks:
      - logging

  # ログ収集・転送
  logstash:
    image: docker.elastic.co/logstash/logstash:8.8.0
    volumes:
      - ./logstash/pipeline:/usr/share/logstash/pipeline:ro
      - ./logstash/config/logstash.yml:/usr/share/logstash/config/logstash.yml:ro
    ports:
      - "5044:5044"
      - "9600:9600"
    environment:
      - "LS_JAVA_OPTS=-Xms512m -Xmx512m"
    depends_on:
      - elasticsearch
    networks:
      - logging

  # ログ収集エージェント
  filebeat:
    image: docker.elastic.co/beats/filebeat:8.8.0
    user: root
    volumes:
      - ./filebeat/filebeat.yml:/usr/share/filebeat/filebeat.yml:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      - ELASTICSEARCH_HOST=elasticsearch:9200
      - LOGSTASH_HOST=logstash:5044
    depends_on:
      - logstash
    networks:
      - logging

volumes:
  elasticsearch_data:

networks:
  logging:
    driver: bridge
```

## セキュリティ・SSL設定

### 1. Nginx SSL設定
```nginx
# nginx/templates/nginx.conf.template
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
    use epoll;
    multi_accept on;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # ログフォーマット
    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for" '
                    'rt=$request_time uct="$upstream_connect_time" '
                    'uht="$upstream_header_time" urt="$upstream_response_time"';

    access_log /var/log/nginx/access.log main;

    # パフォーマンス最適化
    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 10M;

    # Gzip圧縮
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/javascript
        application/xml+rss
        application/json;

    # セキュリティヘッダー
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self';" always;

    # レート制限
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;
    limit_req_zone $binary_remote_addr zone=login:10m rate=1r/s;

    # アップストリーム定義
    upstream frontend {
        server ${FRONTEND_HOST}:3000;
        keepalive 32;
    }

    upstream backend {
        server ${BACKEND_HOST}:8080;
        keepalive 32;
    }

    # HTTPからHTTPSへリダイレクト
    server {
        listen 80;
        server_name _;
        return 301 https://$host$request_uri;
    }

    # HTTPS設定
    server {
        listen 443 ssl http2;
        server_name localhost;

        # SSL証明書
        ssl_certificate /etc/ssl/certs/server.crt;
        ssl_certificate_key /etc/ssl/certs/server.key;

        # SSL設定
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384:DHE-RSA-AES256-GCM-SHA384;
        ssl_prefer_server_ciphers off;
        ssl_session_cache shared:SSL:10m;
        ssl_session_timeout 10m;

        # HSTS
        add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

        # フロントエンド
        location / {
            proxy_pass http://frontend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
            proxy_cache_bypass $http_upgrade;

            # キャッシュ設定
            location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
                expires 1y;
                add_header Cache-Control "public, immutable";
            }
        }

        # API
        location /api {
            limit_req zone=api burst=20 nodelay;

            proxy_pass http://backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;

            # タイムアウト設定
            proxy_connect_timeout 5s;
            proxy_send_timeout 60s;
            proxy_read_timeout 60s;
        }

        # ログイン API（レート制限強化）
        location /api/auth/login {
            limit_req zone=login burst=5 nodelay;
            proxy_pass http://backend;
            # 他の設定は/apiと同じ
        }

        # ヘルスチェック
        location /health {
            access_log off;
            return 200 "healthy\n";
            add_header Content-Type text/plain;
        }
    }
}
```

## バックアップ・災害復旧

### 1. データベースバックアップ
```bash
#!/bin/bash
# scripts/backup-database.sh

set -e

# 設定
BACKUP_DIR="/backup/mysql"
DATE=$(date +%Y%m%d_%H%M%S)
RETENTION_DAYS=7

# バックアップディレクトリ作成
mkdir -p $BACKUP_DIR

# MySQLバックアップ
echo "Creating MySQL backup..."
docker exec mysql mysqldump \
    -u root -p${MYSQL_ROOT_PASSWORD} \
    --routines --triggers --single-transaction \
    --master-data=2 \
    ${DB_NAME} > "${BACKUP_DIR}/backup_${DATE}.sql"

# 圧縮
gzip "${BACKUP_DIR}/backup_${DATE}.sql"

# 古いバックアップ削除
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +$RETENTION_DAYS -delete

echo "Backup completed: backup_${DATE}.sql.gz"

# S3アップロード（オプション）
if [[ -n "$AWS_S3_BUCKET" ]]; then
    aws s3 cp "${BACKUP_DIR}/backup_${DATE}.sql.gz" \
        "s3://${AWS_S3_BUCKET}/backups/mysql/"
    echo "Backup uploaded to S3"
fi
```

### 2. データリストア
```bash
#!/bin/bash
# scripts/restore-database.sh

set -e

if [[ $# -eq 0 ]]; then
    echo "Usage: $0 <backup_file>"
    exit 1
fi

BACKUP_FILE=$1

if [[ ! -f $BACKUP_FILE ]]; then
    echo "Backup file not found: $BACKUP_FILE"
    exit 1
fi

echo "Restoring database from: $BACKUP_FILE"

# 現在のデータベースバックアップ
echo "Creating current database backup..."
./backup-database.sh

# リストア実行
if [[ $BACKUP_FILE == *.gz ]]; then
    gunzip -c $BACKUP_FILE | docker exec -i mysql mysql \
        -u root -p${MYSQL_ROOT_PASSWORD} ${DB_NAME}
else
    docker exec -i mysql mysql \
        -u root -p${MYSQL_ROOT_PASSWORD} ${DB_NAME} < $BACKUP_FILE
fi

echo "Database restore completed"
```

## パフォーマンス最適化

### 1. Docker最適化
```dockerfile
# .dockerignore
node_modules
npm-debug.log*
.git
.gitignore
README.md
.env
.nyc_output
coverage
.coverage
.pytest_cache
**/__pycache__
.vscode
.idea
```

### 2. リソース制限・ヘルスチェック
```yaml
# docker-compose.prod.yml抜粋
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
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
        window: 120s
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  mysql:
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 2G
        reservations:
          cpus: '0.5'
          memory: 1G
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost", "-u", "root", "-p$$MYSQL_ROOT_PASSWORD"]
      interval: 10s
      timeout: 5s
      retries: 5
      start_period: 30s
```

## トラブルシューティング

### 1. ログ監視スクリプト
```bash
#!/bin/bash
# scripts/monitor-logs.sh

# 全サービスログ監視
echo "=== Service Status ==="
docker-compose ps

echo -e "\n=== Recent Errors ==="
docker-compose logs --tail=50 | grep -i error

echo -e "\n=== Resource Usage ==="
docker stats --no-stream

echo -e "\n=== Disk Usage ==="
df -h

echo -e "\n=== Memory Usage ==="
free -h
```

### 2. ヘルスチェックスクリプト
```bash
#!/bin/bash
# scripts/health-check.sh

set -e

BASE_URL=${1:-"http://localhost"}

echo "Checking application health..."

# フロントエンド
if curl -f "${BASE_URL}:3000" > /dev/null 2>&1; then
    echo "✅ Frontend: OK"
else
    echo "❌ Frontend: FAIL"
    exit 1
fi

# バックエンドAPI
if curl -f "${BASE_URL}:8080/health" > /dev/null 2>&1; then
    echo "✅ Backend: OK"
else
    echo "❌ Backend: FAIL"
    exit 1
fi

# データベース
if docker exec mysql mysqladmin ping -h localhost -u root -p${MYSQL_ROOT_PASSWORD} > /dev/null 2>&1; then
    echo "✅ Database: OK"
else
    echo "❌ Database: FAIL"
    exit 1
fi

echo "🎉 All services are healthy!"
```