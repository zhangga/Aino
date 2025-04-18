package service

import "github.com/gin-gonic/gin"

var (
	getHandlers  = make(map[string]gin.HandlerFunc)
	postHandlers = make(map[string]gin.HandlerFunc)
)

func ginHandlers(r *gin.Engine) {
	for path, handler := range getHandlers {
		r.GET(path, handler)
	}
	for path, handler := range postHandlers {
		r.POST(path, handler)
	}
}

func registerGetHandler(name string, handler gin.HandlerFunc) {
	getHandlers[name] = handler
}

func registerPostHandler(name string, handler gin.HandlerFunc) {
	postHandlers[name] = handler
}
