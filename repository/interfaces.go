// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/model"
)

type RepositoryInterface interface {
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (user *model.User)
	GetUserById(ctx context.Context, id int) (user *model.User)

	InsertUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
}
