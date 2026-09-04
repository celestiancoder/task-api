package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"task-api/internal/models"
)

type TaskRepository struct {
	db *pgxpool.Pool   // this lower case is private to the package and unexported
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, name string) (*models.Task, error) {
	query := `
		INSERT INTO tasks (name, completed) 
		VALUES ($1, $2) 
		RETURNING id, created_at, updated_at`

	var t models.Task
	t.Name = name
	t.Completed = false

	err := r.db.QueryRow(ctx, query, name, false).Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) GetAll(ctx context.Context, params models.TaskQueryParams) ([]models.Task, error) {
	offset := (params.Page - 1) * params.Limit

	baseQuery := `SELECT id, name, completed, created_at, updated_at FROM tasks`
	whereClauses := []string{}
	args := []any{}
	argCount := 1

	if params.Completed != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("completed = $%d", argCount))
		args = append(args, params.Completed == "true")
		argCount++
	}

	if params.Search != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argCount))
		args = append(args, "%"+params.Search+"%")
		argCount++
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	validSortColumns := map[string]string{
		"id":         "id",
		"name":       "name",
		"created_at": "created_at",
		"updated_at": "updated_at",
	}

	sortColumn, exists := validSortColumns[params.Sort]
	if !exists {
		sortColumn = "id"
	}

	sortOrder := "ASC"
	if strings.ToLower(params.Order) == "desc" {
		sortOrder = "DESC"
	}

	baseQuery += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", sortColumn, sortOrder, argCount, argCount+1)
	args = append(args, params.Limit, offset)

	rows, err := r.db.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Name, &t.Completed, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) GetByID(ctx context.Context, id int) (*models.Task, error) {
	query := `SELECT id, name, completed, created_at, updated_at FROM tasks WHERE id = $1`
	var t models.Task
	err := r.db.QueryRow(ctx, query, id).Scan(&t.ID, &t.Name, &t.Completed, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) Update(ctx context.Context, id int, name string, completed bool) (*models.Task, error) {
	query := `
		UPDATE tasks 
		SET name = $1, completed = $2, updated_at = NOW() 
		WHERE id = $3 
		RETURNING created_at, updated_at`

	var t models.Task
	t.ID = id
	t.Name = name
	t.Completed = completed

	err := r.db.QueryRow(ctx, query, name, completed, id).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}