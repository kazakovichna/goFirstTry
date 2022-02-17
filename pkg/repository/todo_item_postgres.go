package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kazakovichna/todoListPrjct"
	"strings"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, input todoListPrjct.ItemTable) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf(`INSERT INTO %s 
											(ItemName, UserId, 
											ItemDescription, 
											Done, ItemPosition)
											values ($1, $2, $3, $4, $5)
											RETURNING itemId`,
											itemTable)
	row := tx.QueryRow(createItemQuery, input.ItemName, input.UserId, input.ItemDescription, input.Done, input.ItemPosition)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	createItemListQuery := fmt.Sprintf(`INSERT INTO %s (listId, itemId) values ($1, $2)`, listItemCompression)
	_, err = tx.Exec(createItemListQuery, listId, itemId)

	return itemId, tx.Commit()
}

func (r *TodoItemPostgres) GetAllItems(userId, deskId int) ([]todoListPrjct.AllDesksItems, error) {
	var items []todoListPrjct.AllDesksItems

	getAllItemsQuery := fmt.Sprintf(`SELECT it.itemId, it.userId, it.itemName, it.itemDescription, it.Done, it.itemPosition, lt.listPosition FROM %s it 
											INNER JOIN %s li on it.itemId = li.itemId 
											INNER JOIN %s lt on lt.listId = li.listId
											INNER JOIN %s ld on li.listId = ld.listId 
											INNER JOIN %s ud on ld.deskId = ud.deskId
											WHERE ld.deskId = $1 AND ud.userId = $2`,
											itemTable, listItemCompression, listTable, listDeskCompression, userDeskCompression)

	fmt.Printf("This is query: == %v", getAllItemsQuery)
	if err := r.db.Select(&items, getAllItemsQuery, deskId, userId); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *TodoItemPostgres) GetItemById(userId, itemId int) (todoListPrjct.ItemTable, error) {
	var item todoListPrjct.ItemTable

	getItemByIdQuery := fmt.Sprintf(`SELECT it.itemId, it.userId, it.itemName, it.itemDescription, it.Done, it.itemPosition FROM %s it 
											INNER JOIN %s li on it.itemId = li.itemId 
											INNER JOIN %s ld on li.listId = ld.listId 
											INNER JOIN %s ud on ld.deskId = ud.deskId
											WHERE it.itemId = $1 AND ud.userId = $2`,
		itemTable, listItemCompression, listDeskCompression, userDeskCompression)

	if err := r.db.Get(&item, getItemByIdQuery, itemId, userId); err != nil {
		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) UpdateItem(userId, itemId int, input todoListPrjct.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.ItemName != nil {
		setValues = append(setValues, fmt.Sprintf("itemname=$%d", argId))
		args = append(args, *input.ItemName)
		argId++
	}

	if input.ItemDescription != nil {
		setValues = append(setValues, fmt.Sprintf("itemdescription=$%d", argId))
		args = append(args, *input.ItemDescription)
		argId++
	}

	if input.ItemPosition != nil {
		setValues = append(setValues, fmt.Sprintf("itemposition=$%d", argId))
		args = append(args, *input.ItemPosition)
		argId++
	}

	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateItemQuery := fmt.Sprintf(`UPDATE %s it SET %s FROM %s li, %s ld, %s ud
											WHERE it.itemId = li.itemId AND li.listId = ld.listId
											AND ld.deskId = ud.deskId AND ud.userId = $%d AND it.itemId = $%d`,
											itemTable, setQuery, listItemCompression, listDeskCompression, userDeskCompression, argId, argId + 1)
	args = append(args, userId, itemId)

	_, err := r.db.Exec(updateItemQuery, args...)
	return err
}

func (r *TodoItemPostgres) DeleteItem(userId, itemId int) error {
	deleteItemQuery := fmt.Sprintf(`DELETE FROM %s it USING %s li, %s ld, %s ud
											WHERE it.itemId = li.itemId AND li.listId = ld.listId
											AND ld.deskId = ud.deskId AND ud.userId = $1 AND it.itemId = $2`,
											itemTable, listItemCompression, listDeskCompression, userDeskCompression)
	_, err := r.db.Exec(deleteItemQuery, userId, itemId)
	return err
}