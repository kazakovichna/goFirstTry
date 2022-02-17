package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/kazakovichna/todoListPrjct"
)

type Authorization interface {
	CreateUser(user todoListPrjct.User) (int, error)
	GetUser(username, password string) (todoListPrjct.User, error)
	SetSessions(username, refreshToken string, expiresAt int) error
	GetUserByRefreshToken(refreshToken string) (todoListPrjct.UserRefreshToken, error)
}

type TodoDesk interface {
	Create(userId int, desk todoListPrjct.DeskTable) (int, error)
	GetAll(userId int) ([]todoListPrjct.DeskTable, error)
	GetDeskById(userId, id int) (todoListPrjct.DeskTable, error)
	Delete(userId, id int) error
	Update(userId, deskId int, input todoListPrjct.UpdateDeskInput) error
}

type TodoList interface {
	Create(deskId int, list todoListPrjct.ListTable) (int, error)
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

type Repository struct {
	Authorization
	TodoDesk
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoDesk: NewTodoDeskPostgres(db),
		TodoList: NewTodoListPostgres(db),
		TodoItem: NewTodoItemPostgres(db),
	}
}