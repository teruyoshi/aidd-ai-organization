# フロントエンド技術コンテキスト

## 技術スタック詳細

### 核心技術
- **React 18.2+**: Concurrent Features・Suspense・Error Boundaries活用
- **TypeScript 5.0+**: 厳格型チェック・最新型機能活用
- **Vite 5.0+**: 高速ビルド・HMR・ESM最適化

### 状態管理戦略
```typescript
// サーバー状態: TanStack Query (React Query)
const { data, isLoading, error } = useQuery({
  queryKey: ['users', { page, filter }],
  queryFn: ({ queryKey }) => fetchUsers(queryKey[1]),
  staleTime: 5 * 60 * 1000,
});

// クライアント状態: Context + useReducer
interface AppState {
  theme: 'light' | 'dark';
  sidebar: { isOpen: boolean };
  notifications: Notification[];
}

// フォーム状態: react-hook-form
const { register, handleSubmit, formState: { errors } } = useForm<UserForm>({
  resolver: zodResolver(userSchema),
  mode: 'onChange',
});
```

## アーキテクチャパターン

### 1. Feature-Based構造
```
src/
├── features/
│   ├── auth/
│   │   ├── components/
│   │   ├── hooks/
│   │   ├── services/
│   │   └── types/
│   └── users/
│       ├── components/
│       ├── hooks/
│       ├── services/
│       └── types/
├── shared/
│   ├── components/ui/
│   ├── hooks/
│   ├── services/
│   └── types/
└── app/
    ├── providers/
    ├── router/
    └── store/
```

### 2. コンポーネント設計原則
```typescript
// Compound Component Pattern
export const UserCard = {
  Root: UserCardRoot,
  Header: UserCardHeader,
  Content: UserCardContent,
  Actions: UserCardActions,
};

// 使用例
<UserCard.Root>
  <UserCard.Header>
    <Avatar src={user.avatar} />
    <UserCard.Title>{user.name}</UserCard.Title>
  </UserCard.Header>
  <UserCard.Content>
    <p>{user.bio}</p>
  </UserCard.Content>
  <UserCard.Actions>
    <Button onClick={onEdit}>Edit</Button>
  </UserCard.Actions>
</UserCard.Root>
```

## パフォーマンス最適化

### 1. バンドル最適化
```typescript
// vite.config.ts
export default defineConfig({
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          'vendor': ['react', 'react-dom'],
          'ui': ['@radix-ui/react-dialog', '@radix-ui/react-dropdown-menu'],
          'utils': ['lodash', 'date-fns'],
        },
      },
    },
  },
  optimizeDeps: {
    include: ['react', 'react-dom', '@tanstack/react-query'],
  },
});
```

### 2. レンダリング最適化
```typescript
// メモ化戦略
const UserList = memo(({ users }: UserListProps) => {
  const sortedUsers = useMemo(
    () => users.sort((a, b) => a.name.localeCompare(b.name)),
    [users]
  );

  const handleUserClick = useCallback(
    (userId: string) => {
      onUserSelect(userId);
    },
    [onUserSelect]
  );

  return (
    <div>
      {sortedUsers.map(user => (
        <UserCard
          key={user.id}
          user={user}
          onClick={handleUserClick}
        />
      ))}
    </div>
  );
});
```

### 3. 遅延読み込み
```typescript
// Route-based code splitting
const Dashboard = lazy(() => import('../pages/Dashboard'));
const UserManagement = lazy(() => import('../pages/UserManagement'));

// Component-based lazy loading
const HeavyChart = lazy(() => import('../components/HeavyChart'));

const App = () => (
  <Router>
    <Suspense fallback={<LoadingSpinner />}>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/users" element={<UserManagement />} />
      </Routes>
    </Suspense>
  </Router>
);
```

## データフェッチング戦略

### 1. TanStack Query設定
```typescript
// queryClient.ts
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 5, // 5分
      cacheTime: 1000 * 60 * 30, // 30分
      refetchOnWindowFocus: false,
      retry: (failureCount, error) => {
        if (error.status === 404) return false;
        return failureCount < 3;
      },
    },
    mutations: {
      onError: (error) => {
        toast.error(error.message);
      },
    },
  },
});
```

