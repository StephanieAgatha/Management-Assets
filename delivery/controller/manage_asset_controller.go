package controller

import (
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/model/dto"
	"final-project-enigma-clean/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ManageAssetController struct {
	manageAssetUC usecase.ManageAssetUsecase
	g             *gin.RouterGroup
}

// show assets handler
func (m *ManageAssetController) ShowAllAssetHandler(c *gin.Context) {

	mAssets, err := m.manageAssetUC.ShowAllAsset()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{"Message": "Success", "Data": mAssets})
}

//create a new asset handler

func (m *ManageAssetController) CreateNewAssetHandler(c *gin.Context) {

	var manageAssetReq dto.ManageAssetRequest
	if err := c.ShouldBindJSON(&manageAssetReq); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Bad JSON Format", "error": err.Error()})
		return
	}

	if err := m.manageAssetUC.CreateTransaction(manageAssetReq); err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"Message": "Success"})

}

func (m *ManageAssetController) FindByIdTransaction(c *gin.Context) {
	id := c.Param("id")

	detailAssets, err := m.manageAssetUC.FindByTransactionID(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"Data": detailAssets})
}
func (m *ManageAssetController) FindByName(c *gin.Context) {

	var staff model.Staff
	err := c.ShouldBindJSON(&staff)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": "Bad JSON Format", "error": err.Error()})
		return
	}

	result, err := m.manageAssetUC.FindTransactionByName(staff.Name)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"Data": result})
}

// download handler
func (m *ManageAssetController) DownloadAssetsHandler(c *gin.Context) {
	//set header
	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", `attachment; filename="data-assets.csv"`)
	csvData, err := m.manageAssetUC.DownloadAssets()
	if err != nil {
		c.Error(err)
		return
	}
	c.Data(http.StatusOK, "text/csv", csvData)
}
func (m *ManageAssetController) Route() {
	m.g.GET("/manage-assets/show-all", m.ShowAllAssetHandler)
	m.g.POST("/manage-assets/create-new", m.CreateNewAssetHandler)
	m.g.GET("/manage-assets/find/:id", m.FindByIdTransaction)
	m.g.POST("/manage-assets/find-asset", m.FindByName)
	m.g.GET("/manage-assets/download/list-assets", m.DownloadAssetsHandler)
}

func NewManageAssetController(maUC usecase.ManageAssetUsecase, g *gin.RouterGroup) *ManageAssetController {
	return &ManageAssetController{
		manageAssetUC: maUC,
		g:             g,
	}
}
