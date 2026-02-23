# /project:migrate

## Description
データベースマイグレーションファイルを生成し、スキーマ変更を管理するスラッシュコマンド

## Parameters
- `action` (required): migrate | create | rollback | status
- `name` (optional): マイグレーション名（createアクション時）
- `steps` (optional): ロールバック時のステップ数（デフォルト: 1）

## Implementation

```typescript
import { readdir, readFile, writeFile, mkdir, stat } from 'fs/promises';
import { join, basename } from 'path';
import { execSync } from 'child_process';

export async function projectMigrate(action: string, name?: string, steps: number = 1) {
  const migrationsDir = join(process.cwd(), 'backend', 'migrations');
  const modelsDir = join(process.cwd(), 'backend', 'internal', 'model');

  try {
    // マイグレーションディレクトリ作成
    await mkdir(migrationsDir, { recursive: true });

    switch (action) {
      case 'create':
        return await createMigration(migrationsDir, name, modelsDir);
      case 'migrate':
        return await runMigrations(migrationsDir);
      case 'rollback':
        return await rollbackMigrations(migrationsDir, steps);
      case 'status':
        return await getMigrationStatus(migrationsDir);
      default:
        return { error: 'Invalid action. Use: create, migrate, rollback, or status' };
    }

  } catch (error) {
    console.error('Migration error:', error);
    return { error: 'マイグレーション実行中にエラーが発生しました: ' + error.message };
  }
}

async function createMigration(migrationsDir: string, name?: string, modelsDir?: string) {
  if (!name) {
    return { error: 'マイグレーション名が必要です。例: /project:migrate create create_users_table' };
  }

  const timestamp = new Date().toISOString().replace(/[-T:\.Z]/g, '').substring(0, 14);
  const fileName = `${timestamp}_${name}.go`;
  const filePath = join(migrationsDir, fileName);

  // 既存のモデルから情報を取得
  let modelInfo = '';
  if (modelsDir) {
    try {
      modelInfo = await analyzeModels(modelsDir);
    } catch (error) {
      console.warn('モデル分析でエラー:', error.message);
    }
  }

  // マイグレーションファイル生成
  const migrationContent = generateMigrationTemplate(name, modelInfo);

  await writeFile(filePath, migrationContent);

  // migrate.go ファイル更新
  await updateMigrateFile(migrationsDir, fileName, name);

  return {
    success: true,
    file: filePath,
    name: name,
    timestamp: timestamp,
    nextSteps: [
      '1. 生成されたマイグレーションファイルを編集してスキーマを定義',
      '2. /project:migrate migrate でマイグレーションを実行',
      '3. 必要に応じて backend/internal/model/ でモデルを更新'
    ]
  };
}

async function runMigrations(migrationsDir: string) {
  const results = {
    executed: [],
    skipped: [],
    errors: [],
    totalTime: 0
  };

  const startTime = Date.now();

  try {
    // マイグレーションファイル取得
    const migrationFiles = await getMigrationFiles(migrationsDir);

    // 実行済みマイグレーション確認
    const executedMigrations = await getExecutedMigrations();

    for (const file of migrationFiles) {
      const migrationId = basename(file, '.go');

      if (executedMigrations.includes(migrationId)) {
        results.skipped.push(migrationId);
        continue;
      }

      try {
        // マイグレーション実行
        await executeMigration(file);
        results.executed.push(migrationId);

        // 実行履歴記録
        await recordMigration(migrationId);

      } catch (error) {
        results.errors.push({
          migration: migrationId,
          error: error.message
        });
        break; // エラー時は処理停止
      }
    }

  } catch (error) {
    results.errors.push({
      migration: 'system',
      error: error.message
    });
  }

  results.totalTime = Date.now() - startTime;

  return {
    success: results.errors.length === 0,
    executed: results.executed.length,
    skipped: results.skipped.length,
    errors: results.errors.length,
    details: results,
    summary: generateMigrationSummary(results)
  };
}

async function rollbackMigrations(migrationsDir: string, steps: number) {
  const results = {
    rolledBack: [],
    errors: [],
    totalTime: 0
  };

  const startTime = Date.now();

  try {
    // 実行済みマイグレーション取得（最新から）
    const executedMigrations = await getExecutedMigrations();
    const toRollback = executedMigrations.slice(-steps);

    for (const migrationId of toRollback.reverse()) {
      try {
        const migrationFile = join(migrationsDir, `${migrationId}.go`);
        await rollbackMigration(migrationFile);
        results.rolledBack.push(migrationId);

        // 履歴から削除
        await removeMigrationRecord(migrationId);

      } catch (error) {
        results.errors.push({
          migration: migrationId,
          error: error.message
        });
        break;
      }
    }

  } catch (error) {
    results.errors.push({
      migration: 'system',
      error: error.message
    });
  }

  results.totalTime = Date.now() - startTime;

  return {
    success: results.errors.length === 0,
    rolledBack: results.rolledBack.length,
    errors: results.errors.length,
    details: results
  };
}

async function getMigrationStatus(migrationsDir: string) {
  try {
    const allMigrations = await getMigrationFiles(migrationsDir);
    const executedMigrations = await getExecutedMigrations();

    const status = allMigrations.map(file => {
      const migrationId = basename(file, '.go');
      const isExecuted = executedMigrations.includes(migrationId);

      return {
        id: migrationId,
        name: extractMigrationName(migrationId),
        file: file,
        status: isExecuted ? 'executed' : 'pending',
        executedAt: isExecuted ? getMigrationExecutionTime(migrationId) : null
      };
    });

    return {
      total: allMigrations.length,
      executed: executedMigrations.length,
      pending: allMigrations.length - executedMigrations.length,
      migrations: status
    };

  } catch (error) {
    return { error: 'マイグレーションステータス取得中にエラー: ' + error.message };
  }
}

function generateMigrationTemplate(name: string, modelInfo: string): string {
  const timestamp = new Date().toISOString();
  const migrationName = toPascalCase(name);

  return `package migrations

