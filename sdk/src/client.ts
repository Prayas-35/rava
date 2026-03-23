import axios, { type AxiosInstance } from "axios"
import { query, QueryInput, QueryResponse } from "./query"
import { ingest, IngestInput, IngestResponse } from "./ingest"

export interface ClientConfig {
    apiKey: string
    baseUrl?: string
}

export class RagClient {
    private client: AxiosInstance
    private static instance: RagClient | null = null
    private static config: ClientConfig | null = null

    constructor(config: ClientConfig) {
        this.client = axios.create({
            baseURL: config.baseUrl || "https://api.ragkit.dev",
            headers: {
                Authorization: `Bearer ${config.apiKey}`,
                "Content-Type": "application/json",
            },
        })
    }

    /**
     * Initialize the singleton instance.
     * Call this once at your app startup before using the client.
     */
    static initialize(config: ClientConfig): RagClient {
        if (!RagClient.instance) {
            RagClient.config = config
            RagClient.instance = new RagClient(config)
        }
        return RagClient.instance
    }

    /**
     * Get the singleton instance.
     * Make sure to call initialize() first.
     */
    static getInstance(): RagClient {
        if (!RagClient.instance) {
            throw new Error(
                "RagClient not initialized. Call RagClient.initialize(config) first."
            )
        }
        return RagClient.instance
    }

    /**
     * Reset the singleton instance (useful for testing).
     */
    static reset(): void {
        RagClient.instance = null
        RagClient.config = null
    }

    async query(input: QueryInput) {
        return query(this.client, input)
    }

    async ingest(input: IngestInput) {
        return ingest(this.client, input)
    }
}