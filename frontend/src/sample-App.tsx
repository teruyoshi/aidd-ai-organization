import React, { Suspense } from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { Helmet } from 'react-helmet-async'

import { Layout } from '@/components/layout/sample-Layout'
import { LoadingSpinner } from '@/components/ui/sample-LoadingSpinner'
import { useAuth } from '@/features/auth/hooks/sample-useAuth'

// 遅延ロード（Code Splitting）
const LoginPage = React.lazy(() => import('@/features/auth/components/sample-LoginPage'))
const RegisterPage = React.lazy(() => import('@/features/auth/components/sample-RegisterPage'))
const TodosPage = React.lazy(() => import('@/features/todos/components/sample-TodosPage'))
const ProfilePage = React.lazy(() => import('@/features/auth/components/sample-ProfilePage'))
const NotFoundPage = React.lazy(() => import('@/components/pages/sample-NotFoundPage'))

// プライベートルート保護コンポーネント
const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="large" />
      </div>
    )
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

// パブリックルート（認証済みユーザーはリダイレクト）
const PublicRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <LoadingSpinner size="large" />
      </div>
    )
  }

  if (user) {
    return <Navigate to="/todos" replace />
  }

  return <>{children}</>
}

// 遅延ロード時の読み込み画面
const SuspenseLoader: React.FC = () => (
  <div className="min-h-screen flex items-center justify-center bg-gray-50">
    <div className="text-center">
      <LoadingSpinner size="large" />
      <p className="mt-4 text-gray-600">読み込み中...</p>
    </div>
  </div>
)

const App: React.FC = () => {
  return (
    <>
      <Helmet>
        <title>AIDD TODO Sample</title>
        <meta name="description" content="AI-Driven Development による TODO アプリケーション サンプル実装" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <meta name="theme-color" content="#3B82F6" />

        {/* Open Graph メタタグ */}
        <meta property="og:title" content="AIDD TODO Sample" />
        <meta property="og:description" content="AI-Driven Development による TODO アプリケーション サンプル実装" />
        <meta property="og:type" content="website" />

        {/* Twitter Card メタタグ */}
        <meta name="twitter:card" content="summary" />
        <meta name="twitter:title" content="AIDD TODO Sample" />
        <meta name="twitter:description" content="AI-Driven Development による TODO アプリケーション サンプル実装" />

        {/* PWA関連 */}
        <link rel="manifest" href="/manifest.json" />
        <link rel="apple-touch-icon" href="/icons/icon-192.png" />

        {/* フォント読み込み */}
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="" />
        <link
          href="https://fonts.googleapis.com/css2?family=Inter:wght@300;400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </Helmet>

      <div className="min-h-screen bg-gray-50">
        <Routes>
          {/* パブリックルート（認証不要） */}
          <Route
            path="/login"
            element={
              <PublicRoute>
                <Suspense fallback={<SuspenseLoader />}>
                  <LoginPage />
                </Suspense>
              </PublicRoute>
            }
          />

          <Route
            path="/register"
            element={
              <PublicRoute>
                <Suspense fallback={<SuspenseLoader />}>
                  <RegisterPage />
                </Suspense>
              </PublicRoute>
            }
          />

          {/* プライベートルート（認証必要） */}
          <Route
            path="/"
            element={
              <PrivateRoute>
                <Layout />
              </PrivateRoute>
            }
          >
            {/* デフォルトはTODOページにリダイレクト */}
            <Route index element={<Navigate to="/todos" replace />} />

            <Route
              path="todos"
              element={
                <Suspense fallback={<LoadingSpinner size="large" />}>
                  <TodosPage />
                </Suspense>
              }
            />

            <Route
              path="profile"
              element={
                <Suspense fallback={<LoadingSpinner size="large" />}>
                  <ProfilePage />
                </Suspense>
              }
            />
          </Route>

          {/* 404ページ */}
          <Route
            path="*"
            element={
              <Suspense fallback={<SuspenseLoader />}>
                <NotFoundPage />
              </Suspense>
            }
          />
        </Routes>
      </div>
    </>
  )
}

export default App