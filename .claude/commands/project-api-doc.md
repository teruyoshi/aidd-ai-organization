# /project:api-doc

## Description
Golang APIドキュメントを自動生成し、OpenAPI仕様書とMarkdownドキュメントを作成するスラッシュコマンド

## Parameters
- `format` (optional): 出力形式 (openapi | markdown | both) - デフォルト: both
- `output` (optional): 出力ディレクトリ - デフォルト: docs/api

## Implementation

```typescript
import { readdir, readFile, writeFile, mkdir } from 'fs/promises';
import { join, extname } from 'path';
import * as yaml from 'yaml';

export async function projectApiDoc(format: string = 'both', outputDir: string = 'docs/api') {
  try {
    const backendPath = join(process.cwd(), 'backend');
    const apiDocResults = {
      endpoints: [],
      models: [],
      openapi: null,
      markdown: '',
      outputPath: ''
    };

    // バックエンドディレクトリの存在確認
    try {
      await readdir(backendPath);
    } catch {
      return { error: 'backendディレクトリが見つかりません' };
    }

    // 出力ディレクトリ作成
    const fullOutputDir = join(process.cwd(), outputDir);
    await mkdir(fullOutputDir, { recursive: true });

    // Goファイルからエンドポイント情報を抽出
    apiDocResults.endpoints = await extractEndpoints(backendPath);

    // モデル情報を抽出
    apiDocResults.models = await extractModels(backendPath);

    // OpenAPI仕様書生成
    if (format === 'openapi' || format === 'both') {
      apiDocResults.openapi = generateOpenAPISpec(apiDocResults.endpoints, apiDocResults.models);
      const openApiPath = join(fullOutputDir, 'openapi.yaml');
      await writeFile(openApiPath, yaml.stringify(apiDocResults.openapi));
      apiDocResults.outputPath += `OpenAPI: ${openApiPath}\n`;
    }

    // Markdownドキュメント生成
    if (format === 'markdown' || format === 'both') {
      apiDocResults.markdown = generateMarkdownDoc(apiDocResults.endpoints, apiDocResults.models);
      const markdownPath = join(fullOutputDir, 'api-documentation.md');
      await writeFile(markdownPath, apiDocResults.markdown);
      apiDocResults.outputPath += `Markdown: ${markdownPath}\n`;
    }

    return {
      success: true,
      endpoints: apiDocResults.endpoints.length,
      models: apiDocResults.models.length,
      outputPath: apiDocResults.outputPath.trim(),
      summary: generateSummary(apiDocResults)
    };

  } catch (error) {
    console.error('API documentation generation error:', error);
    return { error: 'APIドキュメント生成中にエラーが発生しました: ' + error.message };
  }
}

async function extractEndpoints(backendPath: string): Promise<APIEndpoint[]> {
  const endpoints: APIEndpoint[] = [];
  const handlerFiles = await findGoFiles(join(backendPath, 'internal/handler'));

  for (const file of handlerFiles) {
    try {
      const content = await readFile(file, 'utf-8');
      const extractedEndpoints = parseHandlerFile(content, file);
      endpoints.push(...extractedEndpoints);
    } catch (error) {
      console.error(`Error parsing handler file ${file}:`, error);
    }
  }

  return endpoints;
}

async function extractModels(backendPath: string): Promise<APIModel[]> {
  const models: APIModel[] = [];
  const modelFiles = await findGoFiles(join(backendPath, 'internal/model'));

  for (const file of modelFiles) {
    try {
      const content = await readFile(file, 'utf-8');
      const extractedModels = parseModelFile(content, file);
      models.push(...extractedModels);
    } catch (error) {
      console.error(`Error parsing model file ${file}:`, error);
    }
  }

  return models;
}

function parseHandlerFile(content: string, filePath: string): APIEndpoint[] {
  const endpoints: APIEndpoint[] = [];

  // HTTP ハンドラー関数を検索
  const handlerRegex = /func\s+\(\w+\s+\*\w+\)\s+(\w+)\(w\s+http\.ResponseWriter,\s+r\s+\*http\.Request\)/g;
  const routeRegex = /r\.(\w+)\("([^"]+)",\s+\w+\.(\w+)\)/g;

  let handlerMatch;
  while ((handlerMatch = handlerRegex.exec(content)) !== null) {
    const handlerName = handlerMatch[1];

    // ハンドラーの開始位置から関数本体を取得
    const funcStart = handlerMatch.index;
    const funcBody = extractFunctionBody(content, funcStart);

    // コメントからAPI情報を抽出
    const commentMatch = content.substring(0, funcStart).match(/\/\/.*\n/g);
    const lastComment = commentMatch ? commentMatch[commentMatch.length - 1] : '';

    endpoints.push({
      name: handlerName,
      method: inferHTTPMethod(handlerName),
      path: inferPath(handlerName),
      description: lastComment.replace(/^\/\/\s*/, '').trim(),
      requestModel: extractRequestModel(funcBody),
      responseModel: extractResponseModel(funcBody),
      file: filePath
    });
  }

  // ルーティング情報を追加で解析
  let routeMatch;
  while ((routeMatch = routeRegex.exec(content)) !== null) {
    const method = routeMatch[1].toUpperCase();
    const path = routeMatch[2];
    const handlerName = routeMatch[3];

    // 既存のエンドポイント情報を更新
    const existingEndpoint = endpoints.find(e => e.name === handlerName);
    if (existingEndpoint) {
      existingEndpoint.method = method;
      existingEndpoint.path = path;
    }
  }

  return endpoints;
}

function parseModelFile(content: string, filePath: string): APIModel[] {
  const models: APIModel[] = [];

  // 構造体定義を検索
  const structRegex = /type\s+(\w+)\s+struct\s*\{([^}]+)\}/g;

  let structMatch;
  while ((structMatch = structRegex.exec(content)) !== null) {
    const structName = structMatch[1];
    const structBody = structMatch[2];

    // フィールド情報を解析
    const fields = parseStructFields(structBody);

    // コメントから説明を抽出
    const commentMatch = content.substring(0, structMatch.index).match(/\/\/.*\n/g);
    const lastComment = commentMatch ? commentMatch[commentMatch.length - 1] : '';

    models.push({
      name: structName,
      description: lastComment.replace(/^\/\/\s*/, '').trim(),
      fields: fields,
      file: filePath
    });
  }

  return models;
}

function parseStructFields(structBody: string): ModelField[] {
  const fields: ModelField[] = [];
  const fieldLines = structBody.split('\n');

  for (const line of fieldLines) {
    const trimmedLine = line.trim();
    if (!trimmedLine || trimmedLine.startsWith('//')) continue;

    // フィールド定義をパース: FieldName Type `tags`
    const fieldMatch = trimmedLine.match(/(\w+)\s+([^\s`]+)(?:\s+`([^`]*)`)?/);
    if (fieldMatch) {
      const fieldName = fieldMatch[1];
      const fieldType = fieldMatch[2];
      const tags = fieldMatch[3] || '';

      // JSONタグから情報を抽出
      const jsonTag = extractJSONTag(tags);
      const gormTag = extractGORMTag(tags);
      const validateTag = extractValidateTag(tags);

      fields.push({
        name: fieldName,
        type: fieldType,
        jsonName: jsonTag.name,
        required: validateTag.includes('required'),
        description: extractFieldComment(line),
        constraints: validateTag
      });
    }
  }

  return fields;
}

function generateOpenAPISpec(endpoints: APIEndpoint[], models: APIModel[]): any {
  const spec = {
    openapi: '3.0.0',
    info: {
      title: 'AIDD Application API',
      description: 'AI-Driven Development アプリケーションのAPI仕様書',
      version: '1.0.0',
      contact: {
        name: 'AIDD Team',
        email: 'team@aidd.example.com'
      }
    },
    servers: [
      {
        url: 'http://localhost:8080',
        description: '開発環境'
      },
      {
        url: 'https://api.aidd.example.com',
        description: '本番環境'
      }
    ],
    paths: {},
    components: {
      schemas: {},
      securitySchemes: {
        bearerAuth: {
          type: 'http',
          scheme: 'bearer',
          bearerFormat: 'JWT'
        }
      }
    }
  };

  // エンドポイント情報をOpenAPI形式に変換
  for (const endpoint of endpoints) {
    if (!spec.paths[endpoint.path]) {
      spec.paths[endpoint.path] = {};
    }

    spec.paths[endpoint.path][endpoint.method.toLowerCase()] = {
      summary: endpoint.name,
      description: endpoint.description || `${endpoint.method} ${endpoint.path}`,
      tags: [inferTag(endpoint.path)],
      security: endpoint.method !== 'POST' || !endpoint.path.includes('/auth') ? [{ bearerAuth: [] }] : [],
      parameters: extractPathParameters(endpoint.path),
      requestBody: endpoint.requestModel ? generateRequestBody(endpoint.requestModel) : undefined,
      responses: {
        '200': {
          description: 'Success',
          content: {
            'application/json': {
              schema: endpoint.responseModel ? { $ref: `#/components/schemas/${endpoint.responseModel}` } : {}
            }
          }
        },
        '400': {
          description: 'Bad Request',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/ErrorResponse' }
            }
          }
        },
        '401': {
          description: 'Unauthorized',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/ErrorResponse' }
            }
          }
        },
        '500': {
          description: 'Internal Server Error',
          content: {
            'application/json': {
              schema: { $ref: '#/components/schemas/ErrorResponse' }
            }
          }
        }
      }
    };
  }

  // モデル情報をOpenAPI形式に変換
  for (const model of models) {
    spec.components.schemas[model.name] = {
      type: 'object',
      description: model.description,
      properties: {},
      required: []
    };

    for (const field of model.fields) {
      spec.components.schemas[model.name].properties[field.jsonName || field.name] = {
        type: mapGoTypeToOpenAPI(field.type),
        description: field.description
      };

      if (field.required) {
        spec.components.schemas[model.name].required.push(field.jsonName || field.name);
      }
    }
  }

  // 共通エラーレスポンススキーマを追加
  spec.components.schemas.ErrorResponse = {
    type: 'object',
    properties: {
      code: { type: 'string', description: 'エラーコード' },
      message: { type: 'string', description: 'エラーメッセージ' },
      details: { type: 'string', description: 'エラー詳細' }
    },
    required: ['code', 'message']
  };

  return spec;
}

