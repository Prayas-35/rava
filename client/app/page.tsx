import { Button } from '@/components/ui/button'
import { ArrowRight, Code2, Zap, Database, Brain, GitBranch, Shield } from 'lucide-react'
import Link from 'next/link'

export default function Home() {
  return (
    <div className="dark min-h-screen bg-black text-white">
      {/* Navigation */}
      <nav className="border-b border-zinc-900 bg-black/50 backdrop-blur-md sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-linear-to-br from-indigo-500 to-indigo-600 rounded flex items-center justify-center">
              <span className="text-white font-bold text-sm">R</span>
            </div>
            <span className="text-xl font-bold"> Rava </span>
          </div>
          <div className="hidden md:flex items-center gap-8">
            <a href="#how" className="text-sm text-zinc-300 hover:text-white transition">How it works</a>
            <a href="#features" className="text-sm text-zinc-300 hover:text-white transition">Features</a>
            <a href="#usecases" className="text-sm text-zinc-300 hover:text-white transition">Use cases</a>
            <Link href="#" className="px-4 py-2 rounded-full bg-white text-black text-sm font-medium hover:bg-zinc-100 transition">Get API key</Link>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="relative px-4 sm:px-6 lg:px-8 py-20 sm:py-32">
        <div className="max-w-7xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div>
              <div className="inline-block mb-4 px-3 py-1 rounded-full bg-zinc-900 border border-zinc-800">
                <span className="text-xs text-indigo-400">RAG-as-a-Service</span>
              </div>
              <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold mb-6 leading-tight">
                Add RAG to your app in minutes
              </h1>
              <p className="text-lg sm:text-xl text-zinc-400 mb-8 leading-relaxed">
                No complex pipelines. No vector database management. Just one-line SDK integration and you&apos;re ready to retrieve and generate.
              </p>
              <div className="flex flex-col sm:flex-row gap-4">
                <Button className="h-12 px-6 rounded-full bg-white text-black hover:bg-zinc-100 font-medium text-base">
                  Get Started <ArrowRight className="ml-2 w-4 h-4" />
                </Button>
                <Button variant="outline" className="h-12 px-6 rounded-full border-zinc-700 text-white hover:bg-zinc-900 font-medium text-base">
                  View Docs
                </Button>
              </div>
            </div>

            {/* Code Block */}
            <div className="relative">
              <div className="bg-gradient-to-br from-zinc-900 to-black border border-zinc-800 rounded-lg p-6 backdrop-blur">
                <div className="flex items-center gap-2 mb-4 pb-4 border-b border-zinc-800">
                  <div className="w-3 h-3 rounded-full bg-red-500"></div>
                  <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                  <div className="w-3 h-3 rounded-full bg-green-500"></div>
                  <span className="text-xs text-zinc-500 ml-2">quick-start.js</span>
                </div>
                <pre className="text-sm text-zinc-300 font-mono overflow-auto">
                  {`import Rava from '@rava/sdk';

// Initialize with API key
const rava = new Rava({
  apiKey: process.env.RAVA_API_KEY
});

try {
  // Ingest data from source
  await rava.ingest({
    source: 'github',
    repo: 'your-org/repo',
    metadata: { type: 'docs' }
  });

  // Query with RAG
  const result = await rava.query({
    question: 'How do I use RAG?',
    topK: 5
  });

  console.log(result.answer);
} catch (error) {
  console.error('RAG error:', error);
}`}
                </pre>
              </div>
              <div className="absolute -top-4 -right-4 w-20 h-20 bg-indigo-500/10 rounded-full blur-2xl"></div>
            </div>
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section id="how" className="px-4 sm:px-6 lg:px-8 py-20 border-t border-zinc-900">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl sm:text-5xl font-bold mb-4">How it works</h2>
            <p className="text-lg text-zinc-400">Three simple steps to add RAG to your application</p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            {/* Step 1 */}
            <div className="relative">
              <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8 h-full">
                <div className="w-12 h-12 rounded-lg bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center mb-4">
                  <Database className="w-6 h-6 text-indigo-400" />
                </div>
                <h3 className="text-xl font-semibold mb-3">Ingest</h3>
                <p className="text-zinc-400">Upload your data—text files, PDFs, GitHub repos, URLs. Rava handles chunking and embedding automatically.</p>
              </div>
              <div className="hidden md:block absolute -right-4 top-8 w-8 h-0.5 bg-gradient-to-r from-indigo-500 to-transparent"></div>
            </div>

            {/* Step 2 */}
            <div className="relative">
              <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8 h-full">
                <div className="w-12 h-12 rounded-lg bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center mb-4">
                  <Brain className="w-6 h-6 text-indigo-400" />
                </div>
                <h3 className="text-xl font-semibold mb-3">Retrieve</h3>
                <p className="text-zinc-400">Query using natural language. Our vector search instantly finds the most relevant context from your data.</p>
              </div>
              <div className="hidden md:block absolute -right-4 top-8 w-8 h-0.5 bg-gradient-to-r from-indigo-500 to-transparent"></div>
            </div>

            {/* Step 3 */}
            <div>
              <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8 h-full">
                <div className="w-12 h-12 rounded-lg bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center mb-4">
                  <Zap className="w-6 h-6 text-indigo-400" />
                </div>
                <h3 className="text-xl font-semibold mb-3">Generate</h3>
                <p className="text-zinc-400">Get back high-quality answers powered by Groq&apos;s fast LLM inference and your custom data.</p>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Features Grid */}
      <section id="features" className="px-4 sm:px-6 lg:px-8 py-20 border-t border-zinc-900">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl sm:text-5xl font-bold mb-4">Powerful features</h2>
            <p className="text-lg text-zinc-400">Everything you need to build production-ready RAG applications</p>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8">
              <div className="flex items-start gap-4">
                <Code2 className="w-6 h-6 text-indigo-400 mt-1 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold mb-2">One-line SDK</h3>
                  <p className="text-sm text-zinc-400">Drop our npm package into your project and start building instantly.</p>
                </div>
              </div>
            </div>

            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8">
              <div className="flex items-start gap-4">
                <GitBranch className="w-6 h-6 text-indigo-400 mt-1 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold mb-2">Multi-source ingestion</h3>
                  <p className="text-sm text-zinc-400">Ingest from text, files, GitHub repositories, and URLs seamlessly.</p>
                </div>
              </div>
            </div>

            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8">
              <div className="flex items-start gap-4">
                <Database className="w-6 h-6 text-indigo-400 mt-1 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold mb-2">Built on pgvector</h3>
                  <p className="text-sm text-zinc-400">Fast vector search powered by PostgreSQL and pgvector for low latency.</p>
                </div>
              </div>
            </div>

            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8">
              <div className="flex items-start gap-4">
                <Zap className="w-6 h-6 text-indigo-400 mt-1 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold mb-2">Groq-powered inference</h3>
                  <p className="text-sm text-zinc-400">Lightning-fast LLM generation using Groq for sub-second responses.</p>
                </div>
              </div>
            </div>

            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8">
              <div className="flex items-start gap-4">
                <Shield className="w-6 h-6 text-indigo-400 mt-1 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold mb-2">Project isolation</h3>
                  <p className="text-sm text-zinc-400">Each project is isolated with its own vector space and API keys.</p>
                </div>
              </div>
            </div>

            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8">
              <div className="flex items-start gap-4">
                <Code2 className="w-6 h-6 text-indigo-400 mt-1 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold mb-2">Go backend</h3>
                  <p className="text-sm text-zinc-400">Scalable infrastructure built in Go for maximum reliability and performance.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* SDK Example */}
      <section className="px-4 sm:px-6 lg:px-8 py-20 border-t border-zinc-900">
        <div className="max-w-7xl mx-auto">
          <div className="grid lg:grid-cols-2 gap-12 items-center">
            <div>
              <h2 className="text-4xl sm:text-5xl font-bold mb-6">Developer experience</h2>
              <p className="text-lg text-zinc-400 mb-6">
                Simple, intuitive APIs designed for developers. No boilerplate, no headaches.
              </p>
              <ul className="space-y-4 mb-8">
                {['Automatic chunking & embedding', 'Natural language queries', 'Streaming responses', 'Idempotent ingestion'].map((item) => (
                  <li key={item} className="flex items-center gap-3">
                    <div className="w-1.5 h-1.5 rounded-full bg-indigo-400"></div>
                    <span>{item}</span>
                  </li>
                ))}
              </ul>
              <Button className="px-6 py-3 rounded-full bg-white text-black hover:bg-zinc-100 font-medium">
                Read the docs
              </Button>
            </div>

            <div className="bg-gradient-to-br from-zinc-900 to-black border border-zinc-800 rounded-lg p-6 backdrop-blur">
              <div className="flex items-center gap-2 mb-4 pb-4 border-b border-zinc-800">
                <div className="w-3 h-3 rounded-full bg-red-500"></div>
                <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
                <div className="w-3 h-3 rounded-full bg-green-500"></div>
                <span className="text-xs text-zinc-500 ml-2">example.js</span>
              </div>
              <pre className="text-sm text-zinc-300 font-mono overflow-auto">
                {`// Query with streaming
const result = await rava.query({
  question: 'What is RAG?',
  topK: 5,
  streaming: true
});

// Handle streaming response
for await (const chunk of result.stream()) {
  process.stdout.write(chunk.delta);
}

// Get sources and metadata
const answer = await result.text();
const sources = result.sources?.map(
  s => s.filename
) ?? [];

return {
  answer,
  sources
};`}
              </pre>
            </div>
          </div>
        </div>
      </section>

      {/* Use Cases */}
      <section id="usecases" className="px-4 sm:px-6 lg:px-8 py-20 border-t border-zinc-900">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl sm:text-5xl font-bold mb-4">Built for every use case</h2>
            <p className="text-lg text-zinc-400">From chatbots to internal tools, Rava powers them all</p>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            {[
              { title: 'AI Chatbots', desc: 'Build context-aware chatbots that understand your data' },
              { title: 'Documentation Search', desc: 'Semantic search over your entire documentation' },
              { title: 'Internal Tools', desc: 'Create internal assistants for your team' },
              { title: 'AI Copilots', desc: 'Embed AI assistance directly into your products' }
            ].map((usecase) => (
              <div key={usecase.title} className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-6">
                <h3 className="text-lg font-semibold mb-2">{usecase.title}</h3>
                <p className="text-zinc-400">{usecase.desc}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Architecture */}
      <section className="px-4 sm:px-6 lg:px-8 py-20 border-t border-zinc-900">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl sm:text-5xl font-bold mb-4">Built for scale</h2>
            <p className="text-lg text-zinc-400">Enterprise-grade infrastructure powering production applications</p>
          </div>

          <div className="grid md:grid-cols-3 gap-8">
            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8 text-center">
              <div className="text-3xl font-bold text-indigo-400 mb-2">Go</div>
              <h3 className="font-semibold mb-2">Scalable Backend</h3>
              <p className="text-sm text-zinc-400">High-performance backend built in Go for reliability and scale</p>
            </div>
            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8 text-center">
              <div className="text-3xl font-bold text-indigo-400 mb-2">pgvector</div>
              <h3 className="font-semibold mb-2">Vector Storage</h3>
              <p className="text-sm text-zinc-400">Distributed vector search with PostgreSQL and pgvector</p>
            </div>
            <div className="bg-zinc-900/50 border border-zinc-800 rounded-lg p-8 text-center">
              <div className="text-3xl font-bold text-indigo-400 mb-2">Groq</div>
              <h3 className="font-semibold mb-2">LLM Inference</h3>
              <p className="text-sm text-zinc-400">Sub-second LLM responses via Groq&apos;s inference network</p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="px-4 sm:px-6 lg:px-8 py-20 border-t border-zinc-900">
        <div className="max-w-4xl mx-auto text-center">
          <h2 className="text-4xl sm:text-5xl font-bold mb-6">Ready to build with RAG?</h2>
          <p className="text-lg text-zinc-400 mb-8">
            Get your API key and start building in minutes. No credit card required.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button className="h-12 px-8 rounded-full bg-white text-black hover:bg-zinc-100 font-medium text-base">
              Get API Key <ArrowRight className="ml-2 w-4 h-4" />
            </Button>
            <Button variant="outline" className="h-12 px-8 rounded-full border-zinc-700 text-white hover:bg-zinc-900 font-medium text-base">
              View Documentation
            </Button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="border-t border-zinc-900 px-4 sm:px-6 lg:px-8 py-12 bg-black/50">
        <div className="max-w-7xl mx-auto">
          <div className="grid md:grid-cols-4 gap-8 mb-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="w-6 h-6 bg-gradient-to-br from-indigo-500 to-indigo-600 rounded flex items-center justify-center">
                  <span className="text-white font-bold text-xs">R</span>
                </div>
                <span className="font-bold">Rava</span>
              </div>
              <p className="text-sm text-zinc-500">RAG made simple.</p>
            </div>
            <div>
              <h4 className="font-semibold mb-4 text-sm">Product</h4>
              <ul className="space-y-2 text-sm text-zinc-400">
                <li><a href="#" className="hover:text-white transition">Features</a></li>
                <li><a href="#" className="hover:text-white transition">Pricing</a></li>
                <li><a href="#" className="hover:text-white transition">Documentation</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4 text-sm">Company</h4>
              <ul className="space-y-2 text-sm text-zinc-400">
                <li><a href="#" className="hover:text-white transition">Blog</a></li>
                <li><a href="#" className="hover:text-white transition">Twitter</a></li>
                <li><a href="#" className="hover:text-white transition">GitHub</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-4 text-sm">Legal</h4>
              <ul className="space-y-2 text-sm text-zinc-400">
                <li><a href="#" className="hover:text-white transition">Privacy</a></li>
                <li><a href="#" className="hover:text-white transition">Terms</a></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-zinc-900 pt-8 flex flex-col md:flex-row items-center justify-between">
            <p className="text-sm text-zinc-500">© 2024 Rava. All rights reserved.</p>
            <div className="flex gap-4 mt-4 md:mt-0">
              <a href="#" className="text-sm text-zinc-500 hover:text-white transition">Status</a>
              <a href="#" className="text-sm text-zinc-500 hover:text-white transition">Security</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  )
}
