package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kazakovichna/todoListPrjct"
	"strings"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(deskId int, list todoListPrjct.ListTable) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	createListQuery := fmt.Sprintf("INSERT INTO %s (ListName, ListPosition, Description) values ($1, $2, $3) RETURNING ListId", listTable)
	row := tx.QueryRow(createListQuery, list.ListName, list.ListPosition, list.Description)
	err = row.Scan(&listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createListDeskQuery := fmt.Sprintf("INSERT INTO %s (DeskId, ListId) values ($1, $2)", listDeskCompression)
	_, err = tx.Exec(createListDeskQuery, deskId, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (r *TodoListPostgres) GetAll(userId, deskId int) ([]todoListPrjct.ListTable, error) {
	var lists []todoListPrjct.ListTable
	query := fmt.Sprintf(`SELECT lt.listId, lt.listName, lt.description, lt.listPosition FROM %s lt INNER JOIN %s ld on lt.ListId = ld.ListId
									INNER JOIN %s ud on ud.deskId = ld.deskId WHERE ld.deskId = $1 AND ud.userId = $2`,
									listTable, listDeskCompression, userDeskCompression)

	if err := r.db.Select(&lists, query, deskId, userId); err != nil {
		return nil, err
	}

	return lists, nil
}

func (r *TodoListPostgres) GetById(userId, deskId, listId int) (todoListPrjct.ListTable, error) {
	var list todoListPrjct.ListTable
	query := fmt.Sprintf(`SELECT lt.listId, lt.listName, lt.description, lt.listposition FROM %s lt 
									INNER JOIN %s ld on lt.listId = ld.listId 
									INNER JOIN %s ud on ld.deskId = ud.deskId
									WHERE ud.userId = $1 AND ld.deskId = $2 AND lt.listId = $3`,
									listTable, listDeskCompression, userDeskCompression)
	if err := r.db.Get(&list, query, userId, deskId, listId); err != nil {
		return list, err
	}

	return list, nil
}

func (r *TodoListPostgres) Delete(userId, deskId, listId int) error {
	query := fmt.Sprintf(`DELETE FROM %s lt USING %s ld, %s ud
									WHERE lt.listId = ld.listId 
									AND ld.deskId = ud.deskId 
									AND ud.userId = $1 
									AND ld.deskId = $2 
									AND lt.listId = $3`,
									listTable, listDeskCompression, userDeskCompression)
	_, err := r.db.Exec(query, userId, deskId, listId)
	return err
}

func (r *TodoListPostgres) Update(userId, deskId, listId int, input todoListPrjct.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.ListName != nil {
		setValues = append(setValues, fmt.Sprintf("listname=$%d", argId))
		args = append(args, *input.ListName)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.ListPosition != nil {
		setValues = append(setValues, fmt.Sprintf("listposition=$%d", argId))
		args = append(args, *input.ListPosition)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	fmt.Printf("this is setQuery %v", setQuery)

	query := fmt.Sprintf(`UPDATE %s lt SET %s FROM %s ld, %s ud 
								WHERE lt.listId = ld.listId 
								AND ld.deskId = ud.deskId
								AND ud.userId = $%d AND ud.deskId = $%d AND lt.listId = $%d`,
								listTable, setQuery, listDeskCompression, userDeskCompression, argId, argId + 1, argId + 2)
	//fmt.Printf("this is query: == %s", query)
	args = append(args, userId, deskId, listId)

	_, err := r.db.Exec(query, args...)
	return err
}











