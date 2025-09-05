package httpserver

import (
	"context"

	"github.com/nguyentantai21042004/kanban-api/internal/middleware"
	"github.com/nguyentantai21042004/kanban-api/pkg/discord"
	"github.com/nguyentantai21042004/kanban-api/pkg/i18n"
	"github.com/nguyentantai21042004/kanban-api/pkg/position"
	"github.com/nguyentantai21042004/kanban-api/pkg/scope"

	boardHTTP "github.com/nguyentantai21042004/kanban-api/internal/boards/delivery/http"
	boardRepository "github.com/nguyentantai21042004/kanban-api/internal/boards/repository/postgres"
	boardUC "github.com/nguyentantai21042004/kanban-api/internal/boards/usecase"

	listHTTP "github.com/nguyentantai21042004/kanban-api/internal/lists/delivery/http"
	listRepository "github.com/nguyentantai21042004/kanban-api/internal/lists/repository/postgres"
	listUC "github.com/nguyentantai21042004/kanban-api/internal/lists/usecase"

	labelHTTP "github.com/nguyentantai21042004/kanban-api/internal/labels/delivery/http"
	labelRepository "github.com/nguyentantai21042004/kanban-api/internal/labels/repository/postgres"
	labelUC "github.com/nguyentantai21042004/kanban-api/internal/labels/usecase"

	cardHTTP "github.com/nguyentantai21042004/kanban-api/internal/cards/delivery/http"
	cardRepository "github.com/nguyentantai21042004/kanban-api/internal/cards/repository/postgres"
	cardUC "github.com/nguyentantai21042004/kanban-api/internal/cards/usecase"

	roleHTTP "github.com/nguyentantai21042004/kanban-api/internal/role/delivery/http"
	roleRepository "github.com/nguyentantai21042004/kanban-api/internal/role/repository/postgres"
	roleUC "github.com/nguyentantai21042004/kanban-api/internal/role/usecase"

	wsHTTP "github.com/nguyentantai21042004/kanban-api/internal/websocket/delivery/http"
	wsService "github.com/nguyentantai21042004/kanban-api/internal/websocket/service"

	uploadHTTP "github.com/nguyentantai21042004/kanban-api/internal/upload/delivery/http"
	uploadRepository "github.com/nguyentantai21042004/kanban-api/internal/upload/repository/postgres"
	uploadUC "github.com/nguyentantai21042004/kanban-api/internal/upload/usecase"

	commentHTTP "github.com/nguyentantai21042004/kanban-api/internal/comments/delivery/http"
	commentRepository "github.com/nguyentantai21042004/kanban-api/internal/comments/repository/postgres"
	commentUC "github.com/nguyentantai21042004/kanban-api/internal/comments/usecase"

	adminHTTP "github.com/nguyentantai21042004/kanban-api/internal/admin/delivery/http"
	adminUC "github.com/nguyentantai21042004/kanban-api/internal/admin/usecase"

	// Import this to execute the init function in docs.go which setups the Swagger docs.
	_ "github.com/nguyentantai21042004/kanban-api/docs" // TODO: Generate docs package

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	userHTTP "github.com/nguyentantai21042004/kanban-api/internal/user/delivery/http"
	userRepository "github.com/nguyentantai21042004/kanban-api/internal/user/repository/postgres"
	userUC "github.com/nguyentantai21042004/kanban-api/internal/user/usecase"

	authHTTP "github.com/nguyentantai21042004/kanban-api/internal/auth/delivery/http"
	authUC "github.com/nguyentantai21042004/kanban-api/internal/auth/usecase"

	"github.com/gin-gonic/gin"
	"github.com/nguyentantai21042004/kanban-api/pkg/response"
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
	boardUC := boardUC.New(srv.l, boardRepo, wsService.GetHub(), userUC, roleUC, nil)
	boardH := boardHTTP.New(srv.l, boardUC, discord)

	listRepo := listRepository.New(srv.l, srv.postgresDB)
	listUC := listUC.New(srv.l, listRepo, wsService.GetHub(), positionUC, boardUC, userUC, roleUC)
	listH := listHTTP.New(srv.l, listUC, discord)
	boardUC.SetList(listUC)

	labelRepo := labelRepository.New(srv.l, srv.postgresDB)
	labelUC := labelUC.New(srv.l, labelRepo)
	labelH := labelHTTP.New(srv.l, labelUC, discord)

	cardRepo := cardRepository.New(srv.l, srv.postgresDB)
	cardUC := cardUC.New(srv.l, cardRepo, wsService.GetHub(), positionUC, boardUC, listUC, userUC, roleUC)
	cardH := cardHTTP.New(srv.l, cardUC, discord)

	commentRepo := commentRepository.New(srv.l, srv.postgresDB)
	commentUC := commentUC.New(srv.l, commentRepo, userUC, cardUC, wsService.GetHub())
	commentH := commentHTTP.New(srv.l, commentUC, discord)

	// Apply locale + metrics middleware
	srv.gin.Use(mw.Locale())
	srv.gin.Use(mw.Metrics())
	srv.gin.Use(mw.Cors())

	// routes
	api := srv.gin.Group(Api)
	boardHTTP.MapBoardRoutes(api.Group("/boards"), boardH, mw)
	listHTTP.MapListRoutes(api.Group("/lists"), listH, mw)
	labelHTTP.MapLabelRoutes(api.Group("/labels"), labelH, mw)
	cardHTTP.MapCardRoutes(api.Group("/cards"), cardH, mw)
	roleHTTP.MapRoleRoutes(api.Group("/roles"), roleH, mw)
	uploadHTTP.MapUploadRoutes(api.Group("/uploads"), uploadH, mw)
	commentHTTP.MapCommentRoutes(api.Group("/comments"), commentH, mw)
	commentHTTP.MapCardCommentRoutes(api.Group("/cards/:id"), commentH, mw)

	// WebSocket routes with special CORS middleware
	websocketGroup := api.Group("/websocket")
	websocketGroup.Use(mw.CorsForWebSocket())
	wsHTTP.MapWebSocketRoutes(websocketGroup, wsH, mw)

	// Admin routes
	adminUC := adminUC.New(srv.l, userUC, boardUC, cardUC, commentUC, roleUC, wsService.GetHub())
	adminH := adminHTTP.New(srv.l, adminUC, discord)
	adminHTTP.MapAdminRoutes(api.Group("/admin"), adminH, mw)

	authHTTP.MapAuthRoutes(api.Group("/auth"), authH, mw)
	userHTTP.MapUserRoutes(api.Group("/users"), userH, mw)

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
