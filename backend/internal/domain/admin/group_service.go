package admin

import "context"

type GroupService struct {
	groupRepo GroupRepository
}

func NewGroupService(groupRepo GroupRepository) *GroupService {
	return &GroupService{
		groupRepo: groupRepo,
	}
}

type CreateGroupRequest struct {
	Name        string  `json:"name"`
	FrontpageID *int64  `json:"frontpage_id"`
}

type UpdateGroupRequest struct {
	Name        string  `json:"name"`
	FrontpageID *int64  `json:"frontpage_id"`
}

func (s *GroupService) Create(ctx context.Context, req CreateGroupRequest) (*SystemGroup, error) {
	group := &SystemGroup{
		Name:        req.Name,
		FrontpageID: req.FrontpageID,
	}

	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupService) GetByID(ctx context.Context, id int64) (*SystemGroup, error) {
	group, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupService) List(ctx context.Context, limit, offset int) ([]SystemGroup, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	groups, err := s.groupRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (s *GroupService) Update(ctx context.Context, id int64, req UpdateGroupRequest) (*SystemGroup, error) {
	group, err := s.groupRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	group.Name = req.Name
	group.FrontpageID = req.FrontpageID

	if err := s.groupRepo.Update(ctx, group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *GroupService) Delete(ctx context.Context, id int64) error {
	if err := s.groupRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
