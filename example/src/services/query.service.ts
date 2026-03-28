import { Injectable } from "@nestjs/common"
import { RavaClient, type QueryInput, type QueryResponse } from "rava"

@Injectable()
export class QueryService {
  async process(payload: QueryInput): Promise<QueryResponse> {
    const client = RavaClient.getInstance()
    return client.query(payload)
  }
}
