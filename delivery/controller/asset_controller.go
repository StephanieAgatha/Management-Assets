package controller

import (
	"final-project-enigma-clean/delivery/middleware"
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/model/dto"
	"final-project-enigma-clean/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetController struct {
	usecase usecase.AssetUsecase
	rg      *gin.RouterGroup
}

func (a *AssetController) createAssetHandler(c *gin.Context) {

	var assetRequest model.AssetRequest
	err := c.ShouldBindJSON(&assetRequest)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	err = a.usecase.Create(assetRequest)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, gin.H{"status": "OK", "message": "successfully created Asset", "asset": assetRequest})
	// c.JSON(201, assetRequest)
}

func (a *AssetController) ListAssetHandler(c *gin.Context) {
	name := c.Query("name")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "5"))

	if name != "" {
		assets, err := a.usecase.FindByName(name)
		if err != nil {
			c.Error(err)
			return
		}

		c.JSON(200, gin.H{
			"status": "OK",
			"assets": assets,
		})
		return
	}

	assets, paging, err := a.usecase.Paging(dto.PageRequest{
		Page: page,
		Size: size,
	})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"status": "OK",
		"assets": assets,
		"paging": paging,
	})
}

func (a *AssetController) findByIdHandler(c *gin.Context) {

	id := c.Param("id")

	asset, err := a.usecase.FindById(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{
		"status": "OK",
		"assets": asset,
	})
}

func (a *AssetController) updateHandler(c *gin.Context) {

	var assetRequest model.AssetRequest
	err := c.ShouldBindJSON(&assetRequest)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"status": "Error", "message": err.Error()})
		return
	}

	err = a.usecase.Update(assetRequest)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, gin.H{"status": "OK", "message": "successfully Update Asset"})
}

func (a *AssetController) deleteHandler(c *gin.Context) {

	id := c.Param("id")

	err := a.usecase.Delete(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, gin.H{"status": "OK", "message": "successfully delete asset"})
}

func (a *AssetController) Route() {
	a.rg.POST("/assets", middleware.AuthMiddleware(), a.createAssetHandler)
	a.rg.GET("/assets", middleware.AuthMiddleware(), a.ListAssetHandler)
	a.rg.GET("/assets/:id", middleware.AuthMiddleware(), a.findByIdHandler)
	a.rg.PUT("/assets", middleware.AuthMiddleware(), a.updateHandler)
	a.rg.DELETE("/assets/:id", middleware.AuthMiddleware(), a.deleteHandler)
}

func NewAssetController(usecase usecase.AssetUsecase, rg *gin.RouterGroup) *AssetController {
	return &AssetController{
		usecase: usecase,
		rg:      rg,
	}
}