import (
    "gorm.io/gorm"
    "your-app/internal/model"
)

// ${migrationName} - ${name}
// Generated at: ${timestamp}
type ${migrationName} struct{}

func (m *${migrationName}) Up(db *gorm.DB) error {
    // マイグレーション実行（スキーマ変更）
    // 例: テーブル作成、カラム追加、インデックス追加など

    // テーブル作成例:
    // return db.AutoMigrate(&model.User{})

    // カラム追加例:
    // return db.Migrator().AddColumn(&model.User{}, "new_column")

    // インデックス作成例:
    // return db.Exec("CREATE INDEX idx_users_email ON users(email)").Error

    ${generateMigrationHint(name, modelInfo)}

    return nil
}

func (m *${migrationName}) Down(db *gorm.DB) error {
    // ロールバック処理（マイグレーションの取り消し）
    // 例: テーブル削除、カラム削除、インデックス削除など

    // テーブル削除例:
    // return db.Migrator().DropTable(&model.User{})

    // カラム削除例:
    // return db.Migrator().DropColumn(&model.User{}, "new_column")

    // インデックス削除例:
    // return db.Exec("DROP INDEX idx_users_email").Error

    return nil
}

func (m *${migrationName}) Description() string {
    return "${name}"
}
`;
}

async function analyzeModels(modelsDir: string): Promise<string> {
  try {
    const modelFiles = await readdir(modelsDir);
    const goFiles = modelFiles.filter(file => file.endsWith('.go'));

    let modelInfo = '\n    // 検出されたモデル情報:\n';

    for (const file of goFiles) {
      const content = await readFile(join(modelsDir, file), 'utf-8');
      const structMatches = content.match(/type\s+(\w+)\s+struct/g);

      if (structMatches) {
        for (const match of structMatches) {
          const modelName = match.match(/type\s+(\w+)/)?.[1];
          if (modelName) {
            modelInfo += `    // - ${modelName} (${file})\n`;
          }
        }
      }
    }

    return modelInfo;
  } catch {
    return '\n    // モデル情報の取得に失敗しました\n';
  }
}

function generateMigrationHint(name: string, modelInfo: string): string {
  const lowerName = name.toLowerCase();

  if (lowerName.includes('create') && lowerName.includes('table')) {
    const tableName = extractTableName(name);
    return `    // ${tableName}テーブル作成のヒント:
    // return db.AutoMigrate(&model.${toPascalCase(tableName)}{})${modelInfo}`;
  }

  if (lowerName.includes('add') && lowerName.includes('column')) {
    return `    // カラム追加のヒント:
    // return db.Migrator().AddColumn(&model.YourModel{}, "column_name")${modelInfo}`;
  }

  if (lowerName.includes('create') && lowerName.includes('index')) {
    return `    // インデックス作成のヒント:
    // return db.Exec("CREATE INDEX index_name ON table_name(column_name)").Error${modelInfo}`;
  }

  return `    // マイグレーション処理を実装してください${modelInfo}`;
}

