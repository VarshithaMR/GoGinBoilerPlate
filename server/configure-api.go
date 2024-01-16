package server

import (
	"GoGinBoilerPlate/server/api"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

const endpointv1 = "v1/go-gin"

func configureAPI(goGinHandler api.GoGinHandler, contextRoot string, goGinDomain service.GoGinDomain) api.GoGinHandler {

	// Setup global middleware like Logging, Panic handler
	goGinHandler = goGinHandler.WithGlobalMiddleware(dontPanic, gintrace.Middleware("go-gin-handler"))

	// define all other endpoints as different operations
	createQuizOperation := &api.Operations{
		Method: api.POST,
		Path:   endpointv1,
		Handler: func(ctx *gin.Context) {
			goGinDomain.Create(ctx.Writer, ctx.Request)
		},
		Middleware: []gin.HandlerFunc{},
	}

	goGinHandler = goGinHandler.WithOperations(contextRoot, nil, createQuizOperation)

	// Setup shutdown
	goGinHandler = goGinHandler.WithShutdown(func() {
		goGinHandler.Logger(context.Background(), "Shutting down %s", "quiz-domain-handler server")
	})
	return goGinHandler
}

func dontPanic(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			//log.Error(ctx, "quiz-domain-handler-panic!!!: %v", err)
			ctx.Writer.WriteHeader(http.StatusInternalServerError)
			ctx.Writer.Write([]byte("go-gin-handler-panic!!!"))
		}
	}()
	ctx.Next()
}
