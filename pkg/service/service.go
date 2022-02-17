package service

import (
	"github.com/kazakovichna/todoListPrjct"
	"github.com/kazakovichna/todoListPrjct/pkg/repository"
)

//go:generate go run github.com/golang/mock/mockgen -source=service.go -destination=mocks/mock.go

type Authorization interface {
	CreateUser(user todoListPrjct.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
	CreateSession(username, password string) (TokenResponse, error)
	RefreshTokenServices(refreshToken string) (TokenResponse, error)
}

type TodoDesk interface {
	Create(userId int, desk todoListPrjct.DeskTable) (int, error)
	GetAll(userId int) ([]todoListPrjct.DeskTable, error)
	GetDeskById(userId, id int) (todoListPrjct.DeskTable, error)
	Delete(userId, id int) error
	Update(userId, deskId int, input todoListPrjct.UpdateDeskInput) error
}

type TodoList interface {
	Create(userId, deskId int, list todoListPrjct.ListTable) (int, error)
	GetAll(userId, deskId int) ([]todoListPrjct.ListTable, error)
	GetById(userId, deskId, listId int) (todoListPrjct.ListTable, error)
	Delete(userId, deskId, listId int) error
	Update(userId, deskId, listId int, input todoListPrjct.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, input todoListPrjct.ItemTable) (int, error)
	GetAllItems(userId, deskId int) ([]todoListPrjct.AllDesksItems, error)
	GetItemById(userId, itemId int) (todoListPrjct.ItemTable, error)
	UpdateItem(userId, itemId int, input todoListPrjct.UpdateItemInput) error
	DeleteItem(userId, itemId int) error
}

type Service struct {
	Authorization
	TodoDesk
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoDesk: NewTodoDeskService(repos.TodoDesk),
		TodoList: NewTodoListService(repos.TodoList, repos.TodoDesk),
		TodoItem: NewTodoItemService(repos.TodoItem, repos.TodoList, repos.TodoDesk),
	}
}