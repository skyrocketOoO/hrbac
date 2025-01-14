package boot

import (
	"github.com/skyrocketOoO/serverx/docs/openapi"
	"github.com/skyrocketOoO/serverx/internal/global"
)

func InitSwagger() {
	openapi.SwaggerInfo.Title = "Swagger API"
	openapi.SwaggerInfo.Version = "1.0"
	openapi.SwaggerInfo.Host = ""
	openapi.SwaggerInfo.BasePath = global.ApiVersion
	openapi.SwaggerInfo.Schemes = []string{"http"}
}
