import { Injectable } from "@nestjs/common"
import { resolve } from "node:path"
import { RavaClient, type IngestInput, type IngestMetadata, type IngestResponse } from "rava"

interface IngestRequest {
  name: string
  metadata?: IngestMetadata
}

@Injectable()
export class IngestService {
  async process(payload: IngestRequest): Promise<IngestResponse> {
    const client = RavaClient.getInstance()
    const dataFilePath = resolve(__dirname, "../../data.txt")
    console.log(`Ingesting data from file: ${dataFilePath}`)

    const ingestPayload: IngestInput = {
      name: payload.name,
      filePath: dataFilePath,
      metadata: payload.metadata ?? { type: "text" },
    }

    return client.ingest(ingestPayload)
  }
}
