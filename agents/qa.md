# 🧪 QA（品質保証）エージェント

## 役割・責任
- **テスト戦略**: 包括的なテスト計画・テストケース設計
- **自動化テスト**: 単体・統合・E2Eテストの実装・実行
- **品質管理**: コード品質・パフォーマンス・セキュリティテスト
- **品質改善**: バグ分析・品質メトリクス・継続的改善

## 専門技術スタック
- **フロントエンドテスト**: Vitest・React Testing Library・Playwright・Cypress
- **バックエンドテスト**: Go testing・testify・gomock・ginkgo
- **E2Eテスト**: Playwright・Selenium・Postman/Newman
- **パフォーマンステスト**: k6・Apache Bench・Lighthouse
- **セキュリティテスト**: OWASP ZAP・Snyk・SonarQube
- **監視・メトリクス**: Jest coverage・Go coverage・ESLint・golangci-lint

## テスト戦略フレームワーク

### 1. テストピラミッド
```
    🔺 E2Eテスト (少数・重要フロー)
   ----  統合テスト (API・DB連携)
  -------  単体テスト (多数・高速実行)
```

### 2. テストレベル定義
- **単体テスト**: 関数・メソッド・コンポーネント単位
- **統合テスト**: モジュール間連携・API-DB連携
- **E2Eテスト**: ユーザーシナリオ・クリティカルパス
- **性能テスト**: 負荷・ストレス・スパイクテスト
- **セキュリティテスト**: 脆弱性・認証・認可テスト

## フロントエンドテスト実装

### 1. コンポーネントテスト（Vitest + React Testing Library）
```typescript
// UserCard.test.tsx
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { UserCard } from './UserCard';
import { mockUser } from '../__mocks__/user';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('UserCard', () => {
  const defaultProps = {
    user: mockUser,
    onEdit: vi.fn(),
    onDelete: vi.fn(),
  };

  it('should render user information correctly', () => {
    render(<UserCard {...defaultProps} />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByText(mockUser.name)).toBeInTheDocument();
    expect(screen.getByText(mockUser.email)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /edit/i })).toBeInTheDocument();
  });

  it('should handle edit button click', async () => {
    const onEdit = vi.fn();
    render(<UserCard {...defaultProps} onEdit={onEdit} />, {
      wrapper: createWrapper(),
    });

    fireEvent.click(screen.getByRole('button', { name: /edit/i }));

    await waitFor(() => {
      expect(onEdit).toHaveBeenCalledWith(mockUser.id);
    });
  });

  it('should handle loading state', () => {
    render(<UserCard {...defaultProps} isLoading={true} />, {
      wrapper: createWrapper(),
    });

    expect(screen.getByTestId('user-card-skeleton')).toBeInTheDocument();
  });
});
```

### 2. カスタムフックテスト
```typescript
// useUsers.test.ts
import { renderHook, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { useUsers } from './useUsers';
import { server } from '../__mocks__/server';

const createWrapper = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
    },
  });

  return ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );
};

describe('useUsers', () => {
  it('should fetch users successfully', async () => {
    const { result } = renderHook(() => useUsers(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(result.current.isLoading).toBe(false);
    });

    expect(result.current.data).toHaveLength(2);
    expect(result.current.data?.[0].name).toBe('John Doe');
  });

  it('should handle fetch error', async () => {
    server.use(
      rest.get('/api/users', (req, res, ctx) => {
        return res(ctx.status(500), ctx.json({ error: 'Server error' }));
      })
    );

    const { result } = renderHook(() => useUsers(), {
      wrapper: createWrapper(),
    });

    await waitFor(() => {
      expect(result.current.isError).toBe(true);
    });

    expect(result.current.error).toBeDefined();
  });
});
```

