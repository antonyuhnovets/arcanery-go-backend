// Set up application before run

package app

import (
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/hetonei/arcanery-go-backend/config"
	_ "github.com/hetonei/arcanery-go-backend/docs"
	v1 "github.com/hetonei/arcanery-go-backend/internal/controller/http/v1"
	"github.com/hetonei/arcanery-go-backend/internal/middleware"
	"github.com/hetonei/arcanery-go-backend/internal/service/lobby"
)

// @title         Arcanery
// @version       0.0.1
// @description   Give me request and I'll give you a power
// @license.name  MIT
func Run(cfg *config.Config) {
	app := gin.New()

	// run hub
	h := lobby.NewHub()
	go h.Run()

	m := map[string]string{
		"open":   "OnConnect",
		"chat":   "OnMessage",
		"create": "OnCreate",
		"close":  "OnDisconnect",
	}
	eventHub := lobby.NewEventHub(h, m)

	go eventHub.Run()

	// swagger docs
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// set routes
	v1.NewRouter(app, eventHub)

	// use middlewares
	app.Use(middleware.Cors())

	// sentry
	app.Use(sentrygin.New(sentrygin.Options{
		Repanic: true,
	}))

	app.Use(func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.Scope().SetTag("someRandomTag", "maybeYouNeedIt")
		}
		ctx.Next()
	})

	app.GET("/", func(ctx *gin.Context) {
		if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetExtra("unwantedQuery", "someQueryDataMaybe")
				hub.CaptureMessage("User provided unwanted query string, but we recovered just fine")
			})
		}
		ctx.Status(http.StatusOK)
	})

	app.GET("/foo", func(ctx *gin.Context) {
		// sentrygin handler will catch it just fine. Also, because we attached "someRandomTag"
		// in the middleware before, it will be sent through as well
		panic("y tho")
	})

	// run on given host, port
	app.Run(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port))
}
