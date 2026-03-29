import Link from 'next/link'
import { ArrowLeft, BookOpen, Braces, Terminal } from 'lucide-react'
import { Button } from '@/components/ui/button'

const installSnippet = `pnpm add rava`

const initSnippet = `import { RavaClient } from 'rava'

RavaClient.initialize({
  apiKey: process.env.RAVA_API_KEY!,
  // Optional. Defaults to https://api.rava.dev
  // Use your local engine during development if needed.
  baseUrl: process.env.RAVA_BASE_URL ?? 'http://localhost:8080',
})`

const ingestSnippet = `const client = RavaClient.getInstance()

await client.ingest({
  name: 'kb-intro',
  content: 'Rava uses your data to answer questions.',
  metadata: { type: 'text' },
})

await client.ingest({
  name: 'kb-file',
  filePath: './data.txt',
  metadata: { type: 'file' },
})`

const querySnippet = `const client = RavaClient.getInstance()

const response = await client.query({
  question: 'What did I ingest?',
  history: [
    { role: 'user', content: 'Answer in one sentence.' },
  ],
  top_k: 5,
})

console.log(response.answer)`

const getApiKeyFlowSnippet = `// 1) Sign up (or sign in)
const auth = await fetch('http://localhost:8080/api/auth/signup', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        email: 'dev@example.com',
        password: 'strong-password',
        name: 'Dev User',
    }),
}).then((r) => r.json())

const token = auth.token

// 2) Create project (JWT determines user_id on backend)
const project = await fetch('http://localhost:8080/api/projects', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        Authorization: \`Bearer \${token}\`,
    },
    body: JSON.stringify({
        name: 'my-first-project',
        agent_prompt: 'You are a helpful assistant.',
    }),
}).then((r) => r.json())

// 3) Generate API key for the project
const key = await fetch(\`http://localhost:8080/api/projects/\${project.id}/keys\`, {
    method: 'POST',
    headers: {
        Authorization: \`Bearer \${token}\`,
    },
}).then((r) => r.json())

console.log(key.api_key)`