async function updateMigrateFile(migrationsDir: string, fileName: string, migrationName: string) {
  const migrateFilePath = join(migrationsDir, 'migrate.go');
  const migrationClassName = toPascalCase(migrationName);

  // migrate.go ファイルが存在しない場合は作成
  let migrateContent = '';
  try {
    migrateContent = await readFile(migrateFilePath, 'utf-8');
  } catch {
    migrateContent = generateBaseMigrateFile();
  }

  // 新しいマイグレーションを追加
  const importSection = migrateContent.match(/import \(([\s\S]*?)\)/)?.[1] || '';
  const migrationsArray = migrateContent.match(/migrations := \[\]Migration\{([\s\S]*?)\}/)?.[1] || '';

  // インポートに追加は不要（同一パッケージ内のため）

  // マイグレーション配列に追加
  const newMigrationEntry = `\n        {
            ID:          "${fileName.replace('.go', '')}",
            Description: "${migrationName}",
            Migration:   &${migrationClassName}{},
        },`;

  const updatedMigrationsArray = migrationsArray.trim() + newMigrationEntry;

  const updatedContent = migrateContent.replace(
    /migrations := \[\]Migration\{[\s\S]*?\}/,
    `migrations := []Migration{${updatedMigrationsArray}\n    }`
  );

  await writeFile(migrateFilePath, updatedContent);
}

function generateBaseMigrateFile(): string {
  return `package migrations

import (
    "fmt"
    "time"
    "gorm.io/gorm"
)

type Migration struct {
    ID          string
    Description string
    Migration   MigrationInterface
}

type MigrationInterface interface {
    Up(db *gorm.DB) error
    Down(db *gorm.DB) error
    Description() string
}

type MigrationHistory struct {
    ID          uint      \`gorm:"primaryKey"\`
    MigrationID string    \`gorm:"size:255;not null;uniqueIndex"\`
    Description string    \`gorm:"size:500"\`
    ExecutedAt  time.Time \`gorm:"not null"\`
}

func RunMigrations(db *gorm.DB) error {
    // マイグレーション履歴テーブル作成
    if err := db.AutoMigrate(&MigrationHistory{}); err != nil {
        return fmt.Errorf("failed to create migration history table: %w", err)
    }

    migrations := []Migration{}

    for _, migration := range migrations {
        var count int64
        db.Model(&MigrationHistory{}).Where("migration_id = ?", migration.ID).Count(&count)

        if count == 0 {
            fmt.Printf("Running migration: %s\\n", migration.Description)

            if err := migration.Migration.Up(db); err != nil {
                return fmt.Errorf("failed to run migration %s: %w", migration.ID, err)
            }

            // 履歴記録
            history := &MigrationHistory{
                MigrationID: migration.ID,
                Description: migration.Description,
                ExecutedAt:  time.Now(),
            }
            if err := db.Create(history).Error; err != nil {
                return fmt.Errorf("failed to record migration %s: %w", migration.ID, err)
            }

            fmt.Printf("Completed migration: %s\\n", migration.Description)
        }
    }

    return nil
}

func RollbackMigrations(db *gorm.DB, steps int) error {
    var histories []MigrationHistory
    if err := db.Order("executed_at desc").Limit(steps).Find(&histories).Error; err != nil {
        return fmt.Errorf("failed to get migration history: %w", err)
    }

    migrations := []Migration{}

    for _, history := range histories {
        // マイグレーションを見つけて実行
        for _, migration := range migrations {
            if migration.ID == history.MigrationID {
                fmt.Printf("Rolling back migration: %s\\n", migration.Description)

                if err := migration.Migration.Down(db); err != nil {
                    return fmt.Errorf("failed to rollback migration %s: %w", migration.ID, err)
                }

                // 履歴から削除
                if err := db.Delete(&history).Error; err != nil {
                    return fmt.Errorf("failed to remove migration history %s: %w", migration.ID, err)
                }

                fmt.Printf("Completed rollback: %s\\n", migration.Description)
                break
            }
        }
    }

    return nil
}
`;
}

// ヘルパー関数
function toPascalCase(str: string): string {
  return str.replace(/(?:^|_)([a-z])/g, (_, char) => char.toUpperCase());
}

function extractTableName(migrationName: string): string {
  const match = migrationName.match(/create_(.+)_table/);
  return match ? match[1] : 'unknown';
}

function extractMigrationName(migrationId: string): string {
  return migrationId.substring(15); // タイムスタンプ部分を除去
}

async function getMigrationFiles(migrationsDir: string): Promise<string[]> {
  try {
    const files = await readdir(migrationsDir);
    return files
      .filter(file => file.endsWith('.go') && file !== 'migrate.go')
      .sort()
      .map(file => join(migrationsDir, file));
  } catch {
    return [];
  }
}

