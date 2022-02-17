package service

import (
	"github.com/kazakovichna/todoListPrjct"
	"github.com/kazakovichna/todoListPrjct/pkg/repository"
)

type TodoDeskService struct {
	repo repository.TodoDesk
}

func NewTodoDeskService(repo repository.TodoDesk) *TodoDeskService {
	return &TodoDeskService{repo: repo}
}

func (s *TodoDeskService) Create(userId int, desk todoListPrjct.DeskTable) (int, error) {
	return s.repo.Create(userId, desk)
}

func (s *TodoDeskService) GetAll(userId int) ([]todoListPrjct.DeskTable, error) {
	return s.repo.GetAll(userId)
}

func (s *TodoDeskService) GetDeskById(userId, id int) (todoListPrjct.DeskTable, error) {
	return s.repo.GetDeskById(userId, id)
}

func (s *TodoDeskService) Delete(userId, id int) error {
	return s.repo.Delete(userId, id)
}

func (s *TodoDeskService) Update(userId, deskId int, input todoListPrjct.UpdateDeskInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, deskId, input)
}