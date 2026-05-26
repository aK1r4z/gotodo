package todo

import (
	"context"
	"fmt"
	"strings"
)

// 待办服务
type service struct {
	todoStore Store
}

// 创建待办服务
func NewService(
	todoStore Store,
) (*service, error) {
	return &service{
		todoStore: todoStore,
	}, nil
}

// 创建待办
func (s *service) Create(ctx context.Context, title string, content string) error {
	td := New(title, content)

	err := s.todoStore.Create(ctx, td)
	if err != nil {
		return err
	}

	return nil
}

// 列出所有待办，返回字符串
func (s *service) List(ctx context.Context) (string, error) {
	const limit = 10
	const offset = 0

	list, err := s.todoStore.List(ctx, limit, offset)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}

	for _, r := range list {
		fmt.Fprintf(&builder, "[%d] [ ] %s\n", r.Num, r.Todo.Title)
	}

	return builder.String(), nil
}
