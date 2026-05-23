package postgres

import (
	"context"
	"fmt"

	"github.com/aK1r4z/gotodo/internal/todo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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
			(num, title, content, completed, created_at, deleted_at)
		values
			($1, $2, $3, $4, $5, $6)
		;
	`

	_, err := s.db.Exec(ctx, sql, todo.Num, todo.Title, todo.Content, todo.Completed, todo.CreatedAt, todo.DeletedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *todoStore) List(ctx context.Context, limit uint32, offset uint32) ([]todo.TODO, error) {
	sql := `
		select
			id, num, title, content, completed, created_at, deleted_at
		from
			todos
		order by
			created_at
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

	var result []todo.TODO

	for rows.Next() {
		var td todo.TODO

		if err := rows.Scan(&td.ID, &td.Num, &td.Title, &td.Content, &td.Completed, &td.CreatedAt, &td.DeletedAt); err != nil {
			return nil, err
		}

		result = append(result, td)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
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

func (s *todoStore) Sort(ctx context.Context, list []todo.TODO) error {
	tx, err := s.db.Begin(ctx)
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	temp := `
		create temporary table tmp_order (
			id uuid primary key not null,
			num int not null
		)
		;
	`

	if _, err := s.db.Exec(ctx, temp); err != nil {
		return err
	}

	tempOrders := make([][]any, len(list))
	for i, td := range list {
		tempOrders[i] = []any{td.ID, i + 1}
	}

	cnt, err := s.db.CopyFrom(ctx, pgx.Identifier{"tmp_order"}, []string{"id", "num"}, pgx.CopyFromRows(tempOrders))
	if err != nil {
		return err
	}
	if cnt != int64(len(list)) {
		return fmt.Errorf("store/postgres/todo.go func (s *todoStore) Sort() >> cnt != len(list)")
	}

	sql := `
		update
			todos
		set
			num = tmp.num
		from
			tmp_order as tmp
		where
			todos.id = tmp.id
		;
	`

	cmdTag, err := s.db.Exec(ctx, sql)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() != int64(len(list)) {
		return fmt.Errorf("expected %d rows updated, got %d", len(list), cmdTag.RowsAffected())
	}

	return tx.Commit(ctx)
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
