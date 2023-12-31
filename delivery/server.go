package delivery

import (
	"final-project-enigma-clean/config"
	"final-project-enigma-clean/delivery/controller"
	"final-project-enigma-clean/delivery/middleware"
	"final-project-enigma-clean/manager"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
)

type Server struct {
	um   manager.UsecaseManager
	gin  *gin.Engine
	host string
	log  *logrus.Logger
}

func (s *Server) initMiddlewares() {
	// Create a Zap logger
	logger, _ := zap.NewProduction()
	defer logger.Sync() // Ensure logs are flushed

	// Use the logger
	s.gin.Use(middleware.ZapLogger(logger))
	s.gin.Use(middleware.ErrorHandler())

}

func (s *Server) initControllers() {
	rg := s.gin.Group("/api/v1")
	controller.NewUserController(s.um.UserUsecase(), rg).Route()
	controller.NewTypeAssetController(s.um.TypeAssetUseCase(), rg).Route()
	controller.NewStaffController(s.um.StaffUseCase(), rg).Route()
	controller.NewAssetController(s.um.AssetUsecase(), rg).Route()
	controller.NewCategoryController(s.um.CategoryUsecase(), rg).Route()
	controller.NewManageAssetController(s.um.ManageAssetUsecase(), rg).Route()
}

func (s *Server) Run() {
	s.initMiddlewares()
	s.initControllers()
	err := s.gin.Run(s.host)
	if err != nil {
		panic(err)
	}
}

func NewServer() *Server {

	cfg, err := config.NewDbConfig()
	if err != nil {
		fmt.Printf("Failed on db constructor %v", err.Error())
	}
	//define constructor dari infra
	im, err := manager.NewInfraManager(cfg)
	if err != nil {
		fmt.Printf("Failed on infra constructor %v", err.Error())
	}

	rm := manager.NewRepoManager(im)
	um := manager.NewUsecaseManager(rm)

	//untuk host
	host := fmt.Sprintf("%s:%s", cfg.ApiHost, cfg.ApiPort)
	//gin serv
	g := gin.Default()

	//init log
	log := logrus.New()
	return &Server{
		um:   um,
		gin:  g,
		host: host,
		log:  log,
	}
}
