package healthing

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Status(c *gin.Context) {
	c.String(http.StatusOK, "ok\n")
}

func HttpStatus(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "ok\n")
}
