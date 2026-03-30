'use client'

import { createContext, useContext, useMemo, useSyncExternalStore } from 'react'

const AUTH_TOKEN_KEY = 'rava_jwt'
const AUTH_TOKEN_CHANGED_EVENT = 'rava-auth-token-change'

interface AuthContextValue {
  token: string | null
  isLoading: boolean
  isAuthenticated: boolean
  setToken: (nextToken: string) => void
  clearToken: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

function subscribe(onStoreChange: () => void) {
  if (typeof window === 'undefined') {
    return () => {}
  }

  const handleStorage = () => {
    onStoreChange()
  }

  window.addEventListener('storage', handleStorage)
  window.addEventListener(AUTH_TOKEN_CHANGED_EVENT, handleStorage)

  return () => {
    window.removeEventListener('storage', handleStorage)
    window.removeEventListener(AUTH_TOKEN_CHANGED_EVENT, handleStorage)
  }
}

function getSnapshot() {
  if (typeof window === 'undefined') {
    return null
  }
  return window.localStorage.getItem(AUTH_TOKEN_KEY)
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const token = useSyncExternalStore(subscribe, getSnapshot, () => null)

  const setToken = (nextToken: string) => {
    window.localStorage.setItem(AUTH_TOKEN_KEY, nextToken)
    window.dispatchEvent(new Event(AUTH_TOKEN_CHANGED_EVENT))
  }

  const clearToken = () => {
    window.localStorage.removeItem(AUTH_TOKEN_KEY)
    window.dispatchEvent(new Event(AUTH_TOKEN_CHANGED_EVENT))
  }

  const value = useMemo<AuthContextValue>(() => ({
    token,
    isLoading: false,
    isAuthenticated: !!token,
    setToken,
    clearToken,
  }), [token])

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
