package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func MapPage(c *gin.Context, db *clover.DB) {
	c.HTML(http.StatusOK, "map/index.html", nil)
}
