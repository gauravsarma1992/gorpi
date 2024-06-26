package routing

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gauravsarma1992/go-rest-api/gorpi/api"
	"github.com/gauravsarma1992/go-rest-api/gorpi/middlewares"
	"github.com/gauravsarma1992/go-rest-api/gorpi/models"
	"github.com/gin-gonic/gin"
)

const (
	ContextRouteKey = "route"
)

type (
	Authentication struct{}

	RouteManager struct {
		apiEngine      *gin.Engine
		trie           *Trie
		middlwareStack *middlewares.MiddlewareStack
	}

	Route struct {
		RequestURI    string
		RequestMethod string
		Handler       api.ApiHandlerFunc
		ResourceModel models.ResourceModel
	}
)

func NewRouteManager(apiEngine *gin.Engine, ms *middlewares.MiddlewareStack) (rm *RouteManager, err error) {
	rm = &RouteManager{
		apiEngine:      apiEngine,
		middlwareStack: ms,
	}
	// The noroute handler handles the routes for routes which are
	// not defined. Since we are not defining any routes on the
	// gin context, everything will be handled by the root handler
	rm.apiEngine.NoRoute(rm.RootHandler)
	rm.trie = NewTrie(rm.GetDefaultBaseHandler())
	return
}

func (rm *RouteManager) GetDefaultBaseHandler() (route *Route) {
	route = &Route{
		RequestURI:    "/",
		RequestMethod: "*",
		Handler:       rm.BaseHandler,
	}
	return
}

func (rm *RouteManager) BaseHandler(req *api.Request, resp *api.Response) (err error) {
	return
}

func (rm *RouteManager) RootHandler(c *gin.Context) {
	var (
		route *Route
		err   error
		ctx   context.Context
	)

	path := fmt.Sprintf("%s-%s", c.Request.Method, c.Request.RequestURI)
	if route, err = rm.GetRoute(path); err != nil {
		log.Println(err, "Printing all routes: ", rm.trie.String())
		c.JSON(400, gin.H{
			"request": path,
			"message": "Route not found - " + err.Error(),
		})
		return
	}

	ctx = context.WithValue(context.Background(), ContextRouteKey, route)

	if err = rm.middlwareStack.Exec(ctx, c, route.Handler); err != nil {
		c.JSON(500, gin.H{
			"request": path,
			"message": "Handler errored - " + err.Error(),
		})
	}
	return
}

func (rm *RouteManager) AddRoutes(route *Route) (err error) {
	if _, err = rm.trie.AddPath(route); err != nil {
		return
	}
	log.Println("Adding route path", route.GetName())
	return
}

func (rm *RouteManager) GetRoute(path string) (route *Route, err error) {
	var (
		pathNode *Node
	)
	if pathNode, err = rm.trie.GetNode(path); err != nil {
		return
	}
	if pathNode.Route == nil {
		err = errors.New("Route not found for path" + path)
		return
	}
	route = pathNode.Route
	return
}

func (route *Route) GetName() (name string) {
	name = fmt.Sprintf("%s-%s", route.RequestMethod, route.RequestURI)
	return
}
