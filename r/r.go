package r

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type R struct {
	Code any    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func Q(c *gin.Context) map[string]any {
	queries := make(map[string]any)
	for k := range c.Request.URL.Query() {
		queries[k] = c.Query(k)
	}
	return queries
}

func J(c *gin.Context, data ...any) {
	var r R
	err, ok := data[len(data)-1].(error)
	if ok {
		r.Code = http.StatusInternalServerError
		r.Msg = err.Error()
	} else {
		r.Code = http.StatusOK
		if len(data) == 1 {
			r.Data = data[0]
		} else {
			count := cast.ToInt64(data[1])
			if count == 0 {
				r.Data = data[0]
			} else {
				r.Data = map[string]any{
					"records": data[0],
					"total":   count,
				}
			}
		}
	}
	c.JSON(http.StatusOK, r)
	c.Abort()
}
