import axios, { type AxiosInstance } from "axios"
import { query, QueryInput, QueryResponse } from "./query"
import { ingest, IngestInput, IngestResponse } from "./ingest"

export interface ClientConfig {
    apiKey: string
    baseUrl?: string
}

export class RavaClient {
    private client: AxiosInstance
    private static instance: RavaClient | null = null
    private static config: ClientConfig | null = null

    constructor(config: ClientConfig) {
        this.client = axios.create({
            baseURL: config.baseUrl || "https://api.rava.dev",
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
    static initialize(config: ClientConfig): RavaClient {
        if (!RavaClient.instance) {
            RavaClient.config = config
            RavaClient.instance = new RavaClient(config)
        }
        return RavaClient.instance
    }

    /**
     * Get the singleton instance.
     * Make sure to call initialize() first.
     */
    static getInstance(): RavaClient {
        if (!RavaClient.instance) {
            throw new Error(
                "RavaClient not initialized. Call RavaClient.initialize(config) first."
            )
        }
        return RavaClient.instance
    }

    /**
     * Reset the singleton instance (useful for testing).
     */
    static reset(): void {
        RavaClient.instance = null
        RavaClient.config = null
    }

    async query(input: QueryInput) {
        return query(this.client, input)
    }

    async ingest(input: IngestInput) {
        return ingest(this.client, input)
    }
}