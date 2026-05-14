package admin

import "context"

type ProgramService struct {
	programRepo ProgramRepository
}

func NewProgramService(programRepo ProgramRepository) *ProgramService {
	return &ProgramService{
		programRepo: programRepo,
	}
}

type CreateProgramRequest struct {
	Name        string `json:"name"`
	Controller  string `json:"controller"`
	Description string `json:"description"`
}

type UpdateProgramRequest struct {
	Name        string `json:"name"`
	Controller  string `json:"controller"`
	Description string `json:"description"`
}

func (s *ProgramService) Create(ctx context.Context, req CreateProgramRequest) (*SystemProgram, error) {
	program := &SystemProgram{
		Name:        req.Name,
		Controller:  req.Controller,
		Description: req.Description,
	}

	if err := s.programRepo.Create(ctx, program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) GetByID(ctx context.Context, id int64) (*SystemProgram, error) {
	program, err := s.programRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) List(ctx context.Context, limit, offset int) ([]SystemProgram, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	programs, err := s.programRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return programs, nil
}

func (s *ProgramService) Update(ctx context.Context, id int64, req UpdateProgramRequest) (*SystemProgram, error) {
	program, err := s.programRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	program.Name = req.Name
	program.Controller = req.Controller
	program.Description = req.Description

	if err := s.programRepo.Update(ctx, program); err != nil {
		return nil, err
	}

	return program, nil
}

func (s *ProgramService) Delete(ctx context.Context, id int64) error {
	if err := s.programRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
