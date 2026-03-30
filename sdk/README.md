# Rava SDK

A TypeScript/JavaScript SDK for integrating with the Rava RAG (Retrieval-Augmented Generation) API. Easily ingest documents and query them with AI-powered retrieval.

## Features

- 🚀 **Singleton Pattern**: Efficient client management - initialize once, use everywhere
- 📄 **Document Ingestion**: Ingest text, files, GitHub repositories, and URLs
- 🔍 **Smart Querying**: Query your ingested documents with AI-powered retrieval
- 🔐 **Secure Authentication**: Bearer token authentication with API keys
- 📦 **TypeScript Support**: Full type definitions included
- 🎯 **Simple API**: Intuitive methods for ingest and query operations

## Installation

```bash
npm install @rava-ai/sdk
# or
pnpm add @rava-ai/sdk
# or
yarn add @rava-ai/sdk
```

## Quick Start

### Obtain an API Key

Before initializing the SDK, create an API key from the Rava web app:

1. Go to [ravaai.vercel.app](https://ravaai.vercel.app)
2. Sign in to your account
3. Open your dashboard
4. Create a new key and copy it
5. Store it securely (for example, in your `.env` file as `RAVA_API_KEY`)

### 1. Initialize the Client

Initialize the Rava client once at your application startup:

```typescript
import { RavaClient } from "@rava-ai/sdk"

// At your app startup (e.g., main.ts, server.ts)
RavaClient.initialize({
  apiKey: "your-api-key-here",
  baseUrl: "https://rava-ydvd.onrender.com", // optional, defaults to https://rava-ydvd.onrender.com
})
```

### 2. Get Singleton Instance

Use the same client instance throughout your application:

```typescript
import { RavaClient } from "@rava-ai/sdk"

const client = RavaClient.getInstance()
```

### 3. Ingest a Document

```typescript
const response = await client.ingest({
  name: "my-document",
  content: "Your document content here",
  metadata: {
    type: "text", // 'text', 'github', 'url', or 'file'
  },
})

console.log("Ingestion status:", response.status)
```

### 4. Query Documents

```typescript
const response = await client.query({
  question: "What is the main topic?",
  top_k: 5, // number of relevant results to retrieve
})

console.log("Answer:", response.answer)
```

## API Reference

### RavaClient

#### Constructor

```typescript
new RavaClient(config: ClientConfig)
```

#### Static Methods

##### `initialize(config: ClientConfig): RavaClient`

Initialize the singleton instance. **Call this once at application startup.**

**Parameters:**
- `config.apiKey` (string, required): Your Rava API key
- `config.baseUrl` (string, optional): Base URL for the API (defaults to `https://api.rava.dev`)

**Returns:** The initialized RavaClient instance

```typescript
RavaClient.initialize({
  apiKey: "rag_xxxxxxxxxxxx",
  baseUrl: "https://api.rava.dev",
})
```

##### `getInstance(): RavaClient`

Get the singleton instance. **Must call `initialize()` first.**

**Returns:** The RavaClient instance

**Throws:** Error if `initialize()` was not called

```typescript
const client = RavaClient.getInstance()
```

##### `reset(): void`

Reset the singleton instance. Useful for testing or reinitializing.

```typescript
RavaClient.reset()
RavaClient.initialize({ apiKey: "new-key" })
```

#### Instance Methods

##### `ingest(input: IngestInput): Promise<IngestResponse>`

Ingest a document into your RAG system.

**Parameters:**
- `name` (string, required): Name/identifier for the document
- `content` (string, optional): Document content as string
- `filePath` (string, optional): Path to file to read (Node.js only)
- `metadata` (IngestMetadata, required): Metadata about the document

**Metadata Object:**
- `type` (string, required): One of `'text'`, `'github'`, `'url'`, or `'file'`
- Additional metadata fields can be included as needed

**Returns:** Promise resolving to `{ status: string }`

**Examples:**

Ingest text content:
```typescript
const response = await client.ingest({
  name: "product-documentation",
  content: "This is the product documentation...",
  metadata: {
    type: "text",
    version: "1.0",
    category: "docs",
  },
})
```

Ingest from file (Node.js):
```typescript
const response = await client.ingest({
  name: "research-paper",
  filePath: "./paper.txt",
  metadata: {
    type: "file",
    format: "txt",
  },
})
```

##### `query(input: QueryInput): Promise<QueryResponse>`

Query your ingested documents.

**Parameters:**
- `question` (string, required): Your question/query
- `history` (Array, optional): Conversation history for context
  - Each item: `{ role: 'user' | 'assistant', content: string }`
- `top_k` (number, optional): Number of relevant documents to retrieve (default: 5)

**Returns:** Promise resolving to `{ answer: string }`

**Examples:**

Simple query:
```typescript
const response = await client.query({
  question: "What are the main features?",
})
console.log(response.answer)
```

Query with history:
```typescript
const response = await client.query({
  question: "Tell me more about that",
  history: [
    { role: "user", content: "What is RAG?" },
    { role: "assistant", content: "RAG stands for Retrieval-Augmented Generation..." },
  ],
  top_k: 3,
})
```

## Configuration

### ClientConfig

```typescript
interface ClientConfig {
  apiKey: string        // Required: Your Rava API key
  baseUrl?: string      // Optional: API base URL (default: https://api.rava.dev)
}
```

## Usage Examples

### Express.js

```typescript
import express from "express"
import { RavaClient } from "@rava-ai/sdk"

const app = express()

// Initialize at startup
RavaClient.initialize({
  apiKey: process.env.RAVA_API_KEY!,
})

app.post("/api/ingest", async (req, res) => {
  try {
    const client = RavaClient.getInstance()
    const result = await client.ingest({
      name: req.body.name,
      content: req.body.content,
      metadata: { type: "text" },
    })
    res.json(result)
  } catch (error) {
    res.status(500).json({ error: error.message })
  }
})

app.post("/api/query", async (req, res) => {
  try {
    const client = RavaClient.getInstance()
    const result = await client.query({
      question: req.body.question,
    })
    res.json(result)
  } catch (error) {
    res.status(500).json({ error: error.message })
  }
})

app.listen(3000)
```

### Next.js

```typescript
// lib/rava.ts
import { RavaClient } from "@rava-ai/sdk"

// Initialize client at module load time
if (!RavaClient.getInstance()) {
  RavaClient.initialize({
    apiKey: process.env.NEXT_PUBLIC_RAVA_API_KEY!,
  })
}

export const ravaClient = RavaClient.getInstance()
```

```typescript
// app/api/ingest/route.ts
import { ravaClient } from "@/lib/rava"

export async function POST(req: Request) {
  const body = await req.json()
  
  const result = await ravaClient.ingest({
    name: body.name,
    content: body.content,
    metadata: { type: "text" },
  })
  
  return Response.json(result)
}
```

### Testing

```typescript
import { RavaClient } from "@rava-ai/sdk"

describe("Rava Integration", () => {
  beforeEach(() => {
    RavaClient.reset()
    RavaClient.initialize({ apiKey: "test-key" })
  })

  afterEach(() => {
    RavaClient.reset()
  })

  it("should ingest documents", async () => {
    const client = RavaClient.getInstance()
    const result = await client.ingest({
      name: "test-doc",
      content: "Test content",
      metadata: { type: "text" },
    })
    expect(result.status).toBeDefined()
  })
})
```

## Error Handling

The SDK throws errors in the following scenarios:

1. **Not Initialized**: Calling `getInstance()` before `initialize()`
   ```typescript
   try {
     const client = RavaClient.getInstance()
   } catch (error) {
     console.error(error.message)
     // "RavaClient not initialized. Call RavaClient.initialize(config) first."
   }
   ```

2. **Missing Required Fields**: When `content` and `filePath` are both missing
   ```typescript
   try {
     await client.ingest({
       name: "doc",
       metadata: { type: "text" },
     })
   } catch (error) {
     console.error(error.message)
     // "Either content or filePath is required"
   }
   ```

3. **Missing Metadata**: When metadata is not provided
   ```typescript
   try {
     await client.ingest({
       name: "doc",
       content: "content",
       // metadata missing
     })
   } catch (error) {
     console.error(error.message)
     // "metadata is required"
   }
   ```

4. **Network Errors**: When API calls fail
   ```typescript
   try {
     await client.query({ question: "test" })
   } catch (error) {
     console.error(error.message)
     // Network error details
   }
   ```

## Environment Variables

Store your API key securely in environment variables:

```bash
# .env or .env.local
RAVA_API_KEY=your_api_key_here
```

Then use in your code:

```typescript
RavaClient.initialize({
  apiKey: process.env.RAVA_API_KEY!,
})
```

## Building from Source

Clone the repository and build the SDK:

```bash
pnpm install
pnpm run build
```

This generates:
- `dist/index.js` - CommonJS bundle
- `dist/index.d.ts` - TypeScript type definitions

## Development

Watch for changes and rebuild automatically:

```bash
pnpm run dev
```

## Best Practices

1. **Initialize Once**: Call `RavaClient.initialize()` exactly once at application startup
2. **Use getInstance()**: Never create new RavaClient instances in request handlers
3. **Error Handling**: Always wrap SDK calls in try-catch blocks
4. **API Keys**: Store API keys in environment variables, never hardcode them
5. **Type Safety**: Leverage TypeScript for better IDE support and error detection

## Support

For issues, questions, or contributions, please visit the [Rava repository](https://github.com/Prayas-35/rava).

## License

ISC
