package api

import (
	"encoding/json"

	"github.com/constant-money/constant-web-api/daos"
	"github.com/constant-money/constant-web-api/models"
	"github.com/constant-money/constant-web-api/serializers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ConstantWebhook :...
func (s *Server) ConstantWebhook(c *gin.Context) {

	var req serializers.WebhookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		s.logger.Error("s.ShouldBindJSON", zap.Error(err))
		return
	}

	hookData, err := json.Marshal(req)
	hook := models.Hook{
		Source: models.HookSourcePrimetrust,
		Data:   string(hookData),
		Status: models.HookStatusNew,
	}

	err = daos.WithDB(func(db *gorm.DB) error {
		return s.hookSvc.CreateHook(&hook)
	})

	if err != nil {
		s.logger.Error("s.hookSvc.CreateHook", zap.Error(err))
	}

	switch req.Type {
	case serializers.WebhookTypeOrder:
		var data serializers.WebhookRollbackRequest
		mapstructure.Decode(req.Data, &data)
		err := s.localSrv.RollbackMakerRequest(data.ID, data.CanceledOrderID)
		if err != nil {
			s.logger.Error("s.localSrv.RollbackMakerRequest", zap.Error(err))
		}

	case serializers.WebhookTypeUserWallet:
		var data serializers.WebhookUserWalletRequest
		mapstructure.Decode(req.Data, &data)
		err := s.userSvc.TransferConstant(&data)
		if err != nil {
			s.logger.Error("s.userSvc.TransferConstant", zap.Error(err))
		}

	case serializers.WebhookTypeKYC:
		var data serializers.WebhookKYCRequest
		mapstructure.Decode(req.Data, &data)

		err := s.userSvc.UpdateKYCStatusForUser(data.ID, data.PrimetrustContactStatus, data.PrimetrustContactError)
		if err != nil {
			s.logger.Error("s.userSvc.UpdateKYCStatusForUser", zap.Error(err))
		}

	case serializers.WebhookTypeTxHash:
		var data serializers.WebhookTxHashRequest
		mapstructure.Decode(req.Data, &data)
		err = s.reserveSvc.TxHashWebhook(&data)
		if err != nil {
			s.logger.Error("s.reserveSvc.TxHashWebhook", zap.Error(err))
		}

	case serializers.WebhookTypeReserve:
		var data serializers.WebhookReserveRequest
		mapstructure.Decode(req.Data, &data)

		err := s.reserveSvc.ReserveWebhook(data.ID)
		if err != nil {
			s.logger.Error("s.reserveSvc.ReserveWebhook", zap.Error(err))
		}

	case serializers.WebhookTypeCollateralLoan:
		var data serializers.WebhookCollateralLoanRequest
		mapstructure.Decode(req.Data, &data)
		err := s.collateralLoanSvc.Webhook(&data)
		if err != nil {
			s.logger.Error("s.collateralLoanSvc.Webhook", zap.Error(err))
		}

	}

	errTnx := daos.WithTransaction(func(tx *gorm.DB) error {
		if err != nil {
			hook.Result = err.Error()
		} else {
			hook.Result = "success"
		}

		hook.Status = models.HookStatusProcessed
		updateErr := s.hookSvc.UpdateHook(&hook)
		if updateErr != nil {
			return errors.Wrap(updateErr, "s.hookSvc.UpdateHook")
		}
		return nil
	})

	if errTnx != nil {
		s.logger.Error("WithTransaction", zap.Error(errTnx))
	}

	return
}
