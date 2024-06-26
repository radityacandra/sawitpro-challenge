package repository

import (
	"context"
	"time"

	"github.com/SawitProRecruitment/UserService/model"
)

func (r *Repository) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) *model.User {
	var user model.User
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE phone_number = $1", phoneNumber[1:]).
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

func (r *Repository) GetUserById(ctx context.Context, id int) *model.User {
	var user model.User
	err := r.Db.QueryRowContext(ctx, "SELECT id, full_name, phone_number, password FROM users WHERE id = $1", id).
		Scan(&user.Id, &user.FullName, &user.PhoneNumber, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

func (r *Repository) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := r.Db.ExecContext(ctx, "UPDATE users SET full_name = $1, phone_number = $2, updated_by = $3 WHERE id = $3", user.FullName, user.PhoneNumber, user.Id)
	if err != nil {
		return err
	}

	return nil
}
