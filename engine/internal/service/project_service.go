package service

import (
	"context"
	"time"

	"github.com/Prayas-35/ragkit/engine/internal/database"
)

type Project struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateProjectRequest struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
}

func CreateProject(ctx context.Context, req CreateProjectRequest) (*Project, error) {
	row := database.DB.QueryRow(ctx,
		`INSERT INTO projects (user_id, name) VALUES ($1, $2)
		RETURNING id, user_id, name, created_at`,
		req.UserID,
		req.Name,
	)

	var p Project
	if err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.CreatedAt); err != nil {
		return nil, err
	}

	return &p, nil
}
