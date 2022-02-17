package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kazakovichna/todoListPrjct/pkg/service"
	"time"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	config := cors.Config{
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"Authorization"},
		AllowMethods:     []string{"PUT", "DELETE", "GET", "POST"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
		auth.POST("/refresh", h.refreshToken)
	}

	api := router.Group("/api", h.UserIdentity)
	{
		desks := api.Group("/desks")
		{
			desks.POST("/", h.createDesk)
			desks.GET("/", h.getAllDesks)
			desks.GET("/:id", h.getDeskById)
			desks.PUT("/:id", h.updateDesk)
			desks.DELETE("/:id", h.deleteDesk)

			deskItem := desks.Group(":id/deskItems")
			{
				deskItem.GET("/", h.getAllItems)
			}

			lists := desks.Group(":id/lists")
			{
				lists.POST("/", h.createList)
				lists.GET("/", h.getAllLists)
				lists.GET("/:listId", h.getListById)
				lists.PUT("/:listId", h.updateList)
				lists.DELETE("/:listId", h.deleteList)

				items := lists.Group(":listId/items")
				{
					items.POST("/", h.createItem)
					//items.GET("/", h.getAllItems)
					items.GET("/:itemId", h.getItemById)
					items.PUT("/:itemId", h.updateItem)
					items.DELETE("/:itemId", h.deleteItem)
				}
			}
		}
	}

	return router
}
