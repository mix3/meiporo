package meiporo

import (
	"log"
	"net/http"
	"os"
)

type Meiporo struct {
	Router               *Router
	Middlewares          []Handler
	Logger               *log.Logger
	NotFoundHandler      Handler
	CreateContextHandler Handler
}

func New() *Meiporo {
	m := &Meiporo{
		Router:      NewRouter(),
		Middlewares: []Handler{},
		Logger:      log.New(os.Stdout, "[meiporo] ", 0),
	}
	m.Use(&LoggerMiddleware{})
	m.Use(&RecoverMiddleware{})
	m.Use(&StaticFileMiddleware{PublicPath: "public"})
	return m
}

func (m *Meiporo) CreateContext(res http.ResponseWriter, req *http.Request) *Context {
	route, params := m.Router.Match(req.Method, req.URL.Path)
	c := &Context{
		Meiporo: m,
		Res:     NewResponseWriter(res),
		Req:     NewRequest(req),
		Params:  params,
		index:   0,
		Stash:   map[interface{}]interface{}{},
	}
	if m.CreateContextHandler != nil {
		m.CreateContextHandler(c)
	}
	if route != nil {
		c.Handlers = append(m.Middlewares, route.Dest.([]Handler)...)
	} else {
		if m.NotFoundHandler != nil {
			c.Handlers = append(m.Middlewares, m.NotFoundHandler)
		} else {
			c.Handlers = append(m.Middlewares, func(c *Context) {
				code := http.StatusNotFound
				http.Error(c.Res, http.StatusText(code), code)
			})
		}
	}
	return c
}

func (m *Meiporo) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c := m.CreateContext(res, req)
	c.Run()
}

func (m *Meiporo) Run() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	m.Logger.Printf("listening on %s:%s\n", host, port)
	m.Router.Start()
	http.ListenAndServe(host+":"+port, m)
}

func (m *Meiporo) Use(midd Middleware) {
	m.Middlewares = append(m.Middlewares, midd.Handler())
}
