package server

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func InitRouter() *gin.Engine {
	//db.Init(client, redisClient, node)
	router := gin.New()
	router.Use(gin.Recovery())
	router.GET("", func(c *gin.Context) {
		var request struct {
			Username string `from:"username"`
		}
		if err := c.ShouldBindQuery(&request); err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		logrus.Info(request.Username)
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to Sagi"})
	})
	//authMiddleware = prepareAuthMiddleware()
	//api := router.Group("/api/v1")
	//api.Use(authMiddleware.MiddlewareFunc())
	//registerSmsHandler(router, authMiddleware)
	//registerUserHandler(api)
	//registerPromoterHandler(api)
	//registerStarHandler(api)
	////registerModerator(api)
	//registerAnalytics(api)
	//registerCitiesHandler(api)
	//registerOrdersHandler(api)
	//registerBranchesHandler(api)
	//registerAdditionallyServices(api)
	//registerDealsHandler(api)
	//registerCategories(api)
	//registerBusinessUsersHandler(api)
	//integrations := router.Group("/integrations")
	//registerCardIntegration(integrations)
	//registerHelpHandler(api)
	//registerChatsHandler(api)
	//registerOffersHandler(api)
	//registerAutoPaymentHandler(api)
	//registerFavoriteHandler(api)
	//registerFriendsHandler(api)
	//registerImportantDateHandler(api)
	//
	//if !config.Instance.Server.Production {
	//	registerTestCaseHandler(api)
	//}
	//
	//customerPusher = push.NewPusher(config.Instance.Push.CustomerKey)
	//businessPusher = push.NewPusher(config.Instance.Push.BusinessKey)
	//amplitudeService = amplitude.NewAmplitudeService(config.Instance.Analytic.APIKey)
	return router
}
