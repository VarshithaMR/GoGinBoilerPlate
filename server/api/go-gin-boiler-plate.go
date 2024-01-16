package api

import (
	"context"
	"log"
	"net/http"
)

// Method http Method in enum
type Method int

const (
	GET Method = iota
	POST
	PUT
	DELETE
	PATCH
)

// Operations contains the parameters of every endpoint exposed by go-gin-handler
type Operations struct {
	Handler    gin.HandlerFunc
	Method     Method
	Middleware []gin.HandlerFunc
	Path       string
}

type goGinEngine struct {
	*gin.Engine
	// ServerShutdown is called when the HTTP(S) server is shut down and done
	// handling all active connections and does not accept connections any more
	serverShutdown func()
	// User defined logger function.
	logger func(context.Context, string, ...interface{})
}

// GoGinHandler interface is wrapper on the http handler
type GoGinHandler interface {
	http.Handler
	WithLogger(func(context.Context, string, ...interface{})) GoGinHandler
	WithShutdown(func()) GoGinHandler
	WithOperations(rootPath string, authorize gin.HandlerFunc, operations ...*Operations) GoGinHandler
	WithGlobalMiddleware(globalMiddleware ...gin.HandlerFunc) GoGinHandler
	Logger(context.Context, string, ...interface{})
	ServerShutdown()
}

func NewGoGin() GoGinHandler {
	//LOG, _ := zap.NewProduction()
	goGin := &goGinEngine{
		Engine: gin.New(),
		serverShutdown: func() {
			log.Println("Server Shutdown ...")
		},
		logger: func(x context.Context, s string, arg ...interface{}) {
		},
	}
	return goGin
}

// implement all methods of interface
func (d *goGinEngine) Logger(ctx context.Context, message string, arguments ...interface{}) {
	d.logger(ctx, message, arguments)
}

func (d *goGinEngine) ServerShutdown() {
	d.serverShutdown()
}

func (d *goGinEngine) WithLogger(logger func(context.Context, string, ...interface{})) GoGinHandler {
	d.logger = logger
	return d
}

func (d *goGinEngine) WithShutdown(shutdown func()) GoGinHandler {
	d.serverShutdown = shutdown
	return d
}

func (d *goGinEngine) WithOperations(rootPath string, authorize gin.HandlerFunc, operations ...*Operations) GoGinHandler {
	group := d.Group(rootPath)
	if authorize != nil {
		group.Use(authorize)
	}
	for _, operation := range operations {
		createGroupWithMiddleware(group, operation.Middleware...).POST(operation.Path, operation.Handler)
	}
	return d
}

func createGroupWithMiddleware(group *gin.RouterGroup, globalMiddleware ...gin.HandlerFunc) *gin.RouterGroup {
	if len(globalMiddleware) == 0 {
		return group
	}
	return group.Group("", globalMiddleware...)
}
