import { Routes, Route, Navigate } from 'react-router-dom'
import { useAuth } from '@/features/auth/hooks/sample-useAuth'

import LoginPage from '@/features/auth/components/sample-LoginPage'
import RegisterPage from '@/features/auth/components/sample-RegisterPage'
import TodosPage from '@/features/todos/components/sample-TodosPage'

const PrivateRoute: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const { user, isLoading } = useAuth()

  if (isLoading) {
    return <div className="min-h-screen flex items-center justify-center">読み込み中...</div>
  }

  if (!user) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

const App: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-50">
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />
        <Route
          path="/todos"
          element={
            <PrivateRoute>
              <TodosPage />
            </PrivateRoute>
          }
        />
        <Route path="*" element={<Navigate to="/todos" replace />} />
      </Routes>
    </div>
  )
}

export default App
