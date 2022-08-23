package srv

import (
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func Start(addr string, s *gin.Engine) error {
	return endless.ListenAndServe(addr, s)
}
