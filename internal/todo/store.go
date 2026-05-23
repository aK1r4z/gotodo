package todo

import "context"

type Store interface {
	Create(context.Context, *TODO) error

	List(ctx context.Context, limit uint32, offset uint32) ([]TODO, error)

	Update(ctx context.Context, todo *TODO) error
	Sort(ctx context.Context, todos []TODO) error

	Delete(ctx context.Context, todo *TODO) error
}
