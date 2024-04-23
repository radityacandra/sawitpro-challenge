package repository

import (
	"context"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
)

func (r *Repository) GetTestById(ctx context.Context, input GetTestByIdInput) (output GetTestByIdOutput, err error) {
	err = r.Db.QueryRowContext(ctx, "SELECT name FROM test WHERE id = $1", input.Id).Scan(&output.Name)
	if err != nil {
		return
	}
	return
}

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) *model.User {
	var user model.User
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1", phoneNumber).
		Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

func (r *Repository) InsertUser(ctx context.Context, user *model.User) (*model.User, error) {
	var id int

	err := r.Db.QueryRowContext(ctx,
		"INSERT INTO users (full_name, phone_number, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		user.FullName, user.PhoneNumber, user.Password, time.Now()).Scan(&id)
	if err != nil {
		return nil, err
	}

	user.Id = id

	return user, nil
}
