package middleware

import (
	"simple/framework"
)

func Recovery() framework.ControllerHandler {
	return func(c *framework.Context) error {

		defer func() {
			if err := recover(); err != nil {
				c.Json(err).SetStatus(500)
			}
		}()
		c.Next()
		return nil
	}
}