### 2. カスタムフック
```typescript
// useUsers.ts
export const useUsers = (params: UsersParams) => {
  return useQuery({
    queryKey: ['users', params],
    queryFn: async ({ signal }) => {
      const response = await api.get('/users', {
        params,
        signal,
      });
      return response.data;
    },
    enabled: !!params,
    placeholderData: keepPreviousData,
  });
};

// useCreateUser.ts
export const useCreateUser = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async (userData: CreateUserData) => {
      const response = await api.post('/users', userData);
      return response.data;
    },
    onSuccess: (newUser) => {
      // キャッシュ更新
      queryClient.setQueryData(['users'], (oldData: User[]) => [
        ...(oldData || []),
        newUser,
      ]);

      // 関連クエリの無効化
      queryClient.invalidateQueries(['users']);

      toast.success('User created successfully');
    },
    onError: (error) => {
      toast.error(error.message);
    },
  });
};
```

## UI/UX設計システム

### 1. デザイントークン
```typescript
// tokens.ts
export const tokens = {
  colors: {
    primary: {
      50: '#eff6ff',
      500: '#3b82f6',
      900: '#1e3a8a',
    },
    gray: {
      50: '#f9fafb',
      500: '#6b7280',
      900: '#111827',
    },
  },
  spacing: {
    xs: '0.25rem',
    sm: '0.5rem',
    md: '1rem',
    lg: '1.5rem',
    xl: '2rem',
  },
  typography: {
    fontFamily: {
      sans: ['Inter', 'sans-serif'],
      mono: ['Fira Code', 'monospace'],
    },
    fontSize: {
      xs: ['0.75rem', { lineHeight: '1rem' }],
      sm: ['0.875rem', { lineHeight: '1.25rem' }],
      base: ['1rem', { lineHeight: '1.5rem' }],
      lg: ['1.125rem', { lineHeight: '1.75rem' }],
    },
  },
} as const;
```

### 2. コンポーネントライブラリ
```typescript
// Button.tsx
interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost';
  size?: 'sm' | 'md' | 'lg';
  isLoading?: boolean;
  disabled?: boolean;
  children: React.ReactNode;
  onClick?: () => void;
}

export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  isLoading = false,
  disabled = false,
  children,
  onClick,
  ...props
}) => {
  const classes = cn(
    'inline-flex items-center justify-center rounded-md font-medium transition-colors',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    'disabled:opacity-50 disabled:pointer-events-none',
    {
      // Variant styles
      'bg-blue-600 text-white hover:bg-blue-700 focus:ring-blue-500': variant === 'primary',
      'bg-gray-200 text-gray-900 hover:bg-gray-300 focus:ring-gray-500': variant === 'secondary',
      // Size styles
      'h-8 px-3 text-sm': size === 'sm',
      'h-10 px-4': size === 'md',
      'h-12 px-6 text-lg': size === 'lg',
    }
  );

  return (
    <button
      className={classes}
      disabled={disabled || isLoading}
      onClick={onClick}
      {...props}
    >
      {isLoading && <Spinner className="mr-2 h-4 w-4" />}
      {children}
    </button>
  );
};
```

## フォーム処理パターン

### 1. react-hook-form + Zod
```typescript
// schemas/user.ts
export const userSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  confirmPassword: z.string(),
}).refine((data) => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ["confirmPassword"],
});

// UserForm.tsx
export const UserForm: React.FC<UserFormProps> = ({ onSubmit }) => {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    reset,
  } = useForm<UserFormData>({
    resolver: zodResolver(userSchema),
    mode: 'onChange',
  });

  const submitHandler = async (data: UserFormData) => {
    try {
      await onSubmit(data);
      reset();
      toast.success('User created successfully');
    } catch (error) {
      toast.error('Failed to create user');
    }
  };

  return (
    <form onSubmit={handleSubmit(submitHandler)}>
      <Input
        {...register('name')}
        label="Name"
        error={errors.name?.message}
      />
      <Input
        {...register('email')}
        type="email"
        label="Email"
        error={errors.email?.message}
      />
      <Button type="submit" isLoading={isSubmitting}>
        Create User
      </Button>
    </form>
  );
};
```

## エラーハンドリング

### 1. Error Boundaries
```typescript
// ErrorBoundary.tsx
interface ErrorBoundaryState {
  hasError: boolean;
  error?: Error;
}

export class ErrorBoundary extends Component<
  PropsWithChildren<{}>,
  ErrorBoundaryState
> {
  constructor(props: PropsWithChildren<{}>) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): ErrorBoundaryState {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error Boundary caught an error:', error, errorInfo);
    // エラー報告サービスに送信
    errorReportingService.captureException(error, {
      extra: errorInfo,
    });
  }

  render() {
    if (this.state.hasError) {
      return <ErrorFallback error={this.state.error} />;
    }

    return this.props.children;
  }
}
```