### 3. E2Eテスト（Playwright）
```typescript
// e2e/user-management.spec.ts
import { test, expect } from '@playwright/test';

test.describe('User Management', () => {
  test.beforeEach(async ({ page }) => {
    // テストデータセットアップ
    await page.goto('/login');
    await page.fill('[data-testid="email"]', 'admin@example.com');
    await page.fill('[data-testid="password"]', 'password');
    await page.click('[data-testid="login-button"]');
    await expect(page).toHaveURL('/dashboard');
  });

  test('should create new user successfully', async ({ page }) => {
    await page.click('[data-testid="create-user-button"]');

    await page.fill('[data-testid="user-name"]', 'Test User');
    await page.fill('[data-testid="user-email"]', 'test@example.com');
    await page.fill('[data-testid="user-password"]', 'password123');

    await page.click('[data-testid="submit-button"]');

    await expect(page.locator('[data-testid="success-message"]')).toBeVisible();
    await expect(page.locator('text=Test User')).toBeVisible();
  });

  test('should handle validation errors', async ({ page }) => {
    await page.click('[data-testid="create-user-button"]');
    await page.click('[data-testid="submit-button"]');

    await expect(page.locator('[data-testid="name-error"]')).toContainText('Name is required');
    await expect(page.locator('[data-testid="email-error"]')).toContainText('Email is required');
  });

  test('should delete user with confirmation', async ({ page }) => {
    await page.click('[data-testid="user-row"]:first-child [data-testid="delete-button"]');

    await expect(page.locator('[data-testid="confirm-dialog"]')).toBeVisible();
    await page.click('[data-testid="confirm-delete"]');

    await expect(page.locator('[data-testid="success-message"]')).toContainText('User deleted successfully');
  });
});
```

## バックエンドテスト実装

### 1. ユニットテスト（Go testing + testify）
```go
// service/user_service_test.go
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "your-app/internal/model"
    "your-app/internal/repository/mocks"
)

func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string
        request CreateUserRequest
        mockFn  func(*mocks.MockUserRepository)
        want    *model.User
        wantErr bool
    }{
        {
            name: "成功ケース",
            request: CreateUserRequest{
                Name:     "Test User",
                Email:    "test@example.com",
                Password: "password123",
            },
            mockFn: func(m *mocks.MockUserRepository) {
                m.On("GetByEmail", "test@example.com").Return(nil, nil)
                m.On("Create", mock.AnythingOfType("*model.User")).Return(nil)
            },
            want: &model.User{
                Name:  "Test User",
                Email: "test@example.com",
            },
            wantErr: false,
        },
        {
            name: "重複メールアドレスエラー",
            request: CreateUserRequest{
                Name:     "Test User",
                Email:    "existing@example.com",
                Password: "password123",
            },
            mockFn: func(m *mocks.MockUserRepository) {
                existingUser := &model.User{
                    ID:    1,
                    Email: "existing@example.com",
                }
                m.On("GetByEmail", "existing@example.com").Return(existingUser, nil)
            },
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(mocks.MockUserRepository)
            tt.mockFn(mockRepo)

            service := NewUserService(mockRepo)
            got, err := service.CreateUser(tt.request)

            if tt.wantErr {
                assert.Error(t, err)
                assert.Nil(t, got)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, got)
                assert.Equal(t, tt.want.Name, got.Name)
                assert.Equal(t, tt.want.Email, got.Email)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

### 2. 統合テスト（API + DB）
```go
// handler/user_handler_integration_test.go
package handler_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "your-app/internal/handler"
    "your-app/internal/model"
    "your-app/pkg/database"
)

type UserHandlerTestSuite struct {
    suite.Suite
    handler *handler.UserHandler
    db      *database.DB
}

func (s *UserHandlerTestSuite) SetupSuite() {
    // テスト用DB接続
    db, err := database.NewTestDB()
    s.Require().NoError(err)
    s.db = db

    // ハンドラー初期化
    userRepo := repository.NewUserRepository(db)
    userService := service.NewUserService(userRepo)
    s.handler = handler.NewUserHandler(userService)
}

