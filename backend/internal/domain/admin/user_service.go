package admin

import (
	"context"
	"fmt"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

type CreateUserRequest struct {
	Login        string  `json:"login"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     string  `json:"password"`
	SystemUnitID *int64  `json:"system_unit_id"`
	FrontpageID  *int64  `json:"frontpage_id"`
	EquipeID     *int64  `json:"equipe_id"`
	Phone        *string `json:"phone"`
	Address      *string `json:"address"`
}

type UpdateUserRequest struct {
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Password     *string `json:"password"`
	SystemUnitID *int64  `json:"system_unit_id"`
	Active       *bool   `json:"active"`
	FrontpageID  *int64  `json:"frontpage_id"`
	EquipeID     *int64  `json:"equipe_id"`
	Phone        *string `json:"phone"`
	Address      *string `json:"address"`
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*SystemUser, error) {
	// Verificar se login já existe
	_, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err == nil {
		return nil, fmt.Errorf("login already exists")
	}

	// Hash da senha
	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &SystemUser{
		Login:        req.Login,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: passwordHash,
		SystemUnitID: req.SystemUnitID,
		Active:       true,
		FrontpageID:  req.FrontpageID,
		EquipeID:     req.EquipeID,
		Phone:        req.Phone,
		Address:      req.Address,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Remover hash da senha da resposta
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*SystemUser, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Remover hash da senha da resposta
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) List(ctx context.Context, limit, offset int) ([]SystemUser, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	users, err := s.userRepo.List(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Remover hashes de senha das respostas
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

func (s *UserService) Update(ctx context.Context, id int64, req UpdateUserRequest) (*SystemUser, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Atualizar campos
	user.Name = req.Name
	user.Email = req.Email
	user.SystemUnitID = req.SystemUnitID
	user.FrontpageID = req.FrontpageID
	user.EquipeID = req.EquipeID
	user.Phone = req.Phone
	user.Address = req.Address

	if req.Active != nil {
		user.Active = *req.Active
	}

	// Se senha foi fornecida, atualizar
	if req.Password != nil {
		passwordHash, err := HashPassword(*req.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = passwordHash
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Remover hash da senha da resposta
	user.PasswordHash = ""
	return user, nil
}

func (s *UserService) Delete(ctx context.Context, id int64) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
