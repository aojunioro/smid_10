package admin

import "context"

type UnitService struct {
	unitRepo UnitRepository
}

func NewUnitService(unitRepo UnitRepository) *UnitService {
	return &UnitService{
		unitRepo: unitRepo,
	}
}

type CreateUnitRequest struct {
	Name     string  `json:"name"`
	ParentID *int64  `json:"parent_id"`
}

type UpdateUnitRequest struct {
	Name     string  `json:"name"`
	ParentID *int64  `json:"parent_id"`
}

func (s *UnitService) Create(ctx context.Context, req CreateUnitRequest) (*SystemUnit, error) {
	unit := &SystemUnit{
		Name:     req.Name,
		ParentID: req.ParentID,
	}

	if err := s.unitRepo.Create(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitService) GetByID(ctx context.Context, id int64) (*SystemUnit, error) {
	unit, err := s.unitRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitService) List(ctx context.Context, limit, offset int) ([]SystemUnit, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	units, err := s.unitRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return units, nil
}

func (s *UnitService) Update(ctx context.Context, id int64, req UpdateUnitRequest) (*SystemUnit, error) {
	unit, err := s.unitRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	unit.Name = req.Name
	unit.ParentID = req.ParentID

	if err := s.unitRepo.Update(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

func (s *UnitService) Delete(ctx context.Context, id int64) error {
	if err := s.unitRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
