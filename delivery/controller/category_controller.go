package controller

import (
	"final-project-enigma-clean/delivery/middleware"
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/usecase"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryUC usecase.CategoryUsecase
	rg          *gin.RouterGroup
}

func (cc *CategoryController) createHandlerCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	// category.Id = helper.GenerateUUID()
	err := cc.categoryUC.CreateNew(category)
	if err != nil {
		c.Error(err)
		return
	}
	response := gin.H{
		"message": "successfully created category",
	}
	c.JSON(201, response)
}
func (cc *CategoryController) listHandlerCategory(c *gin.Context) {

	category, err := cc.categoryUC.FindAll()
	if err != nil {
		c.Error(err)
		return
	}
	response := gin.H{
		"message": "successfully get category",
		"data":    category,
	}
	c.JSON(200, response)
}
func (cc *CategoryController) getByIdteHandlerCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := cc.categoryUC.FindById(id)
	if err != nil {
		c.Error(err)
		return
	}
	response := gin.H{
		"message": "successfully get by id category",
		"data":    category,
	}
	c.JSON(200, response)
}


func (cc *CategoryController) updateHandlerCategory(c *gin.Context) {
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	err := cc.categoryUC.Update(category)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, gin.H{
		"message": "successfully update category",
	})
}
func (cc *CategoryController) deleteHandlerCategory(c *gin.Context) {
	id := c.Param("id")
	if err := cc.categoryUC.Delete(id); err != nil {
		c.Error(err)
		return
	}
	message := fmt.Sprintf("successfulyy delete category with id %s", id)
	c.JSON(200, gin.H{
		"message": message,
	})
}
func (cc *CategoryController) Route() {
	cc.rg.POST("/categories", middleware.AuthMiddleware(), cc.createHandlerCategory)
	cc.rg.GET("/categories", middleware.AuthMiddleware(), cc.listHandlerCategory)
	cc.rg.GET("/categories/:id", middleware.AuthMiddleware(), cc.getByIdteHandlerCategory)
	cc.rg.PUT("/categories", middleware.AuthMiddleware(), cc.updateHandlerCategory)
	cc.rg.DELETE("/categories/:id", middleware.AuthMiddleware(), cc.deleteHandlerCategory)
}

func NewCategoryController(categoryUC usecase.CategoryUsecase, rg *gin.RouterGroup) *CategoryController {
	return &CategoryController{
		categoryUC: categoryUC,
		rg:          rg,
	}
}
