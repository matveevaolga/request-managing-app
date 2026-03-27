package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/matveevaolga/request-managing-app/internal/domain"
	domainrepo "github.com/matveevaolga/request-managing-app/internal/domain/repository"
)

func seedUsers(ctx context.Context, repo domainrepo.UserRepository) error {
	users, err := loadUsers()
	if err != nil {
		return err
	}
	for _, u := range users {
		if err := createUser(ctx, repo, u); err != nil {
			return err
		}
	}
	return nil
}

func loadUsers() ([]userSeed, error) {
	data, err := seedFS.ReadFile("data/users.json")
	if err != nil {
		return nil, err
	}
	var users []userSeed
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func createUser(ctx context.Context, repo domainrepo.UserRepository, u userSeed) error {
	exists, err := repo.GetByUsername(ctx, u.Username)
	if err != nil && err != domain.ErrUserNotFound {
		return err
	}
	if exists != nil {
		slog.Info("User already exists", "username", u.Username)
		return nil
	}
	role := domain.RoleUser
	if u.Role == "ADMIN" {
		role = domain.RoleAdmin
	}
	user := &domain.User{
		Username:  u.Username,
		Password:  u.Password,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := repo.Create(ctx, user); err != nil {
		return err
	}
	slog.Info("User created", "username", user.Username, "role", user.Role)
	return nil
}
