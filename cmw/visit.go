package cmw

import (
	"bytes"
	"io"
	"net/http"

	"github.com/btagrass/gobiz/sys/mdl"
	"github.com/btagrass/gobiz/sys/svc"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func Visit(visitSvc *svc.VisitSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		visit := mdl.Visit{
			UserId:    cast.ToInt64(c.GetFloat64("userId")),
			UserName:  c.GetString("userName"),
			Ip:        c.ClientIP(),
			Method:    c.Request.Method,
			Url:       c.Request.RequestURI,
			UserAgent: c.Request.UserAgent(),
		}
		if c.Request.Method != http.MethodGet {
			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logrus.Error(err)
			}
			len := lo.Min([]int{len(body), 1000})
			visit.Req = string(body[:len])
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		writer := visitResponseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer
		c.Next()
		body := writer.body.Bytes()
		len := lo.Min([]int{len(body), 1000})
		visit.Resp = string(body[:len])
		err := visitSvc.Save(visit)
		if err != nil {
			logrus.Error(err)
		}
	}
}

type visitResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w visitResponseWriter) Write(data []byte) (int, error) {
	_, err := w.body.Write(data)
	if err != nil {
		logrus.Error(err)
	}
	return w.ResponseWriter.Write(data)
}