### 2. API エラーハンドリング
```typescript
// api/client.ts
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
});

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // 認証エラー: ログアウト処理
      authStore.logout();
      window.location.href = '/login';
    } else if (error.response?.status >= 500) {
      // サーバーエラー: エラートラッキング
      errorReportingService.captureException(error);
    }

    return Promise.reject(error);
  }
);
```

## テスト戦略

### 1. MSW (Mock Service Worker) 設定
```typescript
// __mocks__/handlers.ts
export const handlers = [
  rest.get('/api/users', (req, res, ctx) => {
    return res(
      ctx.status(200),
      ctx.json([
        { id: '1', name: 'John Doe', email: 'john@example.com' },
        { id: '2', name: 'Jane Smith', email: 'jane@example.com' },
      ])
    );
  }),

  rest.post('/api/users', (req, res, ctx) => {
    return res(
      ctx.status(201),
      ctx.json({ id: '3', name: 'New User', email: 'new@example.com' })
    );
  }),
];

// __mocks__/server.ts
export const server = setupServer(...handlers);
```

### 2. カスタムテストユーティリティ
```typescript
// test-utils.tsx
const AllTheProviders: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: { retry: false },
      mutations: { retry: false },
    },
  });

  return (
    <QueryClientProvider client={queryClient}>
      <BrowserRouter>
        <ThemeProvider>
          {children}
        </ThemeProvider>
      </BrowserRouter>
    </QueryClientProvider>
  );
};

export const renderWithProviders = (
  ui: ReactElement,
  options?: Omit<RenderOptions, 'wrapper'>
) => render(ui, { wrapper: AllTheProviders, ...options });
```

## アクセシビリティ

### 1. セマンティックHTML
```typescript
// 適切なHTMLセマンティクス使用
const UserProfile = () => (
  <article>
    <header>
      <h1>{user.name}</h1>
      <p>Member since {formatDate(user.createdAt)}</p>
    </header>

    <section>
      <h2>Contact Information</h2>
      <address>
        <p>Email: <a href={`mailto:${user.email}`}>{user.email}</a></p>
      </address>
    </section>
  </article>
);
```

### 2. ARIA属性・キーボードナビゲーション
```typescript
// Modal.tsx
export const Modal: React.FC<ModalProps> = ({ isOpen, onClose, children }) => {
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape') onClose();
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
      return () => document.removeEventListener('keydown', handleEscape);
    }
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-labelledby="modal-title"
      className="fixed inset-0 z-50"
    >
      <div className="fixed inset-0 bg-black bg-opacity-50" onClick={onClose} />
      <div className="relative max-w-md mx-auto mt-20 p-6 bg-white rounded-lg">
        {children}
      </div>
    </div>
  );
};
```

## 国際化 (i18n)

```typescript
// i18n/resources.ts
export const resources = {
  en: {
    common: {
      save: 'Save',
      cancel: 'Cancel',
      loading: 'Loading...',
    },
    users: {
      title: 'User Management',
      createUser: 'Create New User',
      editUser: 'Edit User',
    },
  },
  ja: {
    common: {
      save: '保存',
      cancel: 'キャンセル',
      loading: '読み込み中...',
    },
    users: {
      title: 'ユーザー管理',
      createUser: '新規ユーザー作成',
      editUser: 'ユーザー編集',
    },
  },
} as const;

// 使用例
const UserPage = () => {
  const { t } = useTranslation();

  return (
    <div>
      <h1>{t('users.title')}</h1>
      <Button>{t('users.createUser')}</Button>
    </div>
  );
};
```

## パフォーマンス監視

```typescript
// performance/monitoring.ts
export const performanceMonitor = {
  // Web Vitals測定
  measureWebVitals: () => {
    getCLS(console.log);
    getFID(console.log);
    getFCP(console.log);
    getLCP(console.log);
    getTTFB(console.log);
  },

  // カスタムメトリクス
  measurePageLoad: (pageName: string) => {
    performance.mark(`${pageName}-start`);

    return () => {
      performance.mark(`${pageName}-end`);
      performance.measure(`${pageName}-duration`, `${pageName}-start`, `${pageName}-end`);

      const measure = performance.getEntriesByName(`${pageName}-duration`)[0];
      console.log(`${pageName} loaded in ${measure.duration}ms`);
    };
  },
};
```