package httpserver

import (
	"context"

	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	"gitlab.com/tantai-kanban/kanban-api/pkg/i18n"
	"gitlab.com/tantai-kanban/kanban-api/pkg/position"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"

	boardHTTP "gitlab.com/tantai-kanban/kanban-api/internal/boards/delivery/http"
	boardRepository "gitlab.com/tantai-kanban/kanban-api/internal/boards/repository/postgres"
	boardUC "gitlab.com/tantai-kanban/kanban-api/internal/boards/usecase"

	listHTTP "gitlab.com/tantai-kanban/kanban-api/internal/lists/delivery/http"
	listRepository "gitlab.com/tantai-kanban/kanban-api/internal/lists/repository/postgres"
	listUC "gitlab.com/tantai-kanban/kanban-api/internal/lists/usecase"

	labelHTTP "gitlab.com/tantai-kanban/kanban-api/internal/labels/delivery/http"
	labelRepository "gitlab.com/tantai-kanban/kanban-api/internal/labels/repository/postgres"
	labelUC "gitlab.com/tantai-kanban/kanban-api/internal/labels/usecase"

	cardHTTP "gitlab.com/tantai-kanban/kanban-api/internal/cards/delivery/http"
	cardRepository "gitlab.com/tantai-kanban/kanban-api/internal/cards/repository/postgres"
	cardUC "gitlab.com/tantai-kanban/kanban-api/internal/cards/usecase"

	roleHTTP "gitlab.com/tantai-kanban/kanban-api/internal/role/delivery/http"
	roleRepository "gitlab.com/tantai-kanban/kanban-api/internal/role/repository/postgres"
	roleUC "gitlab.com/tantai-kanban/kanban-api/internal/role/usecase"

	wsHTTP "gitlab.com/tantai-kanban/kanban-api/internal/websocket/delivery/http"
	wsService "gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"

	uploadHTTP "gitlab.com/tantai-kanban/kanban-api/internal/upload/delivery/http"
	uploadRepository "gitlab.com/tantai-kanban/kanban-api/internal/upload/repository/postgres"
	uploadUC "gitlab.com/tantai-kanban/kanban-api/internal/upload/usecase"

	commentHTTP "gitlab.com/tantai-kanban/kanban-api/internal/comments/delivery/http"
	commentRepository "gitlab.com/tantai-kanban/kanban-api/internal/comments/repository/postgres"
	commentUC "gitlab.com/tantai-kanban/kanban-api/internal/comments/usecase"

	// Import this to execute the init function in docs.go which setups the Swagger docs.
	_ "gitlab.com/tantai-kanban/kanban-api/docs" // TODO: Generate docs package

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	userHTTP "gitlab.com/tantai-kanban/kanban-api/internal/user/delivery/http"
	userRepository "gitlab.com/tantai-kanban/kanban-api/internal/user/repository/postgres"
	userUC "gitlab.com/tantai-kanban/kanban-api/internal/user/usecase"

	authHTTP "gitlab.com/tantai-kanban/kanban-api/internal/auth/delivery/http"
	authUC "gitlab.com/tantai-kanban/kanban-api/internal/auth/usecase"

	"github.com/gin-gonic/gin"
	"gitlab.com/tantai-kanban/kanban-api/pkg/response"
)

const (
	Api         = "/api/v1"
	InternalApi = "internal/api/v1"
)

