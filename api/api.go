package api

import (
	"strconv"

	"github.com/constant-money/constant-web-api/config"
	"github.com/constant-money/constant-web-api/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	g                 *gin.Engine
	userSvc           *services.User
	reserveSvc        *services.ReserveService
	localSrv          *services.LocalService
	logger            *zap.Logger
	storageSvc        *services.StorageService
	hookSvc           *services.HookService
	countrySvc        *services.CountryService
	collateralLoanSvc *services.CollateralLoanService
	config            *config.Config
}

func (s *Server) pagingFromContext(c *gin.Context) (int, int) {
	var (
		pageS  = c.DefaultQuery("page", "1")
		limitS = c.DefaultQuery("limit", "10")
		page   int
		limit  int
		err    error
	)

	page, err = strconv.Atoi(pageS)
	if err != nil {
		page = 1
	}

	limit, err = strconv.Atoi(limitS)
	if err != nil {
		limit = 10
	}

	return page, limit
}

// NewServer : userSvc, reserveSvc, localSrv
func NewServer(g *gin.Engine, userSvc *services.User, reserveSvc *services.ReserveService, localSrv *services.LocalService, logger *zap.Logger, storageSvc *services.StorageService, hookSvc *services.HookService, countrySvc *services.CountryService, collateralLoanSvc *services.CollateralLoanService, config *config.Config) *Server {
	return &Server{
		g:                 g,
		userSvc:           userSvc,
		reserveSvc:        reserveSvc,
		localSrv:          localSrv,
		logger:            logger,
		storageSvc:        storageSvc,
		hookSvc:           hookSvc,
		countrySvc:        countrySvc,
		collateralLoanSvc: collateralLoanSvc,
		config:            config,
	}
}
