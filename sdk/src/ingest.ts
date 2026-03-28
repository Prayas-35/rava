import { AxiosInstance } from "axios"

export interface IngestMetadata {
    type: "text" | "github" | "url" | "file"
    [key: string]: any
}

export interface IngestInput {
    name: string
    content?: string
    filePath?: string
    metadata?: IngestMetadata
}

export interface IngestResponse {
    status: string
}

async function resolveContent(input: IngestInput): Promise<string> {
    if (input.content && input.content.trim() !== "") {
        return input.content
    }

    if (input.filePath && input.filePath.trim() !== "") {
        // Load fs lazily so browser consumers that only use content are unaffected.
        const { readFile } = await import("node:fs/promises")
        return readFile(input.filePath, "utf8")
    }

    throw new Error("Either content or filePath is required")
}

export async function ingest(
    client: AxiosInstance,
    input: IngestInput
): Promise<IngestResponse> {

    if (!input.metadata) {
        throw new Error("metadata is required")
    }

    const content = await resolveContent(input)

    const res = await client.put("/api/ingest", {
        name: input.name,
        content,
        metadata: input.metadata,
    })

    return res.data
}