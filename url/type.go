package url

import "github.com/gin-gonic/gin"

type PatternStruct struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}
