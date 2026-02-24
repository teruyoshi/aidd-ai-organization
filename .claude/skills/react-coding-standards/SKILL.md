---
name: react-coding-standards
description: Reactフロントエンドのコーディング規約とアーキテクチャパターン。TSX・React・TypeScript・Tailwind・TanStack Queryの実装・レビュー時に自動適用。
user-invocable: false
---

# React コーディング規約（AIDD プロジェクト）

## Feature-Based ディレクトリ構造

```
frontend/src/
├── features/              # 機能ごとのモジュール（最重要）
│   └── [feature-name]/
│       ├── components/    # この機能専用コンポーネント
│       ├── hooks/         # カスタムフック（useQuery等）
│       ├── services/      # APIクライアント関数
│       └── types/         # TypeScript型定義
├── components/            # 共有コンポーネント
│   ├── ui/               # 汎用UIパーツ（Button, Input等）
│   ├── forms/            # フォーム部品
│   └── layout/           # レイアウト（Header, Sidebar等）
├── hooks/                 # グローバルカスタムフック
├── services/              # グローバルAPIクライアント
├── types/                 # グローバル型定義
└── utils/                 # ユーティリティ関数
```

## 状態管理の使い分け

- **サーバー状態**: TanStack Query（useQuery / useMutation）
- **クライアントUI状態**: Context API + useReducer
- **フォーム状態**: react-hook-form + zod バリデーション
- **グローバル状態**: 原則 Context のみ（Redux 不使用）

## TanStack Query 規約

```typescript
// クエリキーは配列 + オブジェクト形式
const { data, isLoading } = useQuery({
  queryKey: ['todos', { status, userId }],
  queryFn: () => fetchTodos({ status, userId }),
  staleTime: 5 * 60 * 1000,
});

// Mutation はキャッシュ無効化を忘れずに
const mutation = useMutation({
  mutationFn: createTodo,
  onSuccess: () => queryClient.invalidateQueries({ queryKey: ['todos'] }),
});
```

## TypeScript 規約

- `strict: true` 必須、`any` 禁止
- API レスポンス型は `types/api.ts` に集約
- Props 型は `interface` で定義、コンポーネント直上に記述

## コンポーネント設計

- 1ファイル1コンポーネント原則
- Props は最小限（必要なものだけ受け取る）
- ロジックは Custom Hook に切り出す
- Tailwind CSS のクラスは `cn()` ユーティリティで結合

## テスト規約

- ユニット/統合: Vitest + React Testing Library
- E2E: Playwright（`frontend/e2e/`）
- `userEvent` を使用、`fireEvent` は非推奨

詳細は `docs/context/frontend-context.md` を参照。