func (srv HTTPServer) mapHandlers() error {
	discord, err := discord.New(srv.l, srv.discord, discord.DefaultConfig())
	if err != nil {
		return err
	}
	srv.gin.Use(middleware.Recovery(discord))

	// Health check endpoint
	srv.gin.GET("/health", srv.healthCheck)
	srv.gin.GET("/ready", srv.readyCheck)
	srv.gin.GET("/live", srv.liveCheck)

	// Swagger UI
	srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	scopeUC := scope.New(srv.jwtSecretKey)
	// internalKey, err := srv.encrypter.Encrypt(srv.internalKey)
	// if err != nil {
	// 	srv.l.Fatal(context.Background(), err)
	// 	return err
	// }

	i18n.Init()

	// Middleware
	mw := middleware.New(srv.l, scopeUC)

	// Initialize WebSocket
	if err := wsService.InitWebSocketHub(srv.l); err != nil {
		srv.l.Error(context.Background(), "Failed to initialize WebSocket hub", "error", err)
		return err
	}
	wsH := wsHTTP.New(wsService.GetHub(), scopeUC, srv.l)

	roleRepo := roleRepository.New(srv.l, srv.postgresDB)
	roleUC := roleUC.New(srv.l, roleRepo)
	roleH := roleHTTP.New(srv.l, roleUC, discord)

	uploadRepository := uploadRepository.New(srv.l, srv.postgresDB)
	uploadUC := uploadUC.New(srv.l, uploadRepository, srv.minioClient)
	uploadH := uploadHTTP.New(srv.l, uploadUC, discord)

	userRepo := userRepository.New(srv.l, srv.postgresDB)
	userUC := userUC.New(srv.l, userRepo)
	userH := userHTTP.New(srv.l, userUC, discord)

	authUC := authUC.New(srv.l, srv.encrypter, scopeUC, userUC, roleUC)
	authH := authHTTP.New(srv.l, authUC, discord)

	// Fractical Indexing Algorithm
	positionUC := position.NewPositionManager()

	boardRepo := boardRepository.New(srv.l, srv.postgresDB)
	boardUC := boardUC.New(srv.l, boardRepo, wsService.GetHub(), userUC, roleUC)
	boardH := boardHTTP.New(srv.l, boardUC, discord)

	listRepo := listRepository.New(srv.l, srv.postgresDB)
	listUC := listUC.New(srv.l, listRepo, wsService.GetHub(), positionUC, boardUC)
	listH := listHTTP.New(srv.l, listUC, discord)

	labelRepo := labelRepository.New(srv.l, srv.postgresDB)
	labelUC := labelUC.New(srv.l, labelRepo)
	labelH := labelHTTP.New(srv.l, labelUC, discord)

	cardRepo := cardRepository.New(srv.l, srv.postgresDB)
	cardUC := cardUC.New(srv.l, cardRepo, wsService.GetHub(), positionUC, boardUC, listUC, userUC)
	cardH := cardHTTP.New(srv.l, cardUC, discord)

	commentRepo := commentRepository.New(srv.l, srv.postgresDB)
	commentUC := commentUC.New(srv.l, commentRepo, userUC, cardUC, wsService.GetHub())
	commentH := commentHTTP.New(srv.l, commentUC, discord)

	// Apply locale middleware
	srv.gin.Use(mw.Locale())
	api := srv.gin.Group(Api)

	// WebSocket routes with special CORS middleware
	websocketGroup := api.Group("/websocket")
	websocketGroup.Use(mw.CorsForWebSocket())
	wsHTTP.MapWebSocketRoutes(websocketGroup, wsH, mw)

	// Apply regular CORS middleware for all other routes
	srv.gin.Use(mw.Cors())

	// Routes with CORS middleware
	boardsGroup := api.Group("/boards")
	boardsGroup.Use(mw.Cors())
	boardHTTP.MapBoardRoutes(boardsGroup, boardH, mw)

	listsGroup := api.Group("/lists")
	listsGroup.Use(mw.Cors())
	listHTTP.MapListRoutes(listsGroup, listH, mw)

	labelsGroup := api.Group("/labels")
	labelsGroup.Use(mw.Cors())
	labelHTTP.MapLabelRoutes(labelsGroup, labelH, mw)

	cardsGroup := api.Group("/cards")
	cardsGroup.Use(mw.Cors())
	cardHTTP.MapCardRoutes(cardsGroup, cardH, mw)

	rolesGroup := api.Group("/roles")
	rolesGroup.Use(mw.Cors())
	roleHTTP.MapRoleRoutes(rolesGroup, roleH, mw)

	uploadsGroup := api.Group("/uploads")
	uploadsGroup.Use(mw.Cors())
	uploadHTTP.MapUploadRoutes(uploadsGroup, uploadH, mw)

	commentsGroup := api.Group("/comments")
	commentsGroup.Use(mw.Cors())
	commentHTTP.MapCommentRoutes(commentsGroup, commentH, mw)

	cardCommentsGroup := api.Group("/cards/:id")
	cardCommentsGroup.Use(mw.Cors())
	commentHTTP.MapCardCommentRoutes(cardCommentsGroup, commentH, mw)

	// Map routes with CORS middleware
	authGroup := api.Group("/auth")
	authGroup.Use(mw.Cors())
	authHTTP.MapAuthRoutes(authGroup, authH, mw)

	usersGroup := api.Group("/users")
	usersGroup.Use(mw.Cors())
	userHTTP.MapUserRoutes(usersGroup, userH, mw)

	return nil
}

// healthCheck handles health check requests
// @Summary Health Check
// @Description Check if the API is healthy
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API is healthy"
// @Router /health [get]
func (srv HTTPServer) healthCheck(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "healthy",
		"message": "From Tan Tai API V1 With Love",
		"version": "1.0.0",
		"service": "kanban-api",
	})
}

// readyCheck handles readiness check requests
// @Summary Readiness Check
// @Description Check if the API is ready to serve traffic
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API is ready"
// @Router /ready [get]
func (srv HTTPServer) readyCheck(c *gin.Context) {
	// Check database connection
	if err := srv.postgresDB.PingContext(c.Request.Context()); err != nil {
		c.JSON(503, gin.H{
			"status":  "not ready",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	response.OK(c, gin.H{
		"status":   "ready",
		"message":  "From Tan Tai API V1 With Love",
		"version":  "1.0.0",
		"service":  "kanban-api",
		"database": "connected",
	})
}

// liveCheck handles liveness check requests
// @Summary Liveness Check
// @Description Check if the API is alive
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "API is alive"
// @Router /live [get]
func (srv HTTPServer) liveCheck(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "alive",
		"message": "From Tan Tai API V1 With Love",
		"version": "1.0.0",
		"service": "kanban-api",
	})
}
