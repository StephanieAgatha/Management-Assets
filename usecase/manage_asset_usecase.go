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

type ManageAssetUsecase interface {
	CreateTransaction(payload dto.ManageAssetRequest) error
	ShowAllAsset() ([]model.ManageAsset, error)
	FindByTransactionID(id string) ([]model.ManageAsset, error)
	FindTransactionByName(name string) ([]model.ManageAsset, error)
	DownloadAssets() ([]byte, error)
}

type manageAssetUsecase struct {
	repo    repository.ManageAssetRepository
	staffUC StaffUseCase
	assetUC AssetUsecase
}

// FindTransactionByName implements ManageAssetUsecase.
func (m *manageAssetUsecase) FindTransactionByName(name string) ([]model.ManageAsset, error) {
	if name == "" {
		return nil, exception.BadRequestErr("name cannot empty")
	}

	transactions, transactionDetails, err := m.repo.FindByNameTransaction(name)
	if err != nil {
		return nil, err
	}

	detailMap := make(map[string][]model.ManageDetailAsset)

	// Kelompokkan detail transaksi berdasarkan Id ManageAsset
	for _, detail := range transactionDetails {
		detailMap[detail.ManageAssetId] = append(detailMap[detail.ManageAssetId], detail)
	}

	// Inisialisasi slice datas
	datas := make([]model.ManageAsset, 0)

	// Iterasi melalui transaksi untuk membangun datas
	for _, transaction := range transactions {
		if details, ok := detailMap[transaction.Id]; ok {
			transaction.Detail = details
			datas = append(datas, transaction)
		}
	}
	return datas, nil
}

// CreateTransaction implements ManageAssetUsecase.
func (m *manageAssetUsecase) CreateTransaction(payload dto.ManageAssetRequest) error {
	if payload.NikStaff == "" {
		return exception.BadRequestErr("nik staff cannot empty")
	}

	payload.Id = helper.GenerateUUID()
	var newManageDetail []dto.ManageAssetDetailRequest
	//looping for validation request detail
	for _, detail := range payload.ManageAssetDetailReq {
		if detail.IdAsset == "" {
			return exception.BadRequestErr("id asset cannot empty")
		}

		if detail.Status == "" {
			return exception.BadRequestErr("status cannot empty")
		}

		if detail.TotalItem < 0 {
			return exception.BadRequestErr("total item must equal than 0")
		}

		asset, err := m.assetUC.FindById(detail.IdAsset)
		if err != nil {
			return err
		}
		//valdiation asset amount available or not
		if asset.Available < detail.TotalItem {
			return exception.BadRequestErr("Barang tidak cukup")
		}
		detail.Id = helper.GenerateUUID()
		newManageDetail = append(newManageDetail, detail)
	}
	//validate nikstaff
	_, err := m.staffUC.FindById(payload.NikStaff)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	//reassign value

	payload.ManageAssetDetailReq = newManageDetail

	//comment time.now if you want to run unit testing
	payload.SubmisstionDate = time.Now()
	payload.ReturnDate = payload.SubmisstionDate.AddDate(0, 0, payload.Duration)
	err = m.repo.CreateTransaction(payload)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	//update amount of asset when success
	for _, detail := range payload.ManageAssetDetailReq {
		err = m.assetUC.UpdateAvailable(detail.IdAsset, detail.TotalItem)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *manageAssetUsecase) ShowAllAsset() ([]model.ManageAsset, error) {
	//TODO implement me
	return m.repo.FindAllTransaction()
}

func (m *manageAssetUsecase) FindByTransactionID(id string) ([]model.ManageAsset, error) {
	//TODO implement me
	if id == "" {
		return nil, exception.BadRequestErr("ID is required")
	}

	transactions, transactionDetails, err := m.repo.FindAllByTransId(id)
	if err != nil {
		return nil, err
	}

	detailMap := make(map[string][]model.ManageDetailAsset)

	// Kelompokkan detail transaksi berdasarkan Id ManageAsset
	for _, detail := range transactionDetails {
		detailMap[detail.ManageAssetId] = append(detailMap[detail.ManageAssetId], detail)
	}

	// Inisialisasi slice datas
	datas := make([]model.ManageAsset, 0)

	// Iterasi melalui transaksi untuk membangun datas
	for _, transaction := range transactions {
		if details, ok := detailMap[transaction.Id]; ok {
			transaction.Detail = details
			datas = append(datas, transaction)
		}
	}
	return datas, nil
}

func (m *manageAssetUsecase) DownloadAssets() ([]byte, error) {
	//TODO implement me
	assets, err := m.repo.FindAllTransaction()
	if err != nil {
		return nil, fmt.Errorf("Failed to find assets %v", err.Error())
	}

	//convert data to csv
	csvData, err := helper.ConvertToCSVForAssets(assets)
	return csvData, nil
}

func NewManageAssetUsecase(repo repository.ManageAssetRepository, staffUC StaffUseCase, assetUC AssetUsecase) ManageAssetUsecase {
	return &manageAssetUsecase{
		repo:    repo,
		staffUC: staffUC,
		assetUC: assetUC,
	}
}
