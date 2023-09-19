package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"final-project-enigma-clean/__mock__/usecasemock"
	"final-project-enigma-clean/model"
	"final-project-enigma-clean/usecase"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RegisterControllerTestSuite struct {
	suite.Suite
	usecase *usecasemock.UserCredentialsMock
	router  *gin.Engine
}

func (suite *RegisterControllerTestSuite) SetupTest() {
	suite.usecase = new(usecasemock.UserCredentialsMock)
	suite.router = gin.Default()
}

func TestRegisterUserTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterControllerTestSuite))
}

func (suite *RegisterControllerTestSuite) TestCreateUserHandler_Success() {
	mockData := model.UserRegisterRequest{
		Email:    "elliz@gmail.com",
		Password: "N@ufa282",
		Name:     "pal",
		IsActive: true,
	}

	suite.usecase.On("RegisterUser", mockData).Return(nil)
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	marshal, err := json.Marshal(mockData)
	assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	var userResp model.UserRegisterRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}
func (suite *RegisterControllerTestSuite) TestCreateUserHandler_JSONInvalid() {
	mockData := model.UserRegisterRequest{
		Email:    "elliz@gmail.com",
		Password: "N@ufa282",
		Name:     "pal",
		IsActive: true,
	}

	suite.usecase.On("RegisterUser", mockData).Return(nil)
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	// marshal, err := json.Marshal(mockData)
	// assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/register", nil)
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	var userResp model.UserRegisterRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

func (suite *RegisterControllerTestSuite) TestCreateUser_InvalidEmail() {
	mockData := model.UserRegisterRequest{
		Email:    "elliz", // Invalid email format
		Password: "N@ufa282",
		Name:     "pal",
		IsActive: true,
	}

	suite.usecase.On("RegisterUser", mockData).Return(errors.New("invalid email"))
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	marshal, err := json.Marshal(mockData)
	assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	var userResp model.UserRegisterRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), http.StatusBadRequest, record.Code)
}

// login testing
func (suite *RegisterControllerTestSuite) TestLoginUser_Success() {
	mockData := model.UserLoginRequest{
		Email:    "ellizavad@gmail.com",
		Password: "N@21457asdw",
	}

	suite.usecase.On("LoginUser", mockData).Return(nil)
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	marshal, err := json.Marshal(mockData)
	assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/login", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	var userResp model.UserLoginRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

// login otp
func (suite *RegisterControllerTestSuite) TestLoginUserOTP_Success() {
	// Create a new struct for mockData
	mockData := struct {
		Email string `json:"email"`
		OTP   int    `json:"otp"`
	}{
		Email: "ellizavad@gmail.com",
		OTP:   287303,
	}

	// Set up OTP map with a valid OTP value
	usecase.OTPMap["ellizavad@gmail.com"] = 287303

	suite.usecase.On("LoginUser", mockData).Return(nil)
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	marshal, err := json.Marshal(mockData)
	assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/login/email-otp/start", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	fmt.Println("Response Status Code:", record.Code)

	// Debug: Print the response body
	fmt.Println("Response Body:", string(response))

	var userResp model.UserLoginOTPRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), 200, record.Code)
}

func (suite *RegisterControllerTestSuite) TestLoginUserOTP_Fail() {
	mockData := model.UserLoginOTPRequest{
		Email: "ellizavad@gmail.com",
		OTP:   875502,
	}

	suite.usecase.On("LoginUser", mockData).Return(errors.New("otp is invalid / expired"))
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	marshal, err := json.Marshal(mockData)
	assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/login/email-otp/start", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	var userResp model.UserLoginOTPRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), 401, record.Code)
}

// forgot pass
func (suite *RegisterControllerTestSuite) TestChangePass_Success() {
	mockData := model.UserLoginRequest{
		Email: "ellizavad@gmail.com",
	}

	suite.usecase.On("LoginUser", mockData).Return(nil)
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	record := httptest.NewRecorder()

	marshal, err := json.Marshal(mockData)
	assert.NoError(suite.T(), err)

	request, err := http.NewRequest(http.MethodPost, "/api/v1/change-password", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(record, request)
	response := record.Body.Bytes()

	var userResp model.UserLoginRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), http.StatusOK, record.Code)
}

func (suite *RegisterControllerTestSuite) TestChangePass_InvalidJSONFormat() {
	// Create an invalid JSON request body (missing closing brace)
	requestBody := []byte(`{"email":"ellizavad.com"`)

	suite.usecase.On("LoginUser", requestBody).Return(errors.New("Bad json format"))
	mockRg := suite.router.Group("/api/v1")
	NewUserController(suite.usecase, mockRg).Route()

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	marshal, err := json.Marshal(requestBody)
	assert.NoError(suite.T(), err)

	// Create an HTTP request with the invalid JSON body
	request, err := http.NewRequest(http.MethodPost, "/api/v1/change-password", bytes.NewBuffer(marshal))
	assert.NoError(suite.T(), err)

	suite.router.ServeHTTP(recorder, request)
	response := recorder.Body.Bytes()

	var userResp model.UserLoginRequest
	json.Unmarshal(response, &userResp)
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)
}

// forgot pass otp
func (suite *RegisterControllerTestSuite) TestForgotPassOTP_InvalidJSONFormat() {
	// Create an invalid JSON request body (missing closing brace)
	requestBody := []byte(`{"id":"123", "email":"ellizavad@gmail.com", "otp":123, "old_password":"oldpass", "new_password":"newpass"`)

	// Create an HTTP request with the invalid JSON body
	request, err := http.NewRequest(http.MethodPost, "/api/v1/forgot-password/start", bytes.NewBuffer(requestBody))
	assert.NoError(suite.T(), err)

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Set up the router and controller
	router := gin.Default()
	userController := NewUserController(suite.usecase, router.Group("/api/v1"))
	userController.Route()

	// Serve the HTTP request
	router.ServeHTTP(recorder, request)

	// Check the response status code
	assert.Equal(suite.T(), http.StatusBadRequest, recorder.Code)

}
