package todo

import (
	"context"
	"fmt"
	"strings"
)

type service struct {
	todoStore Store
}

func NewService(
	todoStore Store,
) (*service, error) {
	return &service{
		todoStore: todoStore,
	}, nil
}

func (s *service) Create(ctx context.Context, title string, content string) error {
	td := New(title, content)

	err := s.todoStore.Create(ctx, td)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) List(ctx context.Context) (string, error) {
	const limit = 10
	const offset = 0

	list, err := s.todoStore.List(ctx, limit, offset)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}

	for _, td := range list {
		fmt.Fprintf(&builder, "[%d] [ ] %s\n", td.Num, td.Title)
	}

	return builder.String(), nil
}

func (s *service) Sort(ctx context.Context) error {
	list, err := s.todoStore.List(ctx, 10, 0)
	if err != nil {
		return err
	}

	err = s.todoStore.Sort(ctx, list)
	if err != nil {
		return err
	}

	return nil
}
