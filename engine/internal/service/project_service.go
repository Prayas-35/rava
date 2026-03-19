package service

import (
	"context"
	"time"

	"github.com/Prayas-35/ragkit/engine/internal/database"
)

type Project struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	AgentPrompt string    `json:"agent_prompt"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateProjectRequest struct {
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	AgentPrompt string `json:"agent_prompt"`
}

func CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	row := database.DB.QueryRow(ctx,
		`INSERT INTO projects (user_id, name, agent_prompt) VALUES ($1, $2, $3)
		RETURNING id, user_id, name, agent_prompt, created_at`,
		req.UserID,
		req.Name,
		req.AgentPrompt,
	)

	var p Project
	if err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.AgentPrompt, &p.CreatedAt); err != nil {
		return nil, err
	}

	return &p, nil
}

func GetProjectAgentPrompt(ctx context.Context, projectID string) (string, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT agent_prompt FROM projects WHERE id = $1`,
		projectID,
	)

	var agentPrompt string
	if err := row.Scan(&agentPrompt); err != nil {
		return "", err
	}

	return agentPrompt, nil
}

func UpdateProjectAgentPrompt(ctx context.Context, projectID string, agentPrompt string) error {
	_, err := database.DB.Exec(ctx,
		`UPDATE projects SET agent_prompt = $1 WHERE id = $2`,
		agentPrompt,
		projectID,
	)
	return err
}
