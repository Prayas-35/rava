'use client'

import { useCallback, useEffect, useMemo, useState } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import { ArrowLeft, Copy, Loader2, Plus, LogOut } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { useAuth } from '@/components/auth-provider'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Input } from '@/components/ui/input'
import { Textarea } from '@/components/ui/textarea'

interface Project {
  id: string
  user_id: string
  name: string
  agent_prompt: string
  created_at: string
}

interface ApiKey {
  id: string
  project_id: string
  api_key: string
  created_at: string
  revoked: boolean
}

export default function DashboardPage() {
  const router = useRouter()
  const apiBaseUrl = process.env.NEXT_PUBLIC_ENGINE_URL ?? 'http://localhost:8080'
  const {
    token,
    isLoading: authLoading,
    isAuthenticated,
    clearToken,
  } = useAuth()

  const [projects, setProjects] = useState<Project[]>([])
  const [keysByProject, setKeysByProject] = useState<Record<string, ApiKey[]>>({})
  const [projectName, setProjectName] = useState('')
  const [agentPrompt, setAgentPrompt] = useState('')
  const [isCreateProjectModalOpen, setIsCreateProjectModalOpen] = useState(false)
  const [activeProject, setActiveProject] = useState<Project | null>(null)
  const [activeKey, setActiveKey] = useState<ApiKey | null>(null)
  const [loading, setLoading] = useState(false)
  const [creating, setCreating] = useState(false)
  const [copyText, setCopyText] = useState('Copy API key')
  const [error, setError] = useState('')

  const hasProjects = useMemo(() => projects.length > 0, [projects])

  const authedFetch = useCallback(async (path: string, init?: RequestInit) => {
    const res = await fetch(`${apiBaseUrl}${path}`, {
      ...init,
      headers: {
        ...(init?.headers ?? {}),
        Authorization: `Bearer ${token}`,
      },
    })
    return res
  }, [apiBaseUrl, token])

  const loadDashboard = useCallback(async () => {
    if (!token) return

    setLoading(true)
    setError('')

    try {
      const projectsRes = await authedFetch('/api/projects')
      if (projectsRes.status === 401) {
        clearToken()
        router.replace('/auth')
        return
      }
      if (!projectsRes.ok) {
        throw new Error(await projectsRes.text())
      }

      const projectsData = (await projectsRes.json()) as Project[]
      setProjects(projectsData)

      const keyEntries = await Promise.all(
        projectsData.map(async (project) => {
          const keysRes = await authedFetch(`/api/projects/${project.id}/keys`)
          if (!keysRes.ok) {
            return [project.id, [] as ApiKey[]] as const
          }
          const keys = (await keysRes.json()) as ApiKey[]
          return [project.id, keys] as const
        })
      )

      const map: Record<string, ApiKey[]> = {}
      for (const [projectId, keys] of keyEntries) {
        map[projectId] = keys
      }
      setKeysByProject(map)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dashboard')
    } finally {
      setLoading(false)
    }
  }, [authedFetch, clearToken, router, token])

  useEffect(() => {
    if (authLoading) return

    if (!isAuthenticated) {
      router.replace('/auth')
    }
  }, [authLoading, isAuthenticated, router])

  useEffect(() => {
    if (!authLoading && token) {
      void loadDashboard()
    }
  }, [authLoading, token, loadDashboard])

  async function handleCreateProject(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    if (!projectName.trim()) return

    setCreating(true)
    setError('')

    try {
      const createRes = await authedFetch('/api/projects', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          name: projectName,
          agent_prompt: agentPrompt,
        }),
      })

      if (!createRes.ok) {
        throw new Error(await createRes.text())
      }

      const createdProject = (await createRes.json()) as Project

      const keyRes = await authedFetch(`/api/projects/${createdProject.id}/keys`, {
        method: 'POST',
      })

      if (!keyRes.ok) {
        throw new Error(await keyRes.text())
      }

      setProjectName('')
      setAgentPrompt('You are a helpful assistant.')
      setIsCreateProjectModalOpen(false)
      await loadDashboard()
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create project')
    } finally {
      setCreating(false)
    }
  }

  async function copyApiKey(key: string) {
    await navigator.clipboard.writeText(key)
    setCopyText('Copied!')
    setTimeout(() => {
      setCopyText('Copy API key')
    }, 2000)
  }

  function openApiKeyModal(project: Project) {
    const latestKey = (keysByProject[project.id] ?? [])[0] ?? null
    setActiveProject(project)
    setActiveKey(latestKey)
  }

  function closeApiKeyModal() {
    setActiveProject(null)
    setActiveKey(null)
  }

  function openCreateProjectModal() {
    setIsCreateProjectModalOpen(true)
  }

  function closeCreateProjectModal() {
    setIsCreateProjectModalOpen(false)
  }

  function signOut() {
    clearToken()
    router.replace('/auth')
  }

  if (authLoading) {
    return (
      <div className="dark min-h-screen bg-black text-white flex items-center justify-center">
        <div className="inline-flex items-center text-zinc-300">
          <Loader2 className="w-4 h-4 mr-2 animate-spin" />
          Checking session...
        </div>
      </div>
    )
  }

  return (
    <div className="dark min-h-screen bg-black text-white">
      <header className="border-b border-zinc-900 bg-black/50 backdrop-blur-md sticky top-0 z-50">
        <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <Link href="/" className="inline-flex items-center gap-2 text-sm text-zinc-300 hover:text-white transition">
            <ArrowLeft className="w-4 h-4" />
            Home
          </Link>
          <Button variant="outline" onClick={signOut} className="border-zinc-700 text-white hover:bg-zinc-900 hover:text-white cursor-pointer">
            <LogOut className="w-4 h-4 mr-2" />
            Sign out
          </Button>
        </div>
      </header>

      <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16 space-y-8">
        <section>
          <h1 className="text-4xl sm:text-5xl font-bold mb-3">Dashboard</h1>
          <p className="text-zinc-400">Manage projects and API keys. New projects automatically receive an API key.</p>
        </section>

        <section className="rounded-xl border border-zinc-800 bg-zinc-950 p-6 flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h2 className="text-xl font-semibold mb-1">Create new project</h2>
            <p className="text-sm text-zinc-400">A project API key is generated automatically after creation.</p>
          </div>
          <Button onClick={openCreateProjectModal} className="h-11 rounded-lg bg-white text-black hover:bg-zinc-100 hover:text-black cursor-pointer">
            <Plus className="w-4 h-4 mr-1" />
            Create New Project
          </Button>
        </section>

        <section className="space-y-4">
          <h2 className="text-2xl font-semibold">Your projects</h2>

          {loading && (
            <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-4 text-zinc-300 inline-flex items-center">
              <Loader2 className="w-4 h-4 mr-2 animate-spin" />
              Loading dashboard...
            </div>
          )}

          {!loading && !hasProjects && (
            <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5 text-zinc-400">
              No projects yet. Create your first project above.
            </div>
          )}

          {!loading && hasProjects && projects.map((project) => (
            <article key={project.id} className="rounded-xl border border-zinc-800 bg-zinc-950 p-6">
              <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-2 mb-4">
                <div>
                  <h3 className="text-lg font-semibold">{project.name}</h3>
                  <p className="text-xs text-zinc-500">Project ID: {project.id}</p>
                </div>
                <p className="text-xs text-zinc-500">{new Date(project.created_at).toLocaleString()}</p>
              </div>

              <p className="text-sm text-zinc-400 mb-4">{project.agent_prompt || 'No agent prompt set.'}</p>

              <div className="flex items-center justify-between gap-3">
                <p className="text-sm text-zinc-500">
                  {(keysByProject[project.id] ?? []).length > 0
                    ? `${(keysByProject[project.id] ?? []).length} API key(s) available`
                    : 'No API keys found for this project'}
                </p>
                <Button
                  variant="outline"
                  onClick={() => openApiKeyModal(project)}
                  className="border-zinc-700 text-white hover:bg-zinc-900 hover:text-white cursor-pointer"
                >
                  Show API Key
                </Button>
              </div>
            </article>
          ))}
        </section>

        {error && (
          <section className="rounded-lg border border-red-900/60 bg-red-950/30 p-3 text-sm text-red-200">
            {error}
          </section>
        )}

        <Dialog open={isCreateProjectModalOpen} onOpenChange={setIsCreateProjectModalOpen}>
          <DialogContent className="max-w-3xl border-zinc-800 bg-zinc-950 text-white" showCloseButton={false}>
            <DialogHeader>
              <DialogTitle className='text-xl'>Create New Project</DialogTitle>
              <DialogDescription>
                Fill in the details below to create a project.
              </DialogDescription>
            </DialogHeader>

            <form onSubmit={handleCreateProject} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="project-name" className="text-sm font-medium text-zinc-200">Project Name</Label>
                <Input
                  id="project-name"
                  type="text"
                  value={projectName}
                  onChange={(e) => setProjectName(e.target.value)}
                  required
                  placeholder="e.g., Customer Support Assistant"
                  className="w-full h-11 rounded-lg bg-zinc-900 border border-zinc-800 px-3 text-sm outline-none focus:border-zinc-600"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="agent-prompt" className="text-sm font-medium text-zinc-200">Agent Prompt</Label>
                <Textarea
                  id="agent-prompt"
                  value={agentPrompt}
                  onChange={(e) => setAgentPrompt(e.target.value)}
                  rows={4}
                  placeholder="e.g., You are an expert support assistant. Answer using concise, accurate steps."
                  className="w-full rounded-lg bg-zinc-900 border border-zinc-800 px-3 py-2 text-sm outline-none focus:border-zinc-600"
                />
              </div>

              <DialogFooter>
                <Button
                  type="button"
                  variant="outline"
                  onClick={closeCreateProjectModal}
                  className="border-zinc-700 text-white hover:bg-zinc-900 hover:text-white cursor-pointer"
                >
                  Cancel
                </Button>
                <Button type="submit" size="sm" disabled={creating} className="rounded-lg bg-white text-black hover:bg-zinc-100 hover:text-black cursor-pointer">
                  {creating ? <Loader2 className="w-4 h-4 mr-2 animate-spin" /> : <Plus className="w-4 h-4 mr-1" />}
                  Done
                </Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>

        <Dialog
          open={!!activeProject}
          onOpenChange={(open) => {
            if (!open) closeApiKeyModal()
          }}
        >
          <DialogContent className="max-w-3xl border-zinc-800 bg-zinc-950 text-white" showCloseButton={false}>
            <DialogHeader>
              <DialogTitle>API key</DialogTitle>
              <DialogDescription>
                Project: {activeProject?.name ?? 'Unknown project'}
              </DialogDescription>
            </DialogHeader>

            {!activeKey && (
              <p className="text-sm text-zinc-400">No API key found for this project yet.</p>
            )}

            {activeKey && (
              <div className="space-y-3">
                <p className="text-xs text-zinc-500">Created {new Date(activeKey.created_at).toLocaleString()}</p>
                <textarea
                  readOnly
                  value={activeKey.api_key}
                  rows={3}
                  className="w-full rounded bg-zinc-900 border border-zinc-800 px-3 py-2 text-sm text-emerald-300"
                />
                <Button
                  variant="outline"
                  onClick={() => copyApiKey(activeKey.api_key)}
                  className="border-zinc-700 text-white hover:bg-zinc-900 hover:text-white cursor-pointer"
                >
                  <Copy className="w-4 h-4 mr-2" />
                    {copyText}
                </Button>
              </div>
            )}

            <DialogFooter>
              <Button
                type="button"
                variant="outline"
                onClick={closeApiKeyModal}
                className="border-zinc-700 text-white hover:bg-zinc-900 hover:text-white cursor-pointer"
              >
                Close
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>
      </main>
    </div>
  )
}
