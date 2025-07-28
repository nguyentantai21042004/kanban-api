package httpserver

import (
	// ginSwagger "github.com/swaggo/gin-swagger"

	"gitlab.com/tantai-kanban/kanban-api/internal/middleware"
	"gitlab.com/tantai-kanban/kanban-api/pkg/discord"
	"gitlab.com/tantai-kanban/kanban-api/pkg/i18n"
	"gitlab.com/tantai-kanban/kanban-api/pkg/scope"
	// Import this to execute the init function in docs.go which setups the Swagger docs.
	// swaggerFiles "github.com/swaggo/files"
	// ginSwagger "github.com/swaggo/gin-swagger"
	// _ "gitlab.com/tantai-kanban/kanban-api/docs" // TODO: Generate docs package
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
	// smtpUC := smtpUC.New(srv.l, srv.smtpConfig)

	// Middleware
	mw := middleware.New(srv.l, scopeManager)

	// uploadRepoPostgres := uploadRepoPostgres.New(srv.l, srv.postgresDB)
	// uploadUC := uploadUC.New(srv.l, uploadRepoPostgres, srv.cloudinary, userUC)
	// uploadH := uploadHTTP.New(srv.l, uploadUC, discord)

	// Apply locale middleware
	srv.gin.Use(mw.Locale()).Use(mw.Cors())
	// api := srv.gin.Group(Api)

	return nil
}