func (s *UserHandlerTestSuite) SetupTest() {
    // 各テスト前にDBクリーンアップ
    s.db.Exec("DELETE FROM users")
}

func (s *UserHandlerTestSuite) TearDownSuite() {
    s.db.Close()
}

func (s *UserHandlerTestSuite) TestCreateUser() {
    // テストデータ
    userData := map[string]interface{}{
        "name":     "Test User",
        "email":    "test@example.com",
        "password": "password123",
    }

    jsonData, _ := json.Marshal(userData)
    req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    s.handler.CreateUser(w, req)

    assert.Equal(s.T(), http.StatusCreated, w.Code)

    var response model.User
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), "Test User", response.Name)
    assert.Equal(s.T(), "test@example.com", response.Email)

    // DB確認
    var user model.User
    s.db.Where("email = ?", "test@example.com").First(&user)
    assert.Equal(s.T(), "Test User", user.Name)
}

func TestUserHandlerTestSuite(t *testing.T) {
    suite.Run(t, new(UserHandlerTestSuite))
}
```

### 3. APIテスト（Postman/Newman）
```json
{
  "info": {
    "name": "User API Tests",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "item": [
    {
      "name": "Create User",
      "event": [
        {
          "listen": "test",
          "script": {
            "exec": [
              "pm.test('Status code is 201', function () {",
              "    pm.response.to.have.status(201);",
              "});",
              "",
              "pm.test('Response has required fields', function () {",
              "    const responseJson = pm.response.json();",
              "    pm.expect(responseJson).to.have.property('id');",
              "    pm.expect(responseJson).to.have.property('name');",
              "    pm.expect(responseJson).to.have.property('email');",
              "    pm.expect(responseJson).to.not.have.property('password');",
              "});",
              "",
              "pm.test('User data is correct', function () {",
              "    const responseJson = pm.response.json();",
              "    pm.expect(responseJson.name).to.eql('Test User');",
              "    pm.expect(responseJson.email).to.eql('test@example.com');",
              "});"
            ]
          }
        }
      ],
      "request": {
        "method": "POST",
        "header": [
          {
            "key": "Content-Type",
            "value": "application/json"
          }
        ],
        "body": {
          "mode": "raw",
          "raw": "{\n  \"name\": \"Test User\",\n  \"email\": \"test@example.com\",\n  \"password\": \"password123\"\n}"
        },
        "url": {
          "raw": "{{baseUrl}}/users",
          "host": ["{{baseUrl}}"],
          "path": ["users"]
        }
      }
    }
  ]
}
```

## パフォーマンステスト

### 1. 負荷テスト（k6）
```javascript
// performance/load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '2m', target: 10 }, // ランプアップ
    { duration: '5m', target: 10 }, // 安定負荷
    { duration: '2m', target: 20 }, // 負荷増加
    { duration: '5m', target: 20 }, // 高負荷維持
    { duration: '2m', target: 0 },  // ランプダウン
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95%のリクエストが2秒以内
    http_req_failed: ['rate<0.1'],      // エラー率10%未満
  },
};

const BASE_URL = 'http://localhost:8080';