function generateMarkdownDoc(endpoints: APIEndpoint[], models: APIModel[]): string {
  let markdown = `# API Documentation

> AIDD Application REST API 仕様書

## 概要

このドキュメントは、AIDD（AI-Driven Development）アプリケーションのREST API仕様を記載しています。

### 基本情報

- **Base URL**: \`http://localhost:8080\`
- **認証方式**: Bearer Token (JWT)
- **データ形式**: JSON
- **文字エンコーディング**: UTF-8

## 認証

多くのAPIエンドポイントでは認証が必要です。認証にはJWTトークンを使用します。

\`\`\`http
Authorization: Bearer <your-jwt-token>
\`\`\`

## エンドポイント一覧

`;

  // エンドポイント情報をMarkdown形式で生成
  const groupedEndpoints = groupEndpointsByTag(endpoints);

  for (const [tag, tagEndpoints] of Object.entries(groupedEndpoints)) {
    markdown += `\n### ${tag}\n\n`;

    for (const endpoint of tagEndpoints) {
      markdown += `#### ${endpoint.method} ${endpoint.path}\n\n`;
      markdown += `${endpoint.description}\n\n`;

      // パラメータ
      const pathParams = extractPathParameters(endpoint.path);
      if (pathParams.length > 0) {
        markdown += `**パスパラメータ:**\n\n`;
        for (const param of pathParams) {
          markdown += `- \`${param.name}\` (${param.schema?.type || 'string'}) - ${param.description || 'ID'}\n`;
        }
        markdown += '\n';
      }

      // リクエストボディ
      if (endpoint.requestModel) {
        markdown += `**リクエストボディ:**\n\n`;
        markdown += `\`\`\`json\n`;
        markdown += generateExampleJSON(endpoint.requestModel, models);
        markdown += `\n\`\`\`\n\n`;
      }

      // レスポンス
      markdown += `**レスポンス例:**\n\n`;
      markdown += `\`\`\`json\n`;
      if (endpoint.responseModel) {
        markdown += generateExampleJSON(endpoint.responseModel, models);
      } else {
        markdown += `{\n  "message": "Success"\n}`;
      }
      markdown += `\n\`\`\`\n\n`;

      markdown += `---\n\n`;
    }
  }

  // モデル定義
  if (models.length > 0) {
    markdown += `## データモデル\n\n`;

    for (const model of models) {
      markdown += `### ${model.name}\n\n`;
      if (model.description) {
        markdown += `${model.description}\n\n`;
      }

      markdown += `| フィールド名 | 型 | 必須 | 説明 |\n`;
      markdown += `|--------------|-----|------|------|\n`;

      for (const field of model.fields) {
        const required = field.required ? '✅' : '';
        markdown += `| ${field.jsonName || field.name} | ${field.type} | ${required} | ${field.description || '-'} |\n`;
      }

      markdown += '\n';
    }
  }

  // エラーレスポンス
  markdown += `## エラーレスポンス

APIエラーの場合、以下の形式でエラー情報が返されます：

\`\`\`json
{
  "code": "ERROR_CODE",
  "message": "エラーメッセージ",
  "details": "詳細な説明（オプション）"
}
\`\`\`

### ステータスコード

| コード | 説明 |
|--------|------|
| 200 | 成功 |
| 201 | 作成成功 |
| 400 | リクエストエラー |
| 401 | 認証エラー |
| 403 | 権限エラー |
| 404 | リソース未発見 |
| 409 | 競合エラー |
| 500 | サーバーエラー |

---

*この仕様書は \`/project:api-doc\` コマンドにより自動生成されました*
`;

  return markdown;
}

