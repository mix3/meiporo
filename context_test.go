package meiporo

import "testing"

func TestContext(t *testing.T) {
	func() {
		count := 0
		countUpHandler := func(c *Context) {
			count++
			c.Next()
		}
		c := &Context{
			Handlers: []Handler{
				countUpHandler,
				countUpHandler,
				countUpHandler,
				countUpHandler,
				countUpHandler,
			},
		}
		c.Run()
		expect(t, count, 5)
	}()

	func() {
		count := 0
		countUpHandler := func(c *Context) {
			count++
			c.Next()
			count--
		}
		c := &Context{
			Handlers: []Handler{
				countUpHandler,
				countUpHandler,
				countUpHandler,
				countUpHandler,
				countUpHandler,
			},
		}
		c.Run()
		expect(t, count, 0)
	}()
}
