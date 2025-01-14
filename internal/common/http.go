package cm

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/erx/erx"
	"github.com/skyrocketOoO/serverx/internal/global"
	dm "github.com/skyrocketOoO/serverx/internal/global/domain"
)

func BindAndValidate[T any](c *gin.Context, req *T) bool {
	if c.Request.Body == nil || c.Request.ContentLength == 0 {
		dm.RespErr(c, http.StatusBadRequest, erx.W(dm.ErrEmptyRequest))
		return false
	}

	if err := c.ShouldBindJSON(req); err != nil {
		dm.RespErr(c, http.StatusBadRequest, erx.W(err))
		return false
	}

	if err := global.Validator.Struct(req); err != nil {
		dm.RespErr(c, http.StatusBadRequest, erx.W(err))
		return false
	}
	return true
}

func RespErr(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, dm.ErrResp{Error: err.Error()})
}