// ヘルパー関数
async function findGoFiles(dir: string): Promise<string[]> {
  const files: string[] = [];
  try {
    const entries = await readdir(dir, { withFileTypes: true });
    for (const entry of entries) {
      const fullPath = join(dir, entry.name);
      if (entry.isDirectory()) {
        files.push(...await findGoFiles(fullPath));
      } else if (entry.isFile() && extname(entry.name) === '.go') {
        files.push(fullPath);
      }
    }
  } catch {
    // ディレクトリが存在しない場合は空配列を返す
  }
  return files;
}

function inferHTTPMethod(handlerName: string): string {
  const name = handlerName.toLowerCase();
  if (name.includes('create') || name.includes('post')) return 'POST';
  if (name.includes('update') || name.includes('put')) return 'PUT';
  if (name.includes('delete')) return 'DELETE';
  if (name.includes('get') || name.includes('list')) return 'GET';
  return 'GET';
}

function inferPath(handlerName: string): string {
  const name = handlerName.toLowerCase();
  if (name.includes('user')) return '/api/users';
  if (name.includes('post')) return '/api/posts';
  if (name.includes('auth') || name.includes('login')) return '/api/auth';
  return `/api/${name}`;
}

function inferTag(path: string): string {
  if (path.includes('/users')) return 'Users';
  if (path.includes('/posts')) return 'Posts';
  if (path.includes('/auth')) return 'Authentication';
  return 'General';
}

