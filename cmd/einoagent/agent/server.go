package agent

import (
	"embed"

	"github.com/cloudwego/hertz/pkg/route"
)

//go:embed web
var webContent embed.FS

type ChatRequest struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func BindRoutes(r *route.RouterGroup) error {
	if err := Init(); err != nil {
		return err
	}
	return nil
}
