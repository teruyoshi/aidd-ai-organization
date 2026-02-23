# /project:review

## Description
コードレビューを実行し、品質・セキュリティ・パフォーマンス・保守性を総合的にチェックするスラッシュコマンド

## Parameters
- `target` (optional): レビュー対象（all | frontend | backend | specific-file-path）

## Implementation

```typescript
import { readdir, readFile, stat } from 'fs/promises';
import { join, extname } from 'path';

export async function projectReview(target: string = 'all') {
  const reviewResults = {
    summary: {},
    frontend: {},
    backend: {},
    infrastructure: {},
    issues: [],
    recommendations: []
  };

  try {
    if (target === 'all' || target === 'frontend') {
      reviewResults.frontend = await reviewFrontend();
    }

    if (target === 'all' || target === 'backend') {
      reviewResults.backend = await reviewBackend();
    }

    if (target === 'all') {
      reviewResults.infrastructure = await reviewInfrastructure();
    }

    // 特定ファイルの場合
    if (target !== 'all' && target !== 'frontend' && target !== 'backend') {
      reviewResults.fileReview = await reviewSpecificFile(target);
    }

    reviewResults.summary = generateSummary(reviewResults);
    reviewResults.issues = aggregateIssues(reviewResults);
    reviewResults.recommendations = generateRecommendations(reviewResults);

    return reviewResults;

  } catch (error) {
    console.error('Code review error:', error);
    return { error: 'コードレビュー中にエラーが発生しました: ' + error.message };
  }
}

async function reviewFrontend() {
  const frontendPath = join(process.cwd(), 'frontend');
  const results = {
    typeScript: { score: 0, issues: [] },
    react: { score: 0, issues: [] },
    performance: { score: 0, issues: [] },
    accessibility: { score: 0, issues: [] },
    security: { score: 0, issues: [] },
    testCoverage: { score: 0, issues: [] }
  };

  try {
    // TypeScript品質チェック
    results.typeScript = await checkTypeScriptQuality(frontendPath);

    // React品質チェック
    results.react = await checkReactQuality(frontendPath);

    // パフォーマンスチェック
    results.performance = await checkFrontendPerformance(frontendPath);

    // アクセシビリティチェック
    results.accessibility = await checkAccessibility(frontendPath);

    // セキュリティチェック
    results.security = await checkFrontendSecurity(frontendPath);

    // テストカバレッジチェック
    results.testCoverage = await checkTestCoverage(frontendPath);

  } catch (error) {
    results.error = error.message;
  }

  return results;
}

async function reviewBackend() {
  const backendPath = join(process.cwd(), 'backend');
  const results = {
    go: { score: 0, issues: [] },
    architecture: { score: 0, issues: [] },
    database: { score: 0, issues: [] },
    api: { score: 0, issues: [] },
    security: { score: 0, issues: [] },
    testCoverage: { score: 0, issues: [] }
  };

  try {
    // Go品質チェック
    results.go = await checkGoQuality(backendPath);

    // アーキテクチャチェック
    results.architecture = await checkArchitecture(backendPath);

    // データベース設計チェック
    results.database = await checkDatabaseDesign(backendPath);

    // API設計チェック
    results.api = await checkAPIDesign(backendPath);

    // セキュリティチェック
    results.security = await checkBackendSecurity(backendPath);

    // テストカバレッジチェック
    results.testCoverage = await checkGoTestCoverage(backendPath);

  } catch (error) {
    results.error = error.message;
  }

  return results;
}

async function reviewInfrastructure() {
  const results = {
    docker: { score: 0, issues: [] },
    nginx: { score: 0, issues: [] },
    security: { score: 0, issues: [] },
    monitoring: { score: 0, issues: [] }
  };

  try {
    // Docker設定チェック
    results.docker = await checkDockerConfiguration();

    // Nginx設定チェック
    results.nginx = await checkNginxConfiguration();

    // インフラセキュリティチェック
    results.security = await checkInfraSecurity();

    // 監視設定チェック
    results.monitoring = await checkMonitoring();

  } catch (error) {
    results.error = error.message;
  }

  return results;
}

async function checkTypeScriptQuality(frontendPath: string) {
  const issues = [];
  let score = 100;

  try {
    // tsconfig.json チェック
    const tsconfigPath = join(frontendPath, 'tsconfig.json');
    const tsconfigContent = await readFile(tsconfigPath, 'utf-8');
    const tsconfig = JSON.parse(tsconfigContent);

    // strict mode チェック
    if (!tsconfig.compilerOptions?.strict) {
      issues.push('❌ TypeScript strict modeが無効です');
      score -= 20;
    }

    // any 使用チェック
    const tsFiles = await findFiles(frontendPath, ['.ts', '.tsx']);
    for (const file of tsFiles) {
      const content = await readFile(file, 'utf-8');
      const anyMatches = content.match(/:\s*any\b/g);
      if (anyMatches && anyMatches.length > 0) {
        issues.push(`⚠️ ${file}: any型が${anyMatches.length}箇所使用されています`);
        score -= anyMatches.length * 5;
      }
    }

    // 型定義ファイルチェック
    const typesDir = join(frontendPath, 'src', 'types');
    const hasTyepDefines = await directoryExists(typesDir);
    if (!hasTyepDefines) {
      issues.push('⚠️ 型定義ディレクトリ(src/types)が見つかりません');
      score -= 10;
    }

  } catch (error) {
    issues.push(`❌ TypeScript設定チェックでエラー: ${error.message}`);
    score = 0;
  }

  return { score: Math.max(0, score), issues };
}

async function checkReactQuality(frontendPath: string) {
  const issues = [];
  let score = 100;

  try {
    const jsxFiles = await findFiles(frontendPath, ['.jsx', '.tsx']);

    for (const file of jsxFiles) {
      const content = await readFile(file, 'utf-8');

      // useState/useEffect のベストプラクティスチェック
      const useEffectMatches = content.match(/useEffect\([^,]+,\s*\[\]/g);
      if (useEffectMatches && useEffectMatches.length > 0) {
        // 空の依存配列をチェック（必要に応じて）
      }

      // インラインスタイルチェック
      const inlineStyleMatches = content.match(/style=\{\{/g);
      if (inlineStyleMatches && inlineStyleMatches.length > 0) {
        issues.push(`⚠️ ${file}: インラインスタイルが使用されています`);
        score -= 5;
      }

      // key prop チェック
      const mapWithoutKey = content.match(/\.map\([^}]+\)\s*(?!.*key=)/g);
      if (mapWithoutKey && mapWithoutKey.length > 0) {
        issues.push(`❌ ${file}: .map()でkey propが不足している可能性があります`);
        score -= 15;
      }
    }

  } catch (error) {
    issues.push(`❌ React品質チェックでエラー: ${error.message}`);
    score = 0;
  }

  return { score: Math.max(0, score), issues };
}

async function checkGoQuality(backendPath: string) {
  const issues = [];
  let score = 100;

  try {
    const goFiles = await findFiles(backendPath, ['.go']);

    for (const file of goFiles) {
      const content = await readFile(file, 'utf-8');

      // エラーハンドリングチェック
      const errorChecks = content.match(/if err != nil/g);
      const errorUsage = content.match(/\berr\b/g);
      if (errorUsage && (!errorChecks || errorChecks.length < errorUsage.length / 3)) {
        issues.push(`⚠️ ${file}: エラーハンドリングが不足している可能性があります`);
        score -= 10;
      }

      // 適切なパッケージ構造チェック
      if (content.includes('package main') && !file.includes('cmd/')) {
        issues.push(`❌ ${file}: main packageはcmd/配下に配置してください`);
        score -= 15;
      }

      // コメント不足チェック
      const publicFunctions = content.match(/^func\s+[A-Z][a-zA-Z0-9_]*\(/gm);
      const comments = content.match(/^\/\/\s*[A-Z][a-zA-Z0-9_]*\s/gm);
      if (publicFunctions && (!comments || comments.length < publicFunctions.length)) {
        issues.push(`⚠️ ${file}: 公開関数のコメントが不足しています`);
        score -= 5;
      }
    }

  } catch (error) {
    issues.push(`❌ Go品質チェックでエラー: ${error.message}`);
    score = 0;
  }

  return { score: Math.max(0, score), issues };
}

async function checkDockerConfiguration() {
  const issues = [];
  let score = 100;

  try {
    // docker-compose.yml チェック
    const dockerComposePath = join(process.cwd(), 'docker-compose.yml');
    const dockerComposeExists = await fileExists(dockerComposePath);

    if (!dockerComposeExists) {
      issues.push('❌ docker-compose.ymlが見つかりません');
      return { score: 0, issues };
    }

    const composeContent = await readFile(dockerComposePath, 'utf-8');

    // ヘルスチェック設定チェック
    if (!composeContent.includes('healthcheck:')) {
      issues.push('⚠️ サービスのヘルスチェック設定がありません');
      score -= 20;
    }

    // リソース制限チェック
    if (!composeContent.includes('deploy:') || !composeContent.includes('resources:')) {
      issues.push('⚠️ リソース制限設定がありません');
      score -= 15;
    }

    // ネットワーク設定チェック
    if (!composeContent.includes('networks:')) {
      issues.push('⚠️ カスタムネットワーク設定がありません');
      score -= 10;
    }

  } catch (error) {
    issues.push(`❌ Docker設定チェックでエラー: ${error.message}`);
    score = 0;
  }

  return { score: Math.max(0, score), issues };
}

// ヘルパー関数
async function findFiles(dir: string, extensions: string[]): Promise<string[]> {
  const files: string[] = [];

  try {
    const entries = await readdir(dir, { withFileTypes: true });

    for (const entry of entries) {
      const fullPath = join(dir, entry.name);

      if (entry.isDirectory() && !entry.name.startsWith('.') && entry.name !== 'node_modules') {
        files.push(...await findFiles(fullPath, extensions));
      } else if (entry.isFile() && extensions.includes(extname(entry.name))) {
        files.push(fullPath);
      }
    }
  } catch (error) {
    // ディレクトリが存在しない場合は空配列を返す
  }

  return files;
}

async function fileExists(path: string): Promise<boolean> {
  try {
    await stat(path);
    return true;
  } catch {
    return false;
  }
}

async function directoryExists(path: string): Promise<boolean> {
  try {
    const stats = await stat(path);
    return stats.isDirectory();
  } catch {
    return false;
  }
}

function generateSummary(results: any) {
  const scores = [];

  if (results.frontend && typeof results.frontend === 'object') {
    const frontendScores = Object.values(results.frontend).map((r: any) => r.score || 0);
    scores.push(...frontendScores);
  }

  if (results.backend && typeof results.backend === 'object') {
    const backendScores = Object.values(results.backend).map((r: any) => r.score || 0);
    scores.push(...backendScores);
  }

  const averageScore = scores.length > 0 ? Math.round(scores.reduce((a, b) => a + b, 0) / scores.length) : 0;

  return {
    overallScore: averageScore,
    grade: averageScore >= 90 ? 'A' : averageScore >= 80 ? 'B' : averageScore >= 70 ? 'C' : 'D',
    totalIssues: aggregateIssues(results).length
  };
}

function aggregateIssues(results: any): string[] {
  const allIssues: string[] = [];

  function collectIssues(obj: any) {
    if (obj && typeof obj === 'object') {
      if (Array.isArray(obj.issues)) {
        allIssues.push(...obj.issues);
      }
      for (const value of Object.values(obj)) {
        if (typeof value === 'object') {
          collectIssues(value);
        }
      }
    }
  }

  collectIssues(results);
  return allIssues;
}

function generateRecommendations(results: any): string[] {
  const recommendations: string[] = [];
  const summary = results.summary;

  if (summary.overallScore < 70) {
    recommendations.push('🔥 コード品質が低下しています。緊急対応が必要です');
  } else if (summary.overallScore < 85) {
    recommendations.push('⚠️ コード品質の改善が必要です');
  } else {
    recommendations.push('✅ コード品質は良好です');
  }

  if (summary.totalIssues > 10) {
    recommendations.push('📋 課題が多いため、優先順位をつけて対応してください');
  }

  // エージェント別推奨事項
  if (results.frontend?.typeScript?.score < 80) {
    recommendations.push('🎨 フロントエンドエージェント: TypeScript型安全性の向上が必要です');
  }

  if (results.backend?.go?.score < 80) {
    recommendations.push('⚙️ バックエンドエージェント: Go コード品質の向上が必要です');
  }

  if (results.infrastructure?.docker?.score < 80) {
    recommendations.push('🐳 DevOpsエージェント: Docker設定の改善が必要です');
  }

  return recommendations;
}
```

