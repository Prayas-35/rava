package vector

import (
	"context"

	"github.com/Prayas-35/ragkit/engine/config"
	"google.golang.org/genai"
)

type GeminiEmbedder struct {
	client *genai.Client
}

func NewGeminiEmbedder(ctx context.Context) (*GeminiEmbedder, error) {

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: config.LoadConfig().GEMINI_API_KEY,
	})

	if err != nil {
		return nil, err
	}

	return &GeminiEmbedder{client: client}, nil
}

func (g *GeminiEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {

	var contents []*genai.Content

	for _, t := range texts {
		contents = append(contents,
			genai.NewContentFromText(t, genai.RoleUser))
	}

	resp, err := g.client.Models.EmbedContent(
		ctx,
		"gemini-embedding-001",
		contents,
		nil,
	)

	if err != nil {
		return nil, err
	}

	var vectors [][]float32

	for _, e := range resp.Embeddings {
		vectors = append(vectors, e.Values)
	}

	return vectors, nil
}
