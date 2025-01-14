package cm

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	ope "github.com/skyrocketOoO/gorm-plugin/lib/operator"
	wh "github.com/skyrocketOoO/gorm-plugin/lib/where"
	col "github.com/skyrocketOoO/serverx/internal/gen/column"
	"github.com/skyrocketOoO/serverx/internal/global"
	"github.com/skyrocketOoO/serverx/internal/model"
	"github.com/spf13/viper"
)

func GetSalt() []byte {
	return []byte(viper.GetString("salt"))
}

func GetJwtSecretKey() []byte {
	return []byte(viper.GetString("jwt-secret-key"))
}

func GenerateToken(userID string) (string, error) {
	// Set token claims
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(10000 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(GetJwtSecretKey())
}

func GetOperator(c *gin.Context) (user model.User, err error) {
	operatorName, ok := c.Get("username")
	if !ok {
		return model.User{}, errors.New("username not set in jwt")
	}

	db := global.DB
	if err = db.
		Where(wh.B(col.Users.Name, ope.Eq), operatorName).
		Take(&user).Error; err != nil {
		return
	}

	return user, nil
}
