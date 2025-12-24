package agent

import (
	"embed"

	"github.com/cloudwego/hertz/pkg/route"
)

//go:embed web
var webContent embed.FS

func BindRoutes(r *route.RouterGroup) error {
	return nil
}
