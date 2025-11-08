/* eslint-disable react-refresh/only-export-components */
import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'
import { api } from '@/lib/api'
import type { User } from '@/types'

interface AuthContextType {
  user: User | null
  isLoading: boolean
  isAuthenticated: boolean
  login: (email: string, password: string) => Promise<void>
  signup: (email: string, password: string, nickname: string) => Promise<void>
  logout: () => Promise<void>
  refetch: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  const fetchUser = async () => {
    setIsLoading(true)
    const [data, error] = await api.account.me()

    if (error) {
      setUser(null)
    } else if (data) {
      setUser({
        id: data.id.toString(),
        email: data.email,
        nickname: data.nickname
      })
    }

    setIsLoading(false)
  }

  useEffect(() => {
    fetchUser()
  }, [])

  const login = async (email: string, password: string) => {
    const [, error] = await api.account.login({ email, password })

    if (error) {
      throw new Error(error.message || 'Login failed')
    }

    await fetchUser()
  }

  const signup = async (email: string, password: string, nickname: string) => {
    const [, error] = await api.account.signup({ email, password, nickname })

    if (error) {
      throw new Error(error.message || 'Signup failed')
    }
  }

  const logout = async () => {
    await api.account.logout()
    setUser(null)
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        signup,
        logout,
        refetch: fetchUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
