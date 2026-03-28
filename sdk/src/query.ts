import { AxiosInstance } from "axios"

export interface QueryInput {
    question: string
    history?: Array<{ role: string; content: string }>
    top_k?: number
}

export interface QueryResponse {
    answer: string
}

export async function query(
    client: AxiosInstance,
    input: QueryInput
): Promise<QueryResponse> {
    const res = await client.post("/api/query", {
        question: input.question,
        history: input.history || [],
        top_k: input.top_k || 5,
    })

    return res.data
}