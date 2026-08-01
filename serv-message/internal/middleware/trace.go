package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/reshap0318/serv-message/internal/helpers"
)

// TraceID reads the trace id set by the api-gateway and stores it in the
// request context so every log line for this request can be correlated
// across services. Runs before GatewayAuth so even 401 responses carry one.
// If the header is missing (e.g. a request that didn't come through the
// gateway), a new trace id is generated so logs stay correlatable.
func TraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader(helpers.TraceIDHeader)
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request.Context(), helpers.KeyTraceID, traceID)
		c.Request = c.Request.WithContext(ctx)

		c.Set("trace_id", traceID)
		c.Header(helpers.TraceIDHeader, traceID)

		c.Next()
	}
}
