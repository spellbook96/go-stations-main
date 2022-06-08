package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
		all_read   = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC`
	)

	TODOS := []*model.TODO{}
	fmt.Printf("Running ReadTODO: prevID:%d,size:%d\n", prevID, size)

	if prevID == 0 {
		// Default: read all TODOs
		var rows *sql.Rows
		var err error
		if size == 0 {
			fmt.Println("Running ReadTODO: read all")
			rows, err = s.db.Query(all_read)
			if err != nil {
				return nil, err
			}
		} else {
			rows, err = s.db.Query(read, size)
			if err != nil {
				return nil, err
			}
		}
		defer rows.Close()
		for rows.Next() {
			var todo model.TODO
			err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
			if err != nil {
				return nil, err
			}
			TODOS = append(TODOS, &todo)

		}
	} else {
		rows, _ := s.db.QueryContext(ctx, readWithID, prevID, size)
		// fmt.Printf("%+v\n", TODOS)
		defer rows.Close()
		for rows.Next() {
			TODO := &model.TODO{}
			if err := rows.Scan(&TODO.ID, &TODO.Subject, &TODO.Description, &TODO.CreatedAt, &TODO.UpdatedAt); err != nil {
				return nil, err
			}
			// fmt.Printf("%+v\n", TODO)
			TODOS = append(TODOS, TODO)
		}
	}
	// fmt.Printf("%+v\n", TODOS)
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
	// s.db.PrepareContext(ctx, update)
	// curr, _ := s.db.QueryContext(ctx, confirm, id)
	// currTODO := &model.TODO{}
	// curr.Scan(&currTODO.Subject, &currTODO.Description, &currTODO.CreatedAt, &currTODO.UpdatedAt)
	// createdTime := currTODO.CreatedAt
	result, _ := s.db.ExecContext(ctx, update, subject, description, id)
	num, _ := result.RowsAffected()
	if num == 0 {
		return nil, model.ErrNotFound{Message: "number of rows affected is 0"}
	}
	TODO := &model.TODO{
		ID:          id,
		Subject:     subject,
		Description: description,
		CreatedAt:   time.Now(),
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

	query := fmt.Sprintf(deleteFmt, strings.Repeat(",?", len(ids)-1))
	fmt.Println(query)
	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	args := make([]interface{}, len(ids))
	for i := range args {
		args[i] = ids[i]
	}
	rows, err := stmt.ExecContext(ctx, args...)

	if err != nil {
		return err
	}
	cnt, _ := rows.RowsAffected()
	if cnt == 0 {
		return model.ErrNotFound{Message: "ID not found"}
	}
	return nil
}
