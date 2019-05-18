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

package wkslog

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
func getXtrack(c *gin.Context) {
	xtrack := c.Param("id")
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	// t := time.Now().UTC()
	// l := fmt.Sprintf("log/ezb_wks-%d%d.log", t.Year(), t.YearDay())
	var logfile []string
	err := filepath.Walk(path.Join(exPath, "/log"), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			b, err := ioutil.ReadFile(path)
			if err == nil {
				s := string(b)
				if strings.Contains(s, xtrack) {
					f, err := os.Open(path)
					defer f.Close()
					if err == nil {
						scanner := bufio.NewScanner(f)
						line := 1
						for scanner.Scan() {
							if strings.Contains(scanner.Text(), xtrack) {
								logfile = append(logfile, scanner.Text())
							}
							line++
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusNoContent, err.Error())
	}
	if len(logfile) == 0 {
		c.JSON(http.StatusNoContent, "x-track not found")
	}

	c.JSON(http.StatusOK, logfile)

}
func getLast(c *gin.Context) {
	requestLogger := log.WithFields(log.Fields{"request_id": "request_id", "user_ip": "user_ip"})
	requestLogger.Info("request done!")
	c.JSON(http.StatusOK, "ret")

}
