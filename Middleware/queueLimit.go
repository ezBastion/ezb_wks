// This file is part of ezBastion.

//     ezBastion is free software: you can redistribute it and/or modify
//     it under the terms of the GNU Affero General Public License as published by
//     the Free Software Foundation, either version 3 of the License, or
//     (at your option) any later version.

//     ezBastion is distributed in the hope that it will be useful,
//     but WITHOUT ANY WARRANTY; without even the implied warranty of
//     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//     GNU Affero General Public License for more details.

//     You should have received a copy of the GNU Affero General Public License
//     along with ezBastion.  If not, see <https://www.gnu.org/licenses/>.

package Middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/ezbastion/ezb_wks/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Limit(c *gin.Context) {
	logg := log.WithFields(log.Fields{
		"middleware": "queuelimit",
		// "xtrack":     x.String(),
	})
	t := runtime.NumGoroutine()
	fmt.Println("runtime ", t)
	conf, _ := c.MustGet("conf").(models.Configuration)
	lm := conf.LimitMax
	// 429 Too Many Requests
	if lm > 0 && t > lm {
		logg.Error("Too Many Requests ", t, " threads")
		c.Writer.Header().Set("X-ERROR", "L0001")
		c.AbortWithError(http.StatusTooManyRequests, errors.New("#L0001"))
		return
	}
	lw := conf.LimitWarning
	if lw > 0 && t > lw {
		logg.Warning("Heavy load ", t, " threads")
		c.Writer.Header().Set("X-ERROR", "L0002")
	}
	c.Next()
}