async function getExecutedMigrations(): Promise<string[]> {
  // 実際の実装では、データベースから実行済みマイグレーションを取得
  // ここでは簡易実装
  try {
    // データベース接続してマイグレーション履歴を取得する処理
    return [];
  } catch {
    return [];
  }
}

async function executeMigration(filePath: string): Promise<void> {
  // 実際の実装では、Goプロセスを起動してマイグレーションを実行
  console.log(`Executing migration: ${filePath}`);
}

async function rollbackMigration(filePath: string): Promise<void> {
  // 実際の実装では、Goプロセスを起動してロールバックを実行
  console.log(`Rolling back migration: ${filePath}`);
}

async function recordMigration(migrationId: string): Promise<void> {
  // マイグレーション実行履歴を記録
  console.log(`Recording migration: ${migrationId}`);
}

async function removeMigrationRecord(migrationId: string): Promise<void> {
  // マイグレーション履歴を削除
  console.log(`Removing migration record: ${migrationId}`);
}

function getMigrationExecutionTime(migrationId: string): string | null {
  // 実際の実装では、データベースから実行時刻を取得
  return null;
}

function generateMigrationSummary(results: any): string {
  const { executed, skipped, errors } = results;

  if (errors.length > 0) {
    return `❌ マイグレーション実行中にエラーが発生しました (実行: ${executed.length}, エラー: ${errors.length})`;
  }

  if (executed.length === 0 && skipped.length > 0) {
    return `✅ 全てのマイグレーション済みです (スキップ: ${skipped.length})`;
  }

  return `✅ マイグレーション完了 (実行: ${executed.length}, スキップ: ${skipped.length})`;
}
```

## Usage Example

```bash
# 新しいマイグレーション作成
/project:migrate create create_users_table

# マイグレーション実行
/project:migrate migrate

# 1つ前のマイグレーションをロールバック
/project:migrate rollback

# 3つ前のマイグレーションまでロールバック
/project:migrate rollback 3

# マイグレーション状況確認
/project:migrate status
```

## Expected Output

```markdown
## 🗄️ データベースマイグレーション結果

### ✅ マイグレーション作成完了
- **ファイル**: backend/migrations/20240324102030_create_users_table.go
- **マイグレーション名**: create_users_table
- **タイムスタンプ**: 20240324102030

### 📝 生成されたマイグレーションファイル
```go
package migrations

import (
    "gorm.io/gorm"
    "your-app/internal/model"
)

type CreateUsersTable struct{}

func (m *CreateUsersTable) Up(db *gorm.DB) error {
    // usersテーブル作成のヒント:
    // return db.AutoMigrate(&model.User{})

    // 検出されたモデル情報:
    // - User (user.go)
    // - Profile (profile.go)

    return nil
}

func (m *CreateUsersTable) Down(db *gorm.DB) error {
    // return db.Migrator().DropTable(&model.User{})
    return nil
}
```

### 🚀 次のステップ
1. 生成されたマイグレーションファイルを編集してスキーマを定義
2. `/project:migrate migrate` でマイグレーションを実行
3. 必要に応じて backend/internal/model/ でモデルを更新

### 📊 マイグレーション実行結果
- **実行されたマイグレーション**: 3個
- **スキップされたマイグレーション**: 2個
- **エラー**: 0個
- **実行時間**: 1.2秒

#### 実行詳細
✅ 20240324102030_create_users_table - ユーザーテーブル作成
✅ 20240324103045_create_profiles_table - プロフィールテーブル作成
✅ 20240324104120_add_email_index - メールインデックス追加
⏭️ 20240324101520_initial_schema - 既に実行済み
⏭️ 20240324102015_create_posts_table - 既に実行済み

### 📋 マイグレーション状況
| マイグレーションID | 名前 | 状態 | 実行日時 |
|-------------------|------|------|----------|
| 20240324101520_initial_schema | initial_schema | ✅ executed | 2024-03-24 10:15:20 |
| 20240324102015_create_posts_table | create_posts_table | ✅ executed | 2024-03-24 10:20:15 |
| 20240324102030_create_users_table | create_users_table | ⏳ pending | - |
| 20240324103045_create_profiles_table | create_profiles_table | ⏳ pending | - |

**統計**: 全4個 (実行済み: 2個, 未実行: 2個)

### 💡 推奨事項
- マイグレーション実行前にデータベースのバックアップを取得してください
- 本番環境では段階的にマイグレーションを実行することをお勧めします
- ロールバック処理も適切に実装してください

---
*この結果は `/project:migrate` コマンドにより自動生成されました*
```