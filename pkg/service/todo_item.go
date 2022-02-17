package service

import (
	"github.com/kazakovichna/todoListPrjct"
	"github.com/kazakovichna/todoListPrjct/pkg/repository"
)

type TodoItemService struct {
	repo repository.TodoItem
	listRepo repository.TodoList
	deskRepo repository.TodoDesk
}

func NewTodoItemService(repo repository.TodoItem, listRepo repository.TodoList, deskRepo repository.TodoDesk) *TodoItemService {
	return &TodoItemService{repo: repo, listRepo: listRepo, deskRepo: deskRepo}
}

func (s *TodoItemService) Create(listId int, input todoListPrjct.ItemTable) (int, error) {
	return s.repo.Create(listId, input)
}

func (s *TodoItemService) GetAllItems(userId, deskId int) ([]todoListPrjct.AllDesksItems, error) {
	return s.repo.GetAllItems(userId, deskId)
}

func (s *TodoItemService) GetItemById(userId, itemId int) (todoListPrjct.ItemTable, error) {
	return s.repo.GetItemById(userId, itemId)
}

func (s *TodoItemService) UpdateItem(userId, itemId int, input todoListPrjct.UpdateItemInput) error {
	if err := input.ValidateItem(); err != nil {
		return err
	}
	return s.repo.UpdateItem(userId, itemId, input)
}

func (s *TodoItemService) DeleteItem(userId, itemId int) error {
	return s.repo.DeleteItem(userId, itemId)
}
