package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kazakovichna/todoListPrjct"
	"github.com/sirupsen/logrus"
	"strings"
)

type TodoDeskPostgres struct {
	db *sqlx.DB
}

func NewTodoDeskPostgres(db *sqlx.DB) *TodoDeskPostgres {
	return &TodoDeskPostgres{db: db}
}

func (r *TodoDeskPostgres) Create(userId int, desk todoListPrjct.DeskTable) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createDeskQuery := fmt.Sprintf("INSERT INTO %s (deskName, deskDescription) VALUES ($1, $2) RETURNING deskId", deskTable)
	row := tx.QueryRow(createDeskQuery, desk.DeskName, desk.DeskDescription)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUserDeskQuery := fmt.Sprintf("INSERT INTO %s (userId, deskId) VALUES ($1, $2)", userDeskCompression)
	_, err = tx.Exec(createUserDeskQuery, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TodoDeskPostgres) GetAll(userId int) ([]todoListPrjct.DeskTable, error) {
	var desks []todoListPrjct.DeskTable
	query := fmt.Sprintf("SELECT dl.deskId, dl.deskName, dl.deskDescription FROM %s dl INNER JOIN %s ud on dl.deskId = ud.deskId WHERE ud.userId = $1",
		deskTable, userDeskCompression)
	err := r.db.Select(&desks, query, userId)

	return desks, err
}

func (r *TodoDeskPostgres) GetDeskById(userId, id int) (todoListPrjct.DeskTable, error) {
	var desk todoListPrjct.DeskTable
	query := fmt.Sprintf("SELECT dl.deskId, dl.deskName, dl.deskDescription FROM %s dl INNER JOIN %s ud on dl.deskId = ud.deskId AND dl.deskId = $1 WHERE ud.userId = $2 ",
		deskTable, userDeskCompression)
	err := r.db.Get(&desk, query, id, userId)

	return desk, err
}

func (r *TodoDeskPostgres) Delete(userId, id int) error {
	query := fmt.Sprintf("DELETE FROM %s dt USING %s ud WHERE ud.userId = $1 AND ud.deskId = dt.deskId AND dt.deskId = $2",
		deskTable, userDeskCompression)
	_, err := r.db.Exec(query, userId, id)

	return err
}

func (r *TodoDeskPostgres) Update(userId, deskId int, input todoListPrjct.UpdateDeskInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.DeskName != nil {
		setValues = append(setValues, fmt.Sprintf("deskname=$%d", argId))
		args = append(args, *input.DeskName)
		argId++
	}

	if input.DeskDescription != nil {
		setValues = append(setValues, fmt.Sprintf("deskdescription=$%d", argId))
		args = append(args, *input.DeskDescription)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s dt SET %s FROM %s ud WHERE dt.deskId = ud.deskId AND ud.deskId = $%d AND ud.userId = $%d",
		deskTable, setQuery, userDeskCompression, argId, argId + 1)

	args = append(args, deskId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}




