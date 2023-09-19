package usecase

import (
	"final-project-enigma-clean/exception"
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/model/dto"
	"final-project-enigma-clean/repository"
	"final-project-enigma-clean/util/helper"
	"fmt"
	"time"
)

type AssetUsecase interface {
	Create(payload model.AssetRequest) error
	FindAll() ([]model.Asset, error)
	FindById(id string) (model.Asset, error)
	Update(payload model.AssetRequest) error
	UpdateAvailable(id string, amount int) error
	Delete(id string) error
	FindByName(name string) ([]model.Asset, error)
	Paging(payload dto.PageRequest) ([]model.Asset, dto.Paging, error)
}

type assetUsecase struct {
	repo repository.AssetRepository
	//get category usecase
	categoryUc CategoryUsecase
	//get asset type usecase
	typeAssetUC TypeAssetUseCase
}

func (a *assetUsecase) Paging(payload dto.PageRequest) ([]model.Asset, dto.Paging, error) {
	if payload.Page <= 0 {
		payload.Page = 1
	}

	return a.repo.Paging(payload)
}

// UpdateAmount implements AssetUsecase.
func (a *assetUsecase) UpdateAvailable(id string, amount int) error {

	asset, err := a.FindById(id)
	if err != nil {
		return err
	}

	asset.Available -= amount
	err = a.repo.UpdateAvailable(id, asset.Available)
	if err != nil {
		return fmt.Errorf("failed update amount, %s", err)
	}

	return nil
}

// FindByName implements AssetUsecase.
func (a *assetUsecase) FindByName(name string) ([]model.Asset, error) {
	if name == "" {
		return nil, exception.BadRequestErr("name cannot empty")
	}

	assets, err := a.repo.FindByName(name)
	if err != nil {
		return nil, fmt.Errorf("failed get assets, %s", err)
	}
	return assets, nil
}

// Create implements AssetUsecase.
func (a *assetUsecase) Create(payload model.AssetRequest) error {
	if payload.Name == "" {
		return exception.BadRequestErr("name cannot empty")
	}
	if payload.AssetTypeId == "" || payload.CategoryId == "" {
		return exception.BadRequestErr("asset type id or category id cannot empty")
	}
	if payload.Total < 0 {
		return exception.BadRequestErr("Total cannot negative number")
	}
	if payload.Status == "" {
		return exception.BadRequestErr("status cannot empty")
	}

	//implement asset type find by id
	_, err := a.typeAssetUC.FindById(payload.AssetTypeId)
	if err != nil {
		return err
	}

	//implement category find by id
	_, err = a.categoryUc.FindById(payload.CategoryId)
	if err != nil {
		return err
	}

	//commented id and entrydate for unit testing
	payload.Id = helper.GenerateUUID()
	payload.EntryDate = time.Now()
	payload.Available = payload.Total
	err = a.repo.Save(payload)
	if err != nil {
		return fmt.Errorf("failed save asset %s", err)
	}

	return nil
}

// Delete implements AssetUsecase.
func (a *assetUsecase) Delete(id string) error {
	//find assert first
	_, err := a.FindById(id)
	if err != nil {
		return err
	}

	err = a.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete asset, %s", err)
	}

	return nil
}

// FindAll implements AssetUsecase.
func (a *assetUsecase) FindAll() ([]model.Asset, error) {
	assets, err := a.repo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("failed get assets, %s", err)
	}
	return assets, nil
}

// FindById implements AssetUsecase.
func (a *assetUsecase) FindById(id string) (model.Asset, error) {
	asset, err := a.repo.FindById(id)
	if err != nil {
		return model.Asset{}, exception.BadRequestErr(fmt.Sprintf("asset by id:%s cannot found, err:%s", id, err))
	}

	return asset, nil
}

// Update implements AssetUsecase.
func (a *assetUsecase) Update(payload model.AssetRequest) error {
	if payload.Name == "" {
		return exception.BadRequestErr("name cannot empty")
	}
	if payload.AssetTypeId == "" || payload.CategoryId == "" {
		return exception.BadRequestErr("asset type id or category id cannot empty")
	}
	if payload.Total < 0 {
		return exception.BadRequestErr("amoun cannot negative number")
	}
	if payload.Status == "" {
		return exception.BadRequestErr("status cannot empty")
	}

	//implement asset type find by id
	_, err := a.typeAssetUC.FindById(payload.AssetTypeId)
	if err != nil {
		return err
	}

	//implement category find by id
	_, err = a.categoryUc.FindById(payload.CategoryId)
	if err != nil {
		return err
	}

	asset, err := a.FindById(payload.Id)
	if err != nil {
		return err
	}

	//calculation for update available
	payload.Available = (payload.Total - asset.Total) + asset.Available

	err = a.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed update asset %s", err)
	}

	return nil
}

func NewAssetUsecase(assetRepo repository.AssetRepository, typeAssetUC TypeAssetUseCase, categoryUC CategoryUsecase) AssetUsecase {
	return &assetUsecase{
		repo:        assetRepo,
		categoryUc:  categoryUC,
		typeAssetUC: typeAssetUC,
	}
}
