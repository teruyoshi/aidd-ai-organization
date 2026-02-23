# /project:plan

## Description
SCHEDULE.mdを確認して次のタスクを分解・提示するスラッシュコマンド

## Parameters
None

## Implementation

```typescript
import { readFile } from 'fs/promises';
import { join } from 'path';

export async function projectPlan() {
  try {
    // SCHEDULE.mdを読み込み
    const schedulePath = join(process.cwd(), 'SCHEDULE.md');
    const scheduleContent = await readFile(schedulePath, 'utf-8');

    // 現在の進捗状況を分析
    const currentPhase = analyzeCurrentPhase(scheduleContent);
    const nextTasks = extractNextTasks(scheduleContent);
    const blockers = findBlockers(scheduleContent);

    // レスポンス生成
    return {
      currentPhase,
      nextTasks,
      blockers,
      recommendations: generateRecommendations(currentPhase, nextTasks, blockers)
    };

  } catch (error) {
    console.error('Error reading SCHEDULE.md:', error);
    return { error: 'SCHEDULE.mdが見つかりません。まずSCHEDULE.mdを作成してください。' };
  }
}

function analyzeCurrentPhase(content: string) {
  const lines = content.split('\n');
  const phases = ['Phase 1: 基盤整備', 'Phase 2: バックエンド開発', 'Phase 3: フロントエンド開発', 'Phase 4: 統合・テスト', 'Phase 5: デプロイ・運用'];

  let currentPhase = 'Phase 1: 基盤整備';
  let completedTasks = 0;
  let totalTasks = 0;

  for (const line of lines) {
    if (line.includes('[x]')) completedTasks++;
    if (line.includes('[-]') || line.includes('[ ]') || line.includes('[x]')) totalTasks++;

    for (const phase of phases) {
      if (line.includes(phase)) {
        // このフェーズのタスクが進行中or未着手なら現在のフェーズ
        const phaseContent = extractPhaseContent(content, phase);
        if (phaseContent.includes('[-]') || phaseContent.includes('[ ]')) {
          currentPhase = phase;
          break;
        }
      }
    }
  }

  const progress = totalTasks > 0 ? Math.round((completedTasks / totalTasks) * 100) : 0;

  return {
    name: currentPhase,
    progress: `${progress}% (${completedTasks}/${totalTasks})`
  };
}

function extractNextTasks(content: string): string[] {
  const lines = content.split('\n');
  const nextTasks: string[] = [];

  // 進行中のタスクを抽出
  for (const line of lines) {
    if (line.includes('[-]')) {
      const task = line.replace(/^\s*-\s*\[[-]\]\s*/, '').trim();
      if (task) nextTasks.push(`🏗️ ${task}`);
    }
  }

  // 未着手のタスク（最大5つ）を抽出
  let unfinishedCount = 0;
  for (const line of lines) {
    if (line.includes('[ ]') && unfinishedCount < 5) {
      const task = line.replace(/^\s*-\s*\[\s*\]\s*/, '').trim();
      if (task) {
        nextTasks.push(`📋 ${task}`);
        unfinishedCount++;
      }
    }
  }

  return nextTasks;
}

function findBlockers(content: string): string[] {
  const blockers: string[] = [];

  // ブロッカーキーワードを探す
  const blockerKeywords = ['ブロッカー', 'blocker', '課題', '問題', 'issue', '依存', 'dependency'];
  const lines = content.split('\n');

  for (const line of lines) {
    for (const keyword of blockerKeywords) {
      if (line.toLowerCase().includes(keyword)) {
        blockers.push(line.trim());
        break;
      }
    }
  }

  return blockers;
}

function generateRecommendations(currentPhase: any, nextTasks: string[], blockers: string[]): string[] {
  const recommendations: string[] = [];

  // フェーズ別推奨事項
  if (currentPhase.name.includes('Phase 1')) {
    recommendations.push('🏗️ まずはDocker開発環境を構築することをお勧めします');
    recommendations.push('📋 データベース設計とAPI仕様の策定が重要です');
    recommendations.push('🔧 各エージェントの技術スタック理解を深めましょう');
  } else if (currentPhase.name.includes('Phase 2')) {
    recommendations.push('⚙️ バックエンドAPIの実装に集中しましょう');
    recommendations.push('🗄️ データベースマイグレーションの準備が必要です');
    recommendations.push('🧪 API単体テストの並行実装をお勧めします');
  } else if (currentPhase.name.includes('Phase 3')) {
    recommendations.push('🎨 フロントエンドコンポーネントの実装を進めましょう');
    recommendations.push('🔗 バックエンドAPIとの連携テストが重要です');
    recommendations.push('📱 レスポンシブデザインの確認をお忘れなく');
  } else if (currentPhase.name.includes('Phase 4')) {
    recommendations.push('🧪 E2Eテストの実装と実行が最優先です');
    recommendations.push('📊 パフォーマンステストの実施をお勧めします');
    recommendations.push('🔒 セキュリティテストの実行をお忘れなく');
  } else if (currentPhase.name.includes('Phase 5')) {
    recommendations.push('🚀 本番環境の構築・設定を進めましょう');
    recommendations.push('📈 監視・ログ設定の実装が重要です');
    recommendations.push('📚 運用ドキュメントの整備をお勧めします');
  }

  // ブロッカーがある場合の推奨事項
  if (blockers.length > 0) {
    recommendations.push('⚠️ ブロッカーの解決が最優先です');
    recommendations.push('🤝 必要に応じて他のエージェントとの連携を検討してください');
  }

  // タスク数による推奨事項
  if (nextTasks.length === 0) {
    recommendations.push('🎉 全タスク完了おめでとうございます！次のフェーズに進みましょう');
  } else if (nextTasks.length > 10) {
    recommendations.push('📋 タスクが多いため、優先度付けをお勧めします');
    recommendations.push('👥 複数エージェントでの並行作業を検討してください');
  }

  return recommendations;
}

function extractPhaseContent(content: string, phaseName: string): string {
  const lines = content.split('\n');
  let inPhase = false;
  let phaseContent = '';

  for (const line of lines) {
    if (line.includes(phaseName)) {
      inPhase = true;
      continue;
    }

    if (inPhase) {
      if (line.startsWith('### ') && !line.includes(phaseName)) {
        break; // 次のフェーズに入ったら終了
      }
      phaseContent += line + '\n';
    }
  }

  return phaseContent;
}
```

## Usage Example

```bash
/project:plan
```

## Expected Output

```markdown
## 📊 プロジェクト進捗状況

### 🎯 現在のフェーズ
**Phase 2: バックエンド開発** (進捗: 45% - 9/20タスク完了)

### 📋 次に取り組むべきタスク
🏗️ Golang + go-chi/chi サーバー構築
📋 GORM + MySQL セットアップ
📋 認証・認可システム
📋 REST API実装
📋 単体テスト実装

### ⚠️ ブロッカー・課題
- データベース設計の最終確認が必要
- 認証方式（JWT vs Session）の決定待ち

### 💡 推奨アクション
⚙️ バックエンドAPIの実装に集中しましょう
🗄️ データベースマイグレーションの準備が必要です
🧪 API単体テストの並行実装をお勧めします

### 🚀 次のコマンド提案
```bash
# アーキテクトによる設計確認
Read agents/architect.md and finalize the authentication strategy

# バックエンド実装開始
Read agents/backend.md and docs/context/backend-context.md, then implement user authentication API

# DevOps環境準備
Read agents/devops.md and docs/context/devops-context.md, then create docker-compose.dev.yml
```

---
*このレポートは `/project:plan` コマンドにより自動生成されました*
```