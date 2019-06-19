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

package exec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ezbastion/ezb_wks/models"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
func run(c *gin.Context) {
	// name := c.Param("name")
	// db read script by name
	// checksum
	// asynv / sync
	logg := log.WithFields(log.Fields{
		"controller": "exec",
		// "xtrack":     x.String(),
	})
	conf, _ := c.MustGet("conf").(models.Configuration)
	// var params interface{}
	xtrack := c.GetHeader("X-Track")
	params := make(map[string]string)
	psParams := fmt.Sprintf("-xtrack '%s' ", xtrack)
	err := c.ShouldBindJSON(&params)
	if err != nil {
		logg.Error(err)
	}

	for i, h := range params {
		psParams = fmt.Sprintf("%s -%s '%s' ", psParams, i, h)

	}
	// fmt.Println(psParams)
	var ezbjob models.EzbJobs
	js := strings.NewReader(params["job"])
	json.NewDecoder(js).Decode(&ezbjob)
	fmt.Println(ezbjob.Path)
	psscript := filepath.Join(conf.ScriptPath, ezbjob.Path)

	cmd := exec.Command("powershell", "-NoLogo", "-NonInteractive", "-Command", "&{", psscript, " ", psParams, "}")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()

	if stderr.Len() != 0 {
		errStr := stderr.String()
		log.Printf("Runnning %s  failed with err: \n***************\n%s\n***************\n", psscript, errStr)
		c.JSON(http.StatusInternalServerError, errStr)
	} else {
		ret := json.RawMessage(stdout.Bytes())
		c.JSON(http.StatusOK, ret)
	}
}

func runTaks(c *gin.Context) {

}
