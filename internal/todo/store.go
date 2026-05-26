package todo

import (
	"context"
)

type ListResult struct {
	Num  int32
	Todo TODO
}

type Store interface {
	Create(context.Context, *TODO) error

	List(ctx context.Context, limit uint32, offset uint32) ([]ListResult, error)
	GetByNumber(ctx context.Context, num uint32) (*TODO, error)

	Update(ctx context.Context, todo *TODO) error

	Delete(ctx context.Context, todo *TODO) error
}
