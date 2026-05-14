package admin

import "context"

type RoleService struct {
	roleRepo RoleRepository
}

func NewRoleService(roleRepo RoleRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
	}
}

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type UpdateRoleRequest struct {
	Name string `json:"name"`
}

func (s *RoleService) Create(ctx context.Context, req CreateRoleRequest) (*SystemRole, error) {
	role := &SystemRole{
		Name: req.Name,
	}

	if err := s.roleRepo.Create(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) GetByID(ctx context.Context, id int64) (*SystemRole, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) List(ctx context.Context, limit, offset int) ([]SystemRole, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	roles, err := s.roleRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (s *RoleService) Update(ctx context.Context, id int64, req UpdateRoleRequest) (*SystemRole, error) {
	role, err := s.roleRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	role.Name = req.Name

	if err := s.roleRepo.Update(ctx, role); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *RoleService) Delete(ctx context.Context, id int64) error {
	if err := s.roleRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
