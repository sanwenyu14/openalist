package handles

import (
	"path"
	"strings"

	"github.com/AlliotTech/openalist/internal/errs"
	"github.com/AlliotTech/openalist/internal/model"
	"github.com/AlliotTech/openalist/internal/op"
	"github.com/AlliotTech/openalist/internal/search"
	"github.com/AlliotTech/openalist/pkg/utils"
	"github.com/AlliotTech/openalist/server/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type SearchReq struct {
	model.SearchReq
	Password string `json:"password"`
}

type SearchResp struct {
	model.SearchNode
	Type int `json:"type"`
}

func Search(c *gin.Context) {
	var (
		req SearchReq
		err error
	)
	if err = c.ShouldBind(&req); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	user := c.MustGet("user").(*model.User)
	req.Parent, err = user.JoinPath(req.Parent)
	if err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	if err := req.Validate(); err != nil {
		common.ErrorResp(c, err, 400)
		return
	}
	nodes, total, err := search.Search(c, req.SearchReq)
	if err != nil {
		common.ErrorResp(c, err, 500)
		return
	}
	var filteredNodes []model.SearchNode
	for _, node := range nodes {
		if !strings.HasPrefix(node.Parent, user.BasePath) {
			continue
		}
		meta, err := op.GetNearestMeta(node.Parent)
		if err != nil && !errors.Is(errors.Cause(err), errs.MetaNotFound) {
			continue
		}
		if !common.CanAccess(user, meta, path.Join(node.Parent, node.Name), req.Password) {
			continue
		}
		filteredNodes = append(filteredNodes, node)
	}
	common.SuccessResp(c, common.PageResp{
		Content: utils.MustSliceConvert(filteredNodes, nodeToSearchResp),
		Total:   total,
	})
}

func nodeToSearchResp(node model.SearchNode) SearchResp {
	return SearchResp{
		SearchNode: node,
		Type:       utils.GetObjType(node.Name, node.IsDir),
	}
}
