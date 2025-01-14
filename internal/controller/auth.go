package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/skyrocketOoO/erx/erx"
	"github.com/skyrocketOoO/go-utils/auth"
	ope "github.com/skyrocketOoO/gorm-plugin/lib/operator"
	wh "github.com/skyrocketOoO/gorm-plugin/lib/where"
	cm "github.com/skyrocketOoO/serverx/internal/common"
	col "github.com/skyrocketOoO/serverx/internal/gen/column"
	tb "github.com/skyrocketOoO/serverx/internal/gen/table"
	"github.com/skyrocketOoO/serverx/internal/global"
	dm "github.com/skyrocketOoO/serverx/internal/global/domain"
	"github.com/skyrocketOoO/serverx/internal/model"
	"gorm.io/gorm"
)

// @Param   user  body  controller.Login.Req  true  "Login User"
// @Success 200 {object} controller.Login.Resp "token"
// @Failure 500 {string} dm.ErrResp 
// @Failure 400 {object} dm.ErrResp 
// @Router /login [post]
func (h *Handler) Login(c *gin.Context) {
	type Req struct {
		Name     string `json:"Name" validate:"required"`
		Password string `json:"Password" validate:"required"`
	}

	var req Req
	if ok := cm.BindAndValidate(c, &req); !ok {
		return
	}

	hashedPassword := string(auth.Hash(req.Password, cm.GetSalt()))
	db := global.DB

	var user model.User
	if err := db.
		Table(tb.Users).
		Where(wh.B(col.Users.Name, ope.Eq), req.Name).
		Where(wh.B(col.Users.Password, ope.Eq), hashedPassword).
		Take(&user).
		Error; err != nil {
		dm.RespErr(c, http.StatusInternalServerError, erx.W(dm.ErrLoginFailed))
		return
	}

	token, err := cm.GenerateToken(user.ID)
	if err != nil {
		dm.RespErr(c, http.StatusInternalServerError, erx.W(err))
		return
	}

	token = "Bearer " + token

	type Resp struct {
		Token string `json:"token"`
	}
	c.JSON(http.StatusOK, Resp{Token: token})
}

// @Param   user  body  controller.Register.Req  true  "Register"
// @Success 200
// @Failure 500 {object} dm.ErrResp
// @Failure 400 {object} dm.ErrResp
// @Router /register [post]
func (h *Handler) Register(c *gin.Context) {
	type Req struct {
		Name     string `json:"Name" validate:"required,min=6,max=32"`
		Password string `json:"Password" validate:"required,min=8,max=32"`
	}

	var req Req
	if ok := cm.BindAndValidate(c, &req); !ok {
		return
	}

	db := global.DB
	var existingUser model.User
	if err := db.Where(wh.B(col.Users.Name, ope.Eq), req.Name).
		Take(&existingUser).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			dm.RespErr(c, http.StatusInternalServerError, erx.W(err))
			return
		}
	} else {
		dm.RespErr(c, http.StatusInternalServerError, erx.W(dm.ErrUserNameRepetite))
		return
	}

	if err := db.Create(&model.User{
		Name:     req.Name,
		Password: string(auth.Hash(req.Password, cm.GetSalt())),
	}).Error; err != nil {
		dm.RespErr(c, http.StatusInternalServerError, erx.W(err))
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) ForgetPassword(c *gin.Context) {
	dm.RespErr(c, http.StatusNotFound, erx.W(dm.ErrNotImplement))
	return
}