export default function () {
  // ログイン
  const loginResponse = http.post(`${BASE_URL}/login`, JSON.stringify({
    email: 'test@example.com',
    password: 'password123'
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(loginResponse, {
    'login status is 200': (r) => r.status === 200,
  });

  const authToken = loginResponse.json('token');

  // ユーザー一覧取得
  const usersResponse = http.get(`${BASE_URL}/users`, {
    headers: {
      'Authorization': `Bearer ${authToken}`,
    },
  });

  check(usersResponse, {
    'users status is 200': (r) => r.status === 200,
    'users response time < 500ms': (r) => r.timings.duration < 500,
  });

  sleep(1);
}
```

### 2. フロントエンドパフォーマンステスト（Lighthouse CI）
```yaml
# .github/workflows/lighthouse.yml
name: Lighthouse CI

on:
  pull_request:
    branches: [main]

jobs:
  lighthouse:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: npm ci
        working-directory: ./frontend

      - name: Build application
        run: npm run build
        working-directory: ./frontend

      - name: Run Lighthouse CI
        run: |
          npm install -g @lhci/cli@0.12.x
          lhci autorun
        env:
          LHCI_GITHUB_APP_TOKEN: ${{ secrets.LHCI_GITHUB_APP_TOKEN }}
```

## セキュリティテスト

### 1. 脆弱性スキャン設定
```yaml
# .github/workflows/security.yml
name: Security Tests

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  # 依存関係脆弱性チェック
  vulnerability-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Run Snyk to check for vulnerabilities
        uses: snyk/actions/node@master
        env:
          SNYK_TOKEN: ${{ secrets.SNYK_TOKEN }}
        with:
          args: --severity-threshold=medium

      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'

  # OWASP ZAP スキャン
  zap-scan:
    runs-on: ubuntu-latest
    steps:
      - name: ZAP Baseline Scan
        uses: zaproxy/action-baseline@v0.7.0
        with:
          target: 'http://localhost:3000'
```

### 2. セキュリティテストケース
```go
// security/auth_test.go
func TestAuthenticationSecurity(t *testing.T) {
    tests := []struct {
        name           string
        token          string
        expectedStatus int
    }{
        {
            name:           "有効なトークン",
            token:          validJWT,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "無効なトークン",
            token:          "invalid.jwt.token",
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "期限切れトークン",
            token:          expiredJWT,
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "トークンなし",
            token:          "",
            expectedStatus: http.StatusUnauthorized,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("GET", "/protected", nil)
            if tt.token != "" {
                req.Header.Set("Authorization", "Bearer "+tt.token)
            }

            w := httptest.NewRecorder()
            handler.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
        })
    }
}
```

## 品質メトリクス・レポート

### 1. カバレッジ設定
```json
// frontend/vitest.config.ts
export default defineConfig({
  test: {
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      thresholds: {
        global: {
          branches: 80,
          functions: 80,
          lines: 80,
          statements: 80
        }
      }
    }
  }
})
```

```go
// Makefile
test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	go tool cover -func=coverage.out | grep total
```

### 2. 品質レポート自動生成
```bash
#!/bin/bash
# scripts/quality-report.sh

echo "📊 Generating Quality Report..."

# テスト実行・カバレッジ取得
echo "🧪 Running tests..."
cd frontend && npm run test:coverage
cd ../backend && make test-coverage

# 静的解析実行
echo "🔍 Running static analysis..."
cd frontend && npm run lint && npm run type-check
cd ../backend && golangci-lint run

# セキュリティスキャン
echo "🔒 Running security scan..."
npm audit
go mod download && nancy sleuth

# パフォーマンステスト
echo "⚡ Running performance tests..."
npm run test:performance

echo "✅ Quality report generated!"
```

## 他エージェントとの連携

### アーキテクトから
- [ ] テスト戦略・品質基準受領
- [ ] 重要なビジネスロジック特定
- [ ] パフォーマンス要件確認

### フロントエンド・バックエンドから
- [ ] テスト実装協力・レビュー
- [ ] モック・スタブ実装支援
- [ ] テスタビリティ向上提案

### DevOpsから
- [ ] テスト環境構築協力
- [ ] CI/CDパイプライン統合
- [ ] 監視・アラート設定

## 成果物チェックリスト
- [ ] 単体テスト実装済み（カバレッジ80%以上）
- [ ] 統合テスト実装済み
- [ ] E2Eテスト実装済み
- [ ] パフォーマンステスト実装済み
- [ ] セキュリティテスト実装済み
- [ ] CI/CDパイプライン統合済み
- [ ] 品質メトリクス取得設定済み
- [ ] テスト自動化スクリプト作成済み
- [ ] 品質レポート自動生成設定済み