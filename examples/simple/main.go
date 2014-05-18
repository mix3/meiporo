package main

import (
	"github.com/lestrrat/go-xslate"
	"github.com/mix3/meiporo"
)

func main() {
	m := meiporo.New()
	m.Router.Add("GET", "/xslate",
		func(c *meiporo.Context) {
			data := map[string]interface{}{
				"name": map[string]interface{}{
					"to": "mix3",
				},
			}
			c.Res.Render("xslate.tx", data)
		},
	)
	m.Router.Add("GET", "/error",
		func(c *meiporo.Context) {
			panic("panic")
		},
	)
	m.Router.Add("GET", "/:id",
		func(c *meiporo.Context) {
			c.Res.WriteText("%s", "Hello World!")
		},
	)
	m.Router.Add("GET", "/",
		func(c *meiporo.Context) {
			c.Meiporo.Logger.Printf("before 1")
			c.Next()
			c.Meiporo.Logger.Printf("before 3")
		},
		func(c *meiporo.Context) {
			c.Meiporo.Logger.Printf("before 2")
			c.Next()
			c.Meiporo.Logger.Printf("after 2")
		},
		func(c *meiporo.Context) {
			c.Meiporo.Logger.Printf("before 3")
			c.Res.WriteHTML("<html><body>Hello World!</body></html>")
			c.Meiporo.Logger.Printf("render")
			c.Meiporo.Logger.Printf("after 3")
		},
	)
	meiporo.Xslate.Configure(xslate.Args{
		"Loader": xslate.Args{
			"LoadPaths": []string{"tmpl"},
		},
	})
	m.Run()
}
