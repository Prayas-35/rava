import { BadRequestException, Body, Controller, Post } from "@nestjs/common"
import { type IngestMetadata, type IngestResponse, type QueryInput, type QueryResponse } from "@rava-ai/sdk"
import { IngestService } from "./services/ingest.service"
import { QueryService } from "./services/query.service"

interface IngestRequestBody {
  name: string
  metadata?: IngestMetadata
}

@Controller()
export class AppController {
  constructor(
    private readonly ingestService: IngestService,
    private readonly queryService: QueryService
  ) {}

  @Post("/ingest")
  async ingest(@Body() body: IngestRequestBody): Promise<IngestResponse> {
    if (!body?.name) {
      throw new BadRequestException("Ingest payload must include name")
    }

    return this.ingestService.process(body)
  }

  @Post("/query")
  async query(@Body() body: QueryInput): Promise<QueryResponse> {
    if (!body?.question) {
      throw new BadRequestException("Query payload must include question")
    }

    return this.queryService.process(body)
  }
}
