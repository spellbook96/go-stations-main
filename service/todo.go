package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	fmt.Printf("Running CreateTODO\n")
	if subject == "" {
		return nil, errors.New("subject must not be empty")
	}
	// s.db.Exec(insert, subject, description)
	s.db.PrepareContext(ctx, insert)
	id, err := s.db.ExecContext(ctx, insert, subject, description)
	idint, _ := id.LastInsertId()
	if err != nil {
		return nil, err
	}
	// TODO := s.db.QueryRow(confirm, subject, description)
	TODO := &model.TODO{
		ID:          idint,
		Subject:     subject,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	fmt.Printf("%+v", TODO)

	return TODO, nil
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	if prevID == 0 {
		return nil, errors.New("prevID must not be empty")
	}
	if size == 0 {
		return nil, errors.New("size must not be empty")
	}
	var TODOS []*model.TODO
	for row := s.db.QueryRow(readWithID, prevID, size); row != nil; row = s.db.QueryRow(readWithID, prevID, size) {
		TODO := &model.TODO{}
		row.Scan(&TODO.ID, &TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt)
		TODOS = append(TODOS, TODO)
	}
	return TODOS, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	if subject == "" {
		return nil, model.ErrNotFound{Message: "subject must not be empty"}
	}
	s.db.PrepareContext(ctx, update)
	result, _ := s.db.ExecContext(ctx, update, subject, description, id)
	num, _ := result.RowsAffected()
	if num == 0 {
		return nil, model.ErrNotFound{Message: "number of rows affected is 0"}
	}
	TODO := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   time.Now(), //後で修正する必要がある
		UpdatedAt:   time.Now(),
	}
	return TODO, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	if len(ids) == 0 {
		return errors.New("ids must not be empty")
	}

	return nil
}
