package meiporo

import "github.com/ant0ine/go-urlrouter"

type Router struct {
	Routers map[string]urlrouter.Router
	Routes  map[string][]urlrouter.Route
}

func NewRouter() *Router {
	return &Router{
		Routers: make(map[string]urlrouter.Router),
		Routes:  make(map[string][]urlrouter.Route),
	}
}

func (r *Router) Add(method, pathExp string, handlers ...Handler) {
	r.Routes[method] = append(r.Routes[method], urlrouter.Route{
		PathExp: pathExp,
		Dest:    handlers,
	})
}

func (r *Router) Start() {
	for method, routes := range r.Routes {
		router := urlrouter.Router{Routes: routes}
		if err := router.Start(); err != nil {
			panic(err)
		}
		r.Routers[method] = router
	}
}

func (r *Router) Match(method, path string) (*urlrouter.Route, map[string]string) {
	if router, ok := r.Routers[method]; ok {
		route, params, err := router.FindRoute(path)
		if err == nil {
			return route, params
		}
	}
	return nil, map[string]string{}
}
