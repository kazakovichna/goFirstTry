package repository

import (
	"errors"
	"github.com/kazakovichna/todoListPrjct"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		deskId int
		listId int
		item todoListPrjct.ItemTable
	}
	type mockBehavior func(args args, id int)

	testTable := []struct{
		name string
		mock mockBehavior
		args args
		id int
		want bool
	} {
		{
			name: "OK",
			args: args{
				listId: 1,
				item: todoListPrjct.ItemTable{
					UserId: 1,
					ItemName: "New",
					ItemDescription: "New",
					Done: false,
					ItemPosition: 1,
				},
			},
			id: 2,
			mock: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO itemTable").
					WithArgs(args.item.ItemName,
							 args.item.UserId,
							 args.item.ItemDescription,
							 args.item.Done,
							 args.item.ItemPosition,
						    ).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO listItemCompression").
					WithArgs(args.listId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			want: true,
			args: args{
				listId: 1,
				item: todoListPrjct.ItemTable{
					UserId: 1,
					ItemName: "",
					ItemDescription: "New",
					Done: false,
					ItemPosition: 1,
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))
				mock.ExpectQuery("INSERT INTO itemTable").
					WithArgs(args.item.ItemName,
						args.item.UserId,
						args.item.ItemDescription,
						args.item.Done,
						args.item.ItemPosition,
					).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
		},
		{
			name: "2nd Insert Error",
			want: true,
			args: args{
				listId: 1,
				item: todoListPrjct.ItemTable{
					UserId: 1,
					ItemName: "New",
					ItemDescription: "New",
					Done: false,
					ItemPosition: 1,
				},
			},
			mock: func(args args, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO itemTable").
					WithArgs(args.item.ItemName,
						args.item.UserId,
						args.item.ItemDescription,
						args.item.Done,
						args.item.ItemPosition,
					).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO listItemCompression").
					WithArgs(args.listId, id).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mock(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.listId, testCase.args.item)
			if testCase.want {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, got)
			}
		})
	}

}






