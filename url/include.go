package url

import (
	"bjtubox/app/user"
	"github.com/gin-gonic/gin"
)

func Include(group *gin.RouterGroup, pattern []user.UrlPatternStruct, middlewareList []gin.HandlerFunc) {
	type HttpMethod func(string, ...gin.HandlerFunc) gin.IRoutes
	methods := map[string]HttpMethod{
		"GET":     group.GET,
		"POST":    group.POST,
		"PATCH":   group.PATCH,
		"DELETE":  group.DELETE,
		"PUT":     group.PUT,
		"HEAD":    group.HEAD,
		"OPTIONS": group.OPTIONS,
		"Any":     group.Any,
	}
	if middlewareList != nil {
		for _, middleware := range middlewareList {
			group.Use(middleware)
		}
	}
	for _, pattern := range pattern {
		if _, ok := methods[pattern.Method]; ok {
			methods[pattern.Method](pattern.Path, pattern.Handler)
		}
	}
}
