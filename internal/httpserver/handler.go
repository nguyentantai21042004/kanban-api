package httpserver

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	"gitlab.com/tantai-kanban/kanban-api/pkg/i18n"
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
	wsService.InitWebSocketHub()
	wsH := wsHTTP.New(wsService.GetHub(), scopeUC)

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

	boardRepo := boardRepository.New(srv.l, srv.postgresDB)
	boardUC := boardUC.New(srv.l, boardRepo, userUC, roleUC, wsService.GetHub())
	boardH := boardHTTP.New(srv.l, boardUC, discord)

	listRepo := listRepository.New(srv.l, srv.postgresDB)
	listUC := listUC.New(srv.l, listRepo, wsService.GetHub())
	listH := listHTTP.New(srv.l, listUC, discord)

	labelRepo := labelRepository.New(srv.l, srv.postgresDB)
	labelUC := labelUC.New(srv.l, labelRepo)
	labelH := labelHTTP.New(srv.l, labelUC, discord)

	cardRepo := cardRepository.New(srv.l, srv.postgresDB)
	cardUC := cardUC.New(srv.l, cardRepo, wsService.GetHub())
	cardH := cardHTTP.New(srv.l, cardUC, discord)

	commentRepo := commentRepository.New(srv.l, srv.postgresDB)
	commentUC := commentUC.New(srv.l, commentRepo, userUC, cardUC, wsService.GetHub())
	commentH := commentHTTP.New(srv.l, commentUC, discord)

	// Apply locale middleware
	srv.gin.Use(mw.Locale()).Use(mw.Cors())
	api := srv.gin.Group(Api)

	// WebSocket
	wsHTTP.MapWebSocketRoutes(api.Group("/websocket"), wsH, mw)

	// Routes
	boardHTTP.MapBoardRoutes(api.Group("/boards"), boardH, mw)
	listHTTP.MapListRoutes(api.Group("/lists"), listH, mw)
	labelHTTP.MapLabelRoutes(api.Group("/labels"), labelH, mw)
	cardHTTP.MapCardRoutes(api.Group("/cards"), cardH, mw)
	roleHTTP.MapRoleRoutes(api.Group("/roles"), roleH, mw)
	uploadHTTP.MapUploadRoutes(api.Group("/uploads"), uploadH, mw)
	commentHTTP.MapCommentRoutes(api.Group("/comments"), commentH, mw)
	commentHTTP.MapCardCommentRoutes(api.Group("/cards/:id"), commentH, mw)

	// Map routes
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
