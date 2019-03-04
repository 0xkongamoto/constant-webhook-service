package api

import (
	"strconv"

	"github.com/constant-money/constant-web-api/config"
	"github.com/constant-money/constant-web-api/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Server : struct
type Server struct {
	g                 *gin.Engine
	userSvc           *services.User
	reserveSvc        *services.ReserveService
	localSrv          *services.LocalService
	hookSvc           *services.HookService
	collateralLoanSvc *services.CollateralLoanService
	logger            *zap.Logger
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
func NewServer(g *gin.Engine,
	userSvc *services.User,
	reserveSvc *services.ReserveService,
	localSrv *services.LocalService,
	hookSvc *services.HookService,
	collateralLoanSvc *services.CollateralLoanService,
	logger *zap.Logger,
	config *config.Config) *Server {

	return &Server{
		g:                 g,
		userSvc:           userSvc,
		reserveSvc:        reserveSvc,
		localSrv:          localSrv,
		collateralLoanSvc: collateralLoanSvc,
		hookSvc:           hookSvc,
		logger:            logger,
		config:            config,
	}
}
