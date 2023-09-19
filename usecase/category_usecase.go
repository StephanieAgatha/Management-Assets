package usecase

import (
	"final-project-enigma-clean/exception"
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/repository"
	"final-project-enigma-clean/util/helper"
	"fmt"
)

type CategoryUsecase interface {
	CreateNew(payload model.Category) error
	FindById(id string) (model.Category, error)
	FindAll() ([]model.Category, error)
	Update(payload model.Category) error
	Delete(id string) error
}

type categoryUsecase struct {
	repo repository.CategoryRepository
}

// FindById implements CategoryUseCase.
func (c *categoryUsecase) FindById(id string) (model.Category, error) {
	category, err := c.repo.FindById(id)
	if err != nil {
		return model.Category{}, exception.BadRequestErr("category not found")
	}
	return category, nil

}

// CreateNew implements CategoryUseCase.
func (c *categoryUsecase) CreateNew(payload model.Category) error {
	if payload.Name == "" {
		return exception.BadRequestErr("name is required")
	}

	//commented for unit testing
	payload.Id = helper.GenerateUUID()
	err := c.repo.Save(payload)
	if err != nil {
		return fmt.Errorf("failed to create new category: %v", err)
	}
	return nil
}

// Delete implements CategoryUseCase.
func (c *categoryUsecase) Delete(id string) error {
	Category, err := c.FindById(id)
	if err != nil {
		return err
	}
	err = c.repo.Delete(Category.Id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %v", err)
	}
	return nil
}

// FindAll implements CategoryUseCase.
func (c *categoryUsecase) FindAll() ([]model.Category, error) {
	Category, err := c.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed to find all type asset: %v", err)
	}
	return Category, nil
}

// Update implements CategoryUseCase.
func (c *categoryUsecase) Update(payload model.Category) error {
	if payload.Name == "" {
		return exception.BadRequestErr("name is required")
	}
	_, err := c.FindById(payload.Id)
	if err != nil {
		return err
	}
	err = c.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update category: %v", err)
	}
	return nil
}

func NewCategoryUseCase(repo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		repo: repo,
	}
}
