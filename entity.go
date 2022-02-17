package todoListPrjct

import (
	"errors"
	"time"
)

type UserRefreshToken struct {
	Id int `json:"id" db:"userid"`
	Name string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	//UserEmail string `json:"user_email" binding:"required"`
	PasswordHash string `json:"password" binding:"required"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt int `json:"expires_at"`
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Session struct {
	RefreshToken string    `json:"refreshToken" db:"refreshtoken"`
	ExpiresAt    time.Time `json:"expiresAt" db:"expiresat"`
}

type User struct {
	Id int `json:"id" db:"userid"`
	Name string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt int `json:"expires_at"`
}

type UserDeskCompression struct {
	UserDeskId int `json:"user_desk_id"`
	UserId int `json:"user_id"`
	DeskId int `json:"desk_id"`
}

type DeskTable struct {
	DeskID int 	`json:"desk_id" db:"deskid"`
	DeskName string `json:"desk_name" binding:"required" db:"deskname"`
	DeskDescription string `json:"desk_description" db:"deskdescription"`
}

type ListDeskCompression struct {
	ListDeskId int `json:"list_desk_id"`
	DeskId int `json:"desk_id"`
	ListId int `json:"list_id"`
}

type ListTable struct {
	ListId int `json:"list_id" db:"listid"`
	ListName string `json:"list_name" db:"listname" binding:"required"`
	Description string `json:"description" db:"description" binding:"required"`
	ListPosition int `json:"list_position" db:"listposition" binding:"required"`
}

type ListItemCompression struct {
	ListItemId int `json:"list_item_id"`
	ListId int `json:"list_id"`
	ItemId int `json:"item_id"`
}

type ItemTable struct {
	ItemId int `json:"item_id"`
	UserId int `json:"user_id"`
	ItemName string `json:"item_name"`
	ItemDescription string `json:"item_description"`
	Done bool `json:"done"`
	ItemPosition int `json:"item_position"`
}

type AllDesksItems struct {
	ItemId int `json:"item_id"`
	UserId int `json:"user_id"`
	ItemName string `json:"item_name"`
	ItemDescription string `json:"item_description"`
	Done bool `json:"done"`
	ItemPosition int `json:"item_position"`
	ListPosition int `json:"list_position"`
}

type UpdateDeskInput struct {
	DeskName *string `json:"desk_name"`
	DeskDescription *string `json:"desk_description"`
}

func (i UpdateDeskInput) Validate() error {
	if i.DeskName == nil && i.DeskDescription == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateListInput struct {
	ListName *string `json:"list_name"`
	Description *string `json:"description"`
	ListPosition *int `json:"list_position"`
}

func (i UpdateListInput) ValidateList() error {
	if i.ListName == nil && i.Description == nil && i.ListPosition == nil {
		return errors.New("update structure has no values")
	}

	return nil
}

type UpdateItemInput struct {
	ItemName *string `json:"item_name"`
	ItemDescription *string `json:"item_description"`
	Done *bool `json:"done"`
	ItemPosition *int `json:"item_position"`
}

func (i UpdateItemInput) ValidateItem() error {
	if i.ItemName == nil && i.ItemDescription == nil &&
		i.Done == nil && i.ItemPosition == nil {
		return errors.New("update item structure has no values")
	}

	return nil
}