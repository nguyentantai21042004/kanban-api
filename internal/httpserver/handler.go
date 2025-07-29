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

	// Import this to execute the init function in docs.go which setups the Swagger docs.
	_ "gitlab.com/tantai-kanban/kanban-api/docs" // TODO: Generate docs package

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	Api         = "/api/v1"
	InternalApi = "internal/api/v1"
)

func (srv HTTPServer) mapHandlers() error {
	discord, err := discord.New(srv.l, srv.discord)
	if err != nil {
		return err
	}
	srv.gin.Use(middleware.Recovery(discord))

	// Swagger UI
	srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	scopeManager := scope.NewManager(srv.jwtSecretKey)
	// internalKey, err := srv.encrypter.Encrypt(srv.internalKey)
	// if err != nil {
	// 	srv.l.Fatal(context.Background(), err)
	// 	return err
	// }

	i18n.Init()

	// Middleware
	mw := middleware.New(srv.l, scopeManager)

	// Initialize WebSocket
	wsService.InitWebSocketHub()
	wsH := wsHTTP.New(wsService.GetHub())

	boardRepo := boardRepository.New(srv.l, srv.postgresDB)
	boardUC := boardUC.New(srv.l, boardRepo)
	boardH := boardHTTP.New(srv.l, boardUC, discord)

	listRepo := listRepository.New(srv.l, srv.postgresDB)
	listUC := listUC.New(srv.l, listRepo)
	listH := listHTTP.New(srv.l, listUC, discord)

	labelRepo := labelRepository.New(srv.l, srv.postgresDB)
	labelUC := labelUC.New(srv.l, labelRepo)
	labelH := labelHTTP.New(srv.l, labelUC, discord)

	cardRepo := cardRepository.New(srv.l, srv.postgresDB)
	cardUC := cardUC.New(srv.l, cardRepo, wsService.GetHub())
	cardH := cardHTTP.New(srv.l, cardUC, discord)

	roleRepo := roleRepository.New(srv.l, srv.postgresDB)
	roleUC := roleUC.New(srv.l, roleRepo)
	roleH := roleHTTP.New(srv.l, roleUC, discord)

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

	return nil
}
