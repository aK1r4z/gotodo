package postgres

import (
	"context"

	"github.com/aK1r4z/gotodo/internal/todo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type todoStore struct {
	db *pgxpool.Pool
}

func (d *db) TODOStore() *todoStore {
	return &todoStore{d.pool}
}

func (s *todoStore) Create(ctx context.Context, todo *todo.TODO) error {
	sql := `
		insert into todos
			(title, content, completed, created_at, deleted_at)
		values
			($1, $2, $3, $4, $5)
		;
	`

	_, err := s.db.Exec(ctx, sql, todo.Title, todo.Content, todo.Completed, todo.CreatedAt, todo.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *todoStore) List(ctx context.Context, limit uint32, offset uint32) ([]todo.ListResult, error) {
	sql := `
		select
			row_number() over (order by created_at) as num,
			id, title, content, completed, created_at, deleted_at
		from
			todos
		limit
			$1
		offset
			$2
		;
	`

	rows, err := s.db.Query(ctx, sql, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []todo.ListResult

	for rows.Next() {
		var s todo.ListResult

		if err := rows.Scan(&s.Num, &s.Todo.ID, &s.Todo.Title, &s.Todo.Content, &s.Todo.Completed, &s.Todo.CreatedAt, &s.Todo.DeletedAt); err != nil {
			return nil, err
		}

		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (s *todoStore) GetByNumber(ctx context.Context, num uint32) (*todo.TODO, error) {
	sql := `
		SELECT
			id, title, content, completed, created_at, deleted_at
		FROM (
			SELECT
				ROW_NUMBER() OVER (ORDER BY created_at) AS num,
				id, title, content, completed, created_at, deleted_at
			FROM todos
		) AS t
		WHERE num = $1
		;
	`

	row := s.db.QueryRow(ctx, sql, num)

	var td todo.TODO
	if err := row.Scan(&td.ID, &td.Title, &td.Content, &td.Completed, &td.CreatedAt, &td.DeletedAt); err != nil {
		return nil, err
	}

	return &td, nil
}

func (s *todoStore) Update(ctx context.Context, todo *todo.TODO) error {
	sql := `
		update
			todos
		set
			num = $2, title = $3, content = $4, completed = $5, created_at = $6, deleted_at = $7
		where
			id = $1
		;
	`

	_, err := s.db.Exec(ctx, sql, todo.ID, todo.Title, todo.Content, todo.Completed, todo.CreatedAt, todo.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *todoStore) Delete(ctx context.Context, todo *todo.TODO) error {
	sql := `
		delete from
			todos
		where
			id = $1
		;
	`

	_, err := s.db.Exec(ctx, sql, todo.ID)
	if err != nil {
		return err
	}

	return nil
}
