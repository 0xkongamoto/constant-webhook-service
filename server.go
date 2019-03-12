package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/constant-money/constant-web-api/config"
	"github.com/constant-money/constant-web-api/daos"
	"github.com/constant-money/constant-web-api/services"
	"github.com/constant-money/constant-web-api/services/3rd/coinbase"
	"github.com/constant-money/constant-web-api/services/3rd/primetrust"
	"github.com/constant-money/constant-web-api/services/3rd/sendgrid"
	"github.com/constant-money/constant-web-api/templates/email"
	"github.com/constant-money/constant-webhook-service/api"
)

func main() {
	// load config
	conf := config.GetConfig()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to create zap logger: %v", err)
	}
	defer logger.Sync()

	// init daos
	if err := daos.Init(conf); err != nil {
		panic(err)
	}

	if err := daos.AutoMigrate(); err != nil {
		logger.Fatal("failed to auto migrate", zap.Error(err))
	}

	var (
		mailClient        = sendgrid.Init(conf)
		emailHelper       = email.New(mailClient, conf)
		primetrustService = primetrust.Init(conf.PrimetrustPrefix, conf.PrimetrustEmail, conf.PrimetrustPassword, conf.PrimetrustAccountID)
		gsClient          = services.InitGsClient(conf)
		hubspotService    = services.NewHubspotService(conf.HubspotHapiKey)
		userDAO           = daos.NewUser()
		reserveDAO        = daos.NewReserve()
		taskDAO           = daos.NewTask()

		storageSvc = services.InitStorageService(gsClient, userDAO)
		userSvc    = services.NewUserService(userDAO, reserveDAO, taskDAO, emailHelper, primetrustService, storageSvc, hubspotService, conf)

		masterAddrDAO = daos.NewMasterAddressDAO()
		txDAO         = daos.NewTx()
		hookDAO       = daos.NewHook()

		collateralLoanDAO = daos.NewCollateralLoan()
		collateralDAO     = daos.NewCollateral()
		coinbaseSvc       = coinbase.Init(conf)

		// local service
		orderDAO    = daos.NewOrder()
		makerDAO    = daos.NewMaker()
		shakerDAO   = daos.NewShaker()
		exchangeDAO = daos.NewExchange()
		firebaseDB  = services.InitFirebase(conf)
		localSrv    = services.InitLocalService(userDAO, orderDAO, makerDAO, shakerDAO, exchangeDAO, emailHelper, firebaseDB)

		// task service
		taskSvc = services.NewTaskService(taskDAO, conf)

		// reserve service
		reserveSvc = services.NewReserveService(reserveDAO, userDAO, txDAO, masterAddrDAO, taskSvc, primetrustService, conf)

		// hook service
		hookSvc = services.NewHookService(hookDAO)

		collateralLoanSvc = services.InitCollateralLoanService(userDAO, collateralDAO, collateralLoanDAO, coinbaseSvc, firebaseDB, emailHelper)
	)

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://*", "https://*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
	}))
	svr := api.NewServer(r, userSvc, reserveSvc, localSrv, hookSvc, collateralLoanSvc, logger, conf)
	svr.Routes()

	if err := r.Run(fmt.Sprintf(":%d", conf.Port)); err != nil {
		logger.Fatal("router.Run", zap.Error(err))
	}
}
