package httpserver

import (
	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	"gitlab.com/tantai-kanban/kanban-api/pkg/i18n"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"

	boardHTTP "gitlab.com/tantai-kanban/kanban-api/internal/boards/delivery/http"
	boardRepository "gitlab.com/tantai-kanban/kanban-api/internal/boards/repository/postgres"
	boardUC "gitlab.com/tantai-kanban/kanban-api/internal/boards/usecase"
	wsHTTP "gitlab.com/tantai-kanban/kanban-api/internal/websocket/delivery/http"
	wsService "gitlab.com/tantai-kanban/kanban-api/internal/websocket/service"

	// Import this to execute the init function in docs.go which setups the Swagger docs.
	_ "gitlab.com/tantai-kanban/kanban-api/docs" // TODO: Generate docs package
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

	//swagger api
	// srv.gin.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	scopeManager := scope.NewManager(srv.jwtSecretKey)
	// internalKey, err := srv.encrypter.Encrypt(srv.internalKey)
	// if err != nil {
	// 	srv.l.Fatal(context.Background(), err)
	// 	return err
	// }

	i18n.Init()

	// Middleware
	mw := middleware.New(srv.l, scopeManager)

	boardRepo := boardRepository.New(srv.l, srv.postgresDB)
	boardUC := boardUC.New(srv.l, boardRepo)
	boardH := boardHTTP.New(srv.l, boardUC, discord)

	// Initialize WebSocket
	wsService.InitWebSocketHub()
	wsH := wsHTTP.New(wsService.GetHub())

	// Apply locale middleware
	srv.gin.Use(mw.Locale()).Use(mw.Cors())
	api := srv.gin.Group(Api)
	boardHTTP.MapBoardRoutes(api.Group("/boards"), boardH, mw)
	wsHTTP.MapWebSocketRoutes(api.Group("/websocket"), wsH, mw)

	return nil
}
