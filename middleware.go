package meiporo

import (
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Middleware interface {
	Handler() Handler
}

type LoggerMiddleware struct {
}

func (m *LoggerMiddleware) Handler() Handler {
	return func(c *Context) {
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime).Seconds()
			c.Meiporo.Logger.Printf(`[%.6f] %d "%s"`, duration, c.Res.Status(), c.Req.URL.Path)
		}()
		c.Next()
	}
}

type RecoverMiddleware struct {
	RecoverHandler func(*Context, interface{})
}

func (m *RecoverMiddleware) Handler() Handler {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if m.RecoverHandler != nil {
					m.RecoverHandler(c, err)
				} else {
					code := http.StatusInternalServerError
					http.Error(c.Res, http.StatusText(code), code)
				}
			}
		}()
		c.Next()
	}
}

type StaticFileMiddleware struct {
	PublicPath string
}

func (m *StaticFileMiddleware) Handler() Handler {
	return func(c *Context) {
		path := filepath.Join(m.PublicPath, c.Req.URL.Path)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			c.Res.Header().Del("Content-Type")
			http.ServeFile(c.Res, c.Req.Request, path)
			return
		}
		c.Next()
	}
}
