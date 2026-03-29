package service

import (
	"context"
	"errors"
	"time"

	"github.com/Prayas-35/ragkit/engine/internal/database"
)

var ErrProjectNotFoundOrUnauthorized = errors.New("project not found or unauthorized")

type Project struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	AgentPrompt string    `json:"agent_prompt"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateProjectRequest struct {
	UserID      string `json:"-"`
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

func ListProjectsByUser(ctx context.Context, userID string) ([]Project, error) {
	rows, err := database.DB.Query(ctx,
		`SELECT id, user_id, name, agent_prompt, created_at
		 FROM projects
		 WHERE user_id = $1
		 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]Project, 0)
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.AgentPrompt, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return projects, nil
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

func UserOwnsProject(ctx context.Context, projectID, userID string) (bool, error) {
	row := database.DB.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)`,
		projectID,
		userID,
	)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func UpdateProjectAgentPromptForUser(ctx context.Context, projectID, userID, agentPrompt string) error {
	result, err := database.DB.Exec(ctx,
		`UPDATE projects SET agent_prompt = $1 WHERE id = $2 AND user_id = $3`,
		agentPrompt,
		projectID,
		userID,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProjectNotFoundOrUnauthorized
	}

	return nil
}
