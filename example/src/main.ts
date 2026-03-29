import "reflect-metadata"
import { NestFactory } from "@nestjs/core"
import { AppModule } from "./app.module"
import { RavaClient } from "@rava-ai/sdk"
import * as dotenv from "dotenv"

dotenv.config()

async function bootstrap() {
  const apiKey = process.env.RAVA_API_KEY
  if (!apiKey) {
    throw new Error("RAVA_API_KEY is required")
  }

  const baseUrl = process.env.RAVA_BASE_URL || "http://localhost:8080"
  const port = Number(process.env.PORT || 3001)

  RavaClient.initialize({
    apiKey,
    baseUrl,
  })

  const app = await NestFactory.create(AppModule)
  await app.listen(port)

  console.log(`Example API running on http://localhost:${port}`)
  console.log(`Rava base URL: ${baseUrl}`)
}

bootstrap().catch((error: unknown) => {
  console.error("Failed to start example API", error)
  process.exit(1)
})