export default function DocsPage() {
    return (
        <div className="dark min-h-screen bg-black text-white">
            <header className="border-b border-zinc-900 bg-black/50 backdrop-blur-md sticky top-0 z-50">
                <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
                    <Link href="/" className="inline-flex items-center gap-2 text-sm text-zinc-300 hover:text-white transition">
                        <ArrowLeft className="w-4 h-4" />
                        Back to home
                    </Link>
                    <div className="inline-flex items-center gap-2 text-sm text-zinc-400">
                        <BookOpen className="w-4 h-4" />
                        SDK Documentation
                    </div>
                </div>
            </header>

            <main className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-16 space-y-12">
                <section>
                    <h1 className="text-4xl sm:text-5xl font-bold mb-4">Rava SDK docs</h1>
                    <p className="text-zinc-400 text-lg leading-relaxed max-w-3xl">
                        This page reflects the current SDK implementation in this repository. It documents the exported
                        client, method signatures, default values, and required fields as implemented in the
                        <span className="text-zinc-200"> rava </span>
                        package.
                    </p>
                </section>

                <section className="space-y-4">
                    <h2 className="text-2xl sm:text-3xl font-semibold">1) Install</h2>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5">
                        <div className="flex items-center gap-2 text-zinc-400 text-sm mb-3">
                            <Terminal className="w-4 h-4" />
                            Terminal
                        </div>
                        <pre className="text-sm text-zinc-200 overflow-auto">{installSnippet}</pre>
                    </div>
                </section>

                <section className="space-y-4">
                    <h2 className="text-2xl sm:text-3xl font-semibold">2) Initialize once</h2>
                    <p className="text-zinc-400 leading-relaxed">
                        The SDK uses a singleton. Call <span className="text-zinc-200">RavaClient.initialize()</span> once at app startup,
                        then use <span className="text-zinc-200">RavaClient.getInstance()</span> wherever you need it.
                    </p>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5">
                        <div className="flex items-center gap-2 text-zinc-400 text-sm mb-3">
                            <Braces className="w-4 h-4" />
                            TypeScript
                        </div>
                        <pre className="text-sm text-zinc-200 overflow-auto">{initSnippet}</pre>
                    </div>
                    <div className="rounded-lg border border-amber-900/50 bg-amber-950/20 p-4 text-sm text-amber-200">
                        Calling getInstance() before initialize() throws an error.
                    </div>
                </section>

                <section className="space-y-4">
                    <h2 className="text-2xl sm:text-3xl font-semibold">3) Ingest data</h2>
                    <p className="text-zinc-400 leading-relaxed">
                        Ingest accepts either inline content or a file path. Metadata is required and should include a type
                        such as text, github, url, or file.
                    </p>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5">
                        <pre className="text-sm text-zinc-200 overflow-auto">{ingestSnippet}</pre>
                    </div>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5 text-sm text-zinc-300">
                        <p className="font-medium text-white mb-2">IngestInput</p>
                        <ul className="space-y-2 list-disc list-inside text-zinc-400">
                            <li>name: string (required)</li>
                            <li>content?: string</li>
                            <li>filePath?: string</li>
                            <li>metadata: &#123; type: &apos;text&apos; | &apos;github&apos; | &apos;url&apos; | &apos;file&apos;; ... &#125; (required)</li>
                        </ul>
                    </div>
                </section>

                <section className="space-y-4">
                    <h2 className="text-2xl sm:text-3xl font-semibold">4) Query data</h2>
                    <p className="text-zinc-400 leading-relaxed">
                        Query takes a question and optional chat history. If top_k is not provided, the SDK defaults it to 5.
                    </p>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5">
                        <pre className="text-sm text-zinc-200 overflow-auto">{querySnippet}</pre>
                    </div>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5 text-sm text-zinc-300">
                        <p className="font-medium text-white mb-2">QueryInput</p>
                        <ul className="space-y-2 list-disc list-inside text-zinc-400">
                            <li>question: string (required)</li>
                            <li>history?: Array&#60;&#123; role: string; content: string &#125;&#62;</li>
                            <li>top_k?: number (defaults to 5)</li>
                        </ul>
                        <p className="font-medium text-white mt-4 mb-2">QueryResponse</p>
                        <ul className="space-y-2 list-disc list-inside text-zinc-400">
                            <li>answer: string</li>
                        </ul>
                    </div>
                </section>

                <section className="rounded-lg border border-zinc-800 bg-zinc-950 p-6">
                    <h2 className="text-xl font-semibold mb-3">Notes from the example app</h2>
                    <ul className="space-y-2 list-disc list-inside text-zinc-400">
                        <li>The example Nest app initializes RavaClient once during bootstrap.</li>
                        <li>The ingest flow in the example uses filePath to read local data.txt.</li>
                        <li>The query flow passes QueryInput directly to client.query().</li>
                    </ul>
                </section>

                <section className="space-y-4">
                    <h2 className="text-2xl sm:text-3xl font-semibold">5) Get API key flow</h2>
                    <p className="text-zinc-400 leading-relaxed">
                        Project creation no longer accepts user_id in the body. Authenticate first, then send
                        a Bearer JWT token to create the project and generate API keys.
                    </p>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5">
                        <pre className="text-sm text-zinc-200 overflow-auto">{getApiKeyFlowSnippet}</pre>
                    </div>
                    <div className="rounded-lg border border-zinc-800 bg-zinc-950 p-5 text-sm text-zinc-300">
                        <p className="font-medium text-white mb-2">Required headers</p>
                        <ul className="space-y-2 list-disc list-inside text-zinc-400">
                            <li>Project creation: Authorization: Bearer &lt;jwt&gt;</li>
                            <li>API key generation: Authorization: Bearer &lt;jwt&gt;</li>
                            <li>Ingest/query calls: X-API-Key: &lt;project-api-key&gt;</li>
                        </ul>
                    </div>
                    <Button asChild className="px-6 py-3 rounded-full bg-white text-black hover:bg-zinc-100 font-medium">
                        <Link href="/auth">Open auth and dashboard flow</Link>
                    </Button>
                </section>
            </main>
        </div>
    )
}