function mapGoTypeToOpenAPI(goType: string): string {
  const typeMap: { [key: string]: string } = {
    'string': 'string',
    'int': 'integer',
    'int32': 'integer',
    'int64': 'integer',
    'uint': 'integer',
    'uint32': 'integer',
    'uint64': 'integer',
    'float32': 'number',
    'float64': 'number',
    'bool': 'boolean',
    'time.Time': 'string'
  };

  return typeMap[goType] || 'string';
}

// 型定義
interface APIEndpoint {
  name: string;
  method: string;
  path: string;
  description: string;
  requestModel?: string;
  responseModel?: string;
  file: string;
}

interface APIModel {
  name: string;
  description: string;
  fields: ModelField[];
  file: string;
}

interface ModelField {
  name: string;
  type: string;
  jsonName?: string;
  required: boolean;
  description?: string;
  constraints?: string;
}
```

## Usage Example

```bash
# 全形式でAPIドキュメント生成
/project:api-doc

# OpenAPI仕様書のみ生成
/project:api-doc openapi

# Markdownドキュメントのみ生成
/project:api-doc markdown

# カスタム出力ディレクトリ
/project:api-doc both docs/api-reference
```

## Expected Output

```markdown
## 📚 APIドキュメント生成結果

### ✅ 生成完了
- **エンドポイント数**: 15個
- **データモデル数**: 8個
- **出力ファイル**:
  - OpenAPI: docs/api/openapi.yaml
  - Markdown: docs/api/api-documentation.md

### 📋 検出されたエンドポイント

#### Authentication
- POST /api/auth/login - ユーザーログイン
- POST /api/auth/register - ユーザー登録
- POST /api/auth/refresh - トークンリフレッシュ

#### Users
- GET /api/users - ユーザー一覧取得
- POST /api/users - ユーザー作成
- GET /api/users/{id} - ユーザー詳細取得
- PUT /api/users/{id} - ユーザー更新
- DELETE /api/users/{id} - ユーザー削除

#### Posts
- GET /api/posts - 投稿一覧取得
- POST /api/posts - 投稿作成
- GET /api/posts/{id} - 投稿詳細取得
- PUT /api/posts/{id} - 投稿更新
- DELETE /api/posts/{id} - 投稿削除

### 🏗️ データモデル
- User - ユーザー情報
- Post - 投稿情報
- Profile - プロフィール情報
- CreateUserRequest - ユーザー作成リクエスト
- UpdateUserRequest - ユーザー更新リクエスト
- LoginRequest - ログインリクエスト
- AuthResponse - 認証レスポンス
- ErrorResponse - エラーレスポンス

### 🚀 次のステップ
1. 生成されたOpenAPI仕様書をSwagger UIで確認
2. フロントエンドチームに型定義を共有
3. Postmanコレクションを生成してAPIテストに活用

### 📖 生成されたファイルの活用方法
```bash
# Swagger UIでOpenAPI仕様書を表示
npx swagger-ui-serve docs/api/openapi.yaml

# TypeScript型定義生成
npx openapi-typescript docs/api/openapi.yaml --output src/types/api.ts

# Postmanコレクション生成
npx openapi-to-postman -s docs/api/openapi.yaml -o docs/api/postman-collection.json
```

---
*この結果は `/project:api-doc` コマンドにより自動生成されました*
```