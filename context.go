package meiporo

type Handler func(*Context)

type Context struct {
	Meiporo  *Meiporo
	Res      ResponseWriter
	Req      *Request
	Params   map[string]string
	Handlers []Handler
	index    int
	Stash    map[interface{}]interface{}
}

func (c *Context) Next() {
	if c.index < len(c.Handlers) {
		h := c.Handlers[c.index]
		c.index++
		h(c)
	}
}

func (c *Context) Run() {
	c.Next()
}
