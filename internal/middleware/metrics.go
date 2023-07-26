package middleware

// func (m *Middleware) MetricsMiddleware(metrics metric.Metrics) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		start := time.Now()
// 		err := c.next(c)
// 		var status int
// 		if err != nil {
// 			status = err.(*echo.HTTPError).Code
// 		} else {
// 			status = c.Response().Status
// 		}
// 		metrics.ObserveResponseTime(status, c.Request().Method, c.Path(), time.Since(start).Seconds())
// 		metrics.IncHits(status, c.Request().Method, c.Path())
// 		return err
// 	}

// }
