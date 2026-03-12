package queue

type IngestJob struct {
	DocumentID string                 `json:"document_id"`
	ProjectID  string                 `json:"project_id"`
	Content    string                 `json:"content"`
	Metadata   map[string]interface{} `json:"metadata"`
}
