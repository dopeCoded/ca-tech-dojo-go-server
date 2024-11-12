package repository

import (
	"database/sql"

	"my-go-project/internal/model"
)

// UserRepository はユーザー関連のデータベース操作を定義するインターフェースです。
type UserRepository interface {
	CreateUser(name string) (*model.User, error)
	GetUserByID(id int64) (*model.User, error)
	UpdateUser(id int64, name string) error
}

// userRepository は UserRepository インターフェースを実装する構造体です。
type userRepository struct {
	db *sql.DB
}

// NewUserRepository は新しい UserRepository を生成します。
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

// CreateUser は新しいユーザーを users テーブルに挿入し、作成されたユーザーを返します。
func (r *userRepository) CreateUser(name string) (*model.User, error) {
	result, err := r.db.Exec(`
		INSERT INTO users (name, created_at, updated_at)
		VALUES (?, NOW(), NOW())
	`, name)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:   userID,
		Name: name,
		// CreatedAt と UpdatedAt はデータベースの DEFAULT 値を使用しているため、ここでは省略
	}

	return user, nil
}

// GetUserByID は指定されたユーザーIDに対応するユーザーを取得します。
func (r *userRepository) GetUserByID(id int64) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(`
		SELECT id, name, created_at, updated_at
		FROM users
		WHERE id = ?
	`, id).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser は指定されたユーザーIDの名前を更新します。
func (r *userRepository) UpdateUser(id int64, name string) error {
	_, err := r.db.Exec(`
		UPDATE users
		SET name = ?, updated_at = NOW()
		WHERE id = ?
	`, name, id)
	return err
}
