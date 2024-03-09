package pages

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ostafen/clover/v2"
)

func GridPage(c *gin.Context, db *clover.DB) {
	c.HTML(http.StatusOK, "grid/index.html", nil)
}

func GridSearch(c *gin.Context, db *clover.DB) {
	c.HTML(http.StatusOK, "grid/index.html", nil)
}
