package api

// Routes : ...
func (s *Server) Routes() {
	api := s.g.Group("/api")
	{
		// webhook API group
		webhook := api.Group("/webhook")
		webhook.POST("/constant", s.ConstantWebhook)

		api.POST("/sell-coin", s.CollateralSellCoinWebhook)
	}
}
