package service

import (
	"github.com/kazakovichna/todoListPrjct"
	"github.com/kazakovichna/todoListPrjct/pkg/repository"
)

type TodoListService struct {
	repo repository.TodoList
	deskRepo repository.TodoDesk
}

func NewTodoListService(repo repository.TodoList, deskRepo repository.TodoDesk) *TodoListService {
	return &TodoListService{repo: repo, deskRepo: deskRepo}
}
func (s *TodoListService) Create(userId, deskId int, list todoListPrjct.ListTable) (int, error) {
	_, err := s.deskRepo.GetDeskById(userId, deskId)
	if err != nil {
		// desk doesn't exist or doesn't belong to user
		return 0, err
	}

	return s.repo.Create(deskId, list)
}

func (s *TodoListService) GetAll(userId, deskId int) ([]todoListPrjct.ListTable, error) {
	return s.repo.GetAll(userId, deskId)
}

func (s *TodoListService) GetById(userId, deskId, listId int) (todoListPrjct.ListTable, error) {
	return s.repo.GetById(userId, deskId, listId)
}

func (s *TodoListService) Delete(userId, deskId, listId int) error {
	return s.repo.Delete(userId, deskId, listId)
}

func (s *TodoListService) Update(userId, deskId, listId int, input todoListPrjct.UpdateListInput) error {
	if err := input.ValidateList(); err != nil {
		return err
	}
	return s.repo.Update(userId, deskId, listId, input)
}


