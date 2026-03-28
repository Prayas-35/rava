import { Module } from "@nestjs/common"
import { AppController } from "./app.controller"
import { IngestService } from "./services/ingest.service"
import { QueryService } from "./services/query.service"

@Module({
  controllers: [AppController],
  providers: [IngestService, QueryService],
})
export class AppModule {}