## Usage Example

```bash
# 全体レビュー
/project:review

# フロントエンドのみレビュー
/project:review frontend

# バックエンドのみレビュー
/project:review backend

# 特定ファイルレビュー
/project:review src/components/UserCard.tsx
```

## Expected Output

```markdown
## 📊 コードレビュー結果

### 🎯 総合評価
**スコア: 87/100 (Grade: B)**
検出された課題: 12件

### 📋 フロントエンド品質
- **TypeScript**: 85/100 ✅
- **React**: 92/100 ✅
- **パフォーマンス**: 78/100 ⚠️
- **アクセシビリティ**: 88/100 ✅
- **セキュリティ**: 95/100 ✅
- **テストカバレッジ**: 82/100 ✅

### ⚙️ バックエンド品質
- **Go言語**: 90/100 ✅
- **アーキテクチャ**: 88/100 ✅
- **データベース**: 85/100 ✅
- **API設計**: 90/100 ✅
- **セキュリティ**: 93/100 ✅
- **テストカバレッジ**: 79/100 ⚠️

### 🐳 インフラ品質
- **Docker**: 85/100 ✅
- **Nginx**: 92/100 ✅
- **セキュリティ**: 88/100 ✅
- **監視**: 75/100 ⚠️

### ❌ 重要な課題
- src/components/UserList.tsx: .map()でkey propが不足している可能性があります
- backend/internal/handler/user.go: エラーハンドリングが不足している可能性があります
- docker-compose.yml: リソース制限設定がありません

### 💡 改善推奨事項
✅ コード品質は良好です
📋 課題が多いため、優先順位をつけて対応してください
🎨 フロントエンドエージェント: パフォーマンス最適化に注目してください
⚙️ バックエンドエージェント: テストカバレッジの向上が必要です
🐳 DevOpsエージェント: 監視設定の強化をお勧めします

### 🚀 次のアクション
```bash
# パフォーマンス最適化
Read agents/frontend.md and optimize component rendering performance

# テストカバレッジ向上
Read agents/qa.md and increase backend test coverage to 85%

# 監視設定強化
Read agents/devops.md and docs/context/devops-context.md, then enhance monitoring configuration
```

---
*このレビューは `/project:review` コマンドにより自動生成されました*
```