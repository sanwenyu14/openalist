package server

import (
	"context"
	"path"
	"strings"

	"github.com/AlliotTech/openalist/internal/conf"
	"github.com/AlliotTech/openalist/server/common"
	"github.com/AlliotTech/openalist/server/s3"
	"github.com/gin-gonic/gin"
)

func S3(g *gin.RouterGroup) {
	if !conf.Conf.S3.Enable {
		g.Any("/*path", func(c *gin.Context) {
			common.ErrorStrResp(c, "S3 server is not enabled", 403)
		})
		return
	}
	if conf.Conf.S3.Port != -1 {
		g.Any("/*path", func(c *gin.Context) {
			common.ErrorStrResp(c, "S3 server bound to single port", 403)
		})
		return
	}
	h, _ := s3.NewServer(context.Background())

	g.Any("/*path", func(c *gin.Context) {
		adjustedPath := strings.TrimPrefix(c.Request.URL.Path, path.Join(conf.URL.Path, "/s3"))
		c.Request.URL.Path = adjustedPath
		gin.WrapH(h)(c)
	})
}

func S3Server(g *gin.RouterGroup) {
	h, _ := s3.NewServer(context.Background())
	g.Any("/*path", gin.WrapH(h))
}
