package todo

import (
	"context"
	"fmt"
	"strings"
	"time"
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

// 查询指定顺序编号的待办项的详细信息
func (s *service) GetByNumber(ctx context.Context, num uint32) (string, error) {
	td, err := s.todoStore.GetByNumber(ctx, num)
	if err != nil {
		return "", err
	}

	if td == nil {
		return "Err: Not found", nil
	}

	result := fmt.Sprintf("Title: %s", td.Title)

	if td.Content != "" {
		result += fmt.Sprintf("\nContent: %s", td.Content)
	}

	result += fmt.Sprintf("\nCreated At: %s", td.CreatedAt.Format(time.DateTime))

	if td.DeletedAt != nil {
		result += fmt.Sprintf("\nDeleted At: %s", td.DeletedAt.Format(time.DateTime))
	}

	if td.Completed {
		result += "\nCompleted: TRUE"
	}

	return result, nil
}
