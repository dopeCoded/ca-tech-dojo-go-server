package service

import (
	"my-go-project/internal/model"
	"my-go-project/internal/repository"
)

// UserService はユーザー関連のビジネスロジックを定義するインターフェースです。
type UserService interface {
	CreateUser(name string) (*model.User, error)
	GetUser(id int64) (*model.User, error)
	UpdateUser(id int64, name string) error
}

// userService は UserService インターフェースを実装する構造体です。
type userService struct {
	repo repository.UserRepository
}

// NewUserService は新しい UserService を生成します。
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

// CreateUser は新しいユーザーを作成します。
func (s *userService) CreateUser(name string) (*model.User, error) {
	// 必要に応じて追加のビジネスロジックをここに記述できます。
	return s.repo.CreateUser(name)
}

// GetUser は指定されたユーザーIDに対応するユーザー情報を取得します。
func (s *userService) GetUser(id int64) (*model.User, error) {
	return s.repo.GetUserByID(id)
}

// UpdateUser は指定されたユーザーIDの名前を更新します。
func (s *userService) UpdateUser(id int64, name string) error {
	return s.repo.UpdateUser(id, name)
}
