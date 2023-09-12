package controller

import (
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/usecase"
	"final-project-enigma-clean/util/helper"

	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
)

type UserController struct {
	userUC usecase.UserCredentialUsecase
	rg     *gin.RouterGroup
}

// register handler
func (u *UserController) RegisterUserHandler(c *gin.Context) {
	var userRegist model.UserRegisterRequest

	//bind json
	if err := c.ShouldBindJSON(&userRegist); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	if err := u.userUC.RegisterUser(userRegist); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"Message": "Successfulyy Register"})

}

// login handler
// Di dalam controller layer, perbarui LoginUserHandler
func (u *UserController) LoginUserHandler(c *gin.Context) {
	var userLogin model.UserLoginRequest

	// Bind JSON
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"Error": err.Error()})
		return
	}

	// Memanggil usecase LoginUser
	userID, err := u.userUC.LoginUser(userLogin)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": err.Error()})
		return
	}

	// Generate JWT menggunakan email
	token, err := helper.GenerateJWT(userID)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"Error": "Failed to generate jwt"})
		return
	}

	// Mengirim token JWT sebagai respons
	slog.Infof("New user with email : %v and jwt : %v", userLogin.Email, token)
	c.JSON(200, gin.H{"Message": "Successfully Login", "Token": token})
}

// init route
func (u *UserController) Route() {
	{
		u.rg.POST("/register", u.RegisterUserHandler)
		u.rg.POST("/login", u.LoginUserHandler)
	}
}

func NewUserController(userUC usecase.UserCredentialUsecase, rg *gin.RouterGroup) *UserController {
	return &UserController{
		userUC: userUC,
		rg:     rg,
	}
}
