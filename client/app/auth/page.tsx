'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { ArrowLeft, Loader2 } from 'lucide-react'
import { Button } from '@/components/ui/button'

type AuthMode = 'signup' | 'signin'

interface AuthResponse {
  token: string
}

export default function AuthPage() {
  const router = useRouter()
  const apiBaseUrl = process.env.NEXT_PUBLIC_ENGINE_URL ?? 'http://localhost:8080'

  const [mode, setMode] = useState<AuthMode>('signup')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    const token = window.localStorage.getItem('rava_jwt')
    if (token) {
      router.replace('/dashboard')
    }
  }, [router])

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const endpoint = mode === 'signup' ? '/api/auth/signup' : '/api/auth/signin'
      const payload = mode === 'signup'
        ? { email, password, name }
        : { email, password }

      const res = await fetch(`${apiBaseUrl}${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })

      if (!res.ok) {
        const text = await res.text()
        throw new Error(text || 'Authentication failed')
      }

      const data = (await res.json()) as AuthResponse
      window.localStorage.setItem('rava_jwt', data.token)
      router.replace('/dashboard')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Authentication failed')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="dark min-h-screen bg-black text-white">
      <header className="border-b border-zinc-900 bg-black/50 backdrop-blur-md sticky top-0 z-50">
        <div className="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <Link href="/" className="inline-flex items-center gap-2 text-sm text-zinc-300 hover:text-white transition">
            <ArrowLeft className="w-4 h-4" />
            Back to home
          </Link>
        </div>
      </header>

      <main className="max-w-md mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <h1 className="text-4xl font-bold mb-3">{mode === 'signup' ? 'Create your account' : 'Welcome back'}</h1>
        <p className="text-zinc-400 mb-8">Sign in to access your projects and API keys.</p>

        <div className="inline-flex rounded-lg border border-zinc-800 p-1 bg-zinc-900 mb-6">
          <button
            type="button"
            onClick={() => setMode('signup')}
            className={`px-4 py-2 rounded-md text-sm transition ${mode === 'signup' ? 'bg-white text-black' : 'text-zinc-300 hover:text-white'}`}
          >
            Sign up
          </button>
          <button
            type="button"
            onClick={() => setMode('signin')}
            className={`px-4 py-2 rounded-md text-sm transition ${mode === 'signin' ? 'bg-white text-black' : 'text-zinc-300 hover:text-white'}`}
          >
            Sign in
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-3">
          <input
            type="email"
            required
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full h-11 rounded-lg bg-zinc-900 border border-zinc-800 px-3 text-sm outline-none focus:border-zinc-600"
          />

          {mode === 'signup' && (
            <input
              type="text"
              placeholder="Name (optional)"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="w-full h-11 rounded-lg bg-zinc-900 border border-zinc-800 px-3 text-sm outline-none focus:border-zinc-600"
            />
          )}

          <input
            type="password"
            required
            placeholder="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full h-11 rounded-lg bg-zinc-900 border border-zinc-800 px-3 text-sm outline-none focus:border-zinc-600"
          />

          <Button type="submit" disabled={loading} className="w-full h-11 rounded-lg bg-white text-black hover:bg-zinc-100">
            {loading ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : null}
            {mode === 'signup' ? 'Create account' : 'Sign in'}
          </Button>
        </form>

        {error && (
          <div className="mt-4 rounded-lg border border-red-900/60 bg-red-950/30 p-3 text-sm text-red-200">
            {error}
          </div>
        )}
      </main>
    </div>
  )
}
