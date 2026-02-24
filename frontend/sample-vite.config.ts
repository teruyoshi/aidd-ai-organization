import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react-swc'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],

  // 開発サーバー設定
  server: {
    host: '0.0.0.0',
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://backend:8080',
        changeOrigin: true,
        secure: false
      },
      '/ws': {
        target: 'ws://backend:8080',
        ws: true,
        changeOrigin: true
      }
    }
  },

  // プレビューサーバー設定
  preview: {
    host: '0.0.0.0',
    port: 3000
  },

  // パス解決設定
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
      '@/components': resolve(__dirname, './src/components'),
      '@/features': resolve(__dirname, './src/features'),
      '@/hooks': resolve(__dirname, './src/hooks'),
      '@/services': resolve(__dirname, './src/services'),
      '@/types': resolve(__dirname, './src/types'),
      '@/utils': resolve(__dirname, './src/utils'),
      '@/styles': resolve(__dirname, './src/styles')
    }
  },

  // ビルド設定
  build: {
    outDir: 'dist',
    sourcemap: true,
    target: 'es2020',
    rollupOptions: {
      output: {
        manualChunks: {
          // React関連
          'vendor-react': ['react', 'react-dom', 'react-router-dom'],
          // UI関連
          'vendor-ui': [
            '@headlessui/react',
            '@heroicons/react',
            'framer-motion',
            'lucide-react'
          ],
          // 状態管理・データフェッチ
          'vendor-state': [
            '@tanstack/react-query',
            'zustand',
            'axios'
          ],
          // フォーム・バリデーション
          'vendor-forms': [
            'react-hook-form',
            '@hookform/resolvers',
            'zod'
          ],
          // ユーティリティ
          'vendor-utils': [
            'date-fns',
            'clsx',
            'tailwind-merge'
          ]
        }
      }
    },
    // チャンクサイズ警告の閾値
    chunkSizeWarningLimit: 500
  },

  // 環境変数設定
  envPrefix: 'VITE_',

  // 最適化設定
  optimizeDeps: {
    include: [
      'react',
      'react-dom',
      'react-router-dom',
      '@tanstack/react-query',
      'axios',
      'date-fns',
      'zod'
    ]
  },

  // テスト設定
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
    css: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      exclude: [
        'node_modules/',
        'src/test/',
        '**/*.d.ts',
        '**/*.config.ts',
        '**/*.config.js'
      ],
      thresholds: {
        global: {
          branches: 80,
          functions: 80,
          lines: 80,
          statements: 80
        }
      }
    }
  },

  // CSS設定
  css: {
    modules: {
      localsConvention: 'camelCaseOnly'
    },
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/styles/variables.scss";`
      }
    }
  },

  // PWA設定（将来的な拡張用）
  define: {
    __APP_VERSION__: JSON.stringify(process.env.npm_package_version),
    __BUILD_TIME__: JSON.stringify(new Date().toISOString())
  }
})