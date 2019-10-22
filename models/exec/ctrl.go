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
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ezbastion/ezb_wks/models"
	"github.com/ezbastion/ezb_wks/models/tasks"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

var (
	conf   models.Configuration
	xtrack string
)

func run(c *gin.Context) {
	conf, _ = c.MustGet("conf").(models.Configuration)
	polling := c.GetHeader("X-Polling")
	xtrack = c.GetHeader("X-Track")
	var params EzbParams
	err := c.ShouldBindJSON(&params)
	if err != nil {
		c.String(http.StatusInternalServerError, "#E0002 bind parameters error", err)
		return
	}
	psParams := fmt.Sprintf("-xtrack '%s' ", xtrack)
	for i, h := range params.Data {
		psParams = fmt.Sprintf("%s -%s '%s' ", psParams, i, h)
	}
	psscript := filepath.Join(conf.ScriptPath, params.Meta.Job.Path)
	if polling == "true" {
		runTask(c, psscript, psParams)
	} else {
		runJob(c, psscript, psParams)
	}
}

func runJob(c *gin.Context, psscript string, psParams string) {
	logg := log.WithFields(log.Fields{
		"controller": "exec",
		"xtrack":     xtrack,
	})

	logg.Debug("start")
	cmd := exec.Command("powershell", "-NoLogo", "-NonInteractive", "-Command", "&{", psscript, " ", psParams, "}")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Run()
	if stderr.Len() != 0 {
		errStr := stderr.String()
		logg.Error("#E0001", errStr)
		c.String(http.StatusInternalServerError, "#E0001 Powershell error: %s", errStr)
		return
	}
	if stdout.Len() != 0 {
		ret := json.RawMessage(stdout.Bytes())
		c.JSON(http.StatusOK, ret)
	} else {
		c.JSON(http.StatusNoContent, "")
	}
}

func runTask(c *gin.Context, psscript string, psParams string) {
	tokenID := c.GetHeader("x-ezb-tokenid")
	logg := log.WithFields(log.Fields{
		"controller": "exec",
		"xtrack":     xtrack,
	})

	t := time.Now()
	taskID := fmt.Sprintf("%s%s", t.Format("20060102"), xtrack)
	jobPath := path.Join(strings.Replace(conf.JobPath, "\\", "/", -1), t.Format("2006/01/02"), xtrack)
	if _, err := os.Stat(jobPath); os.IsNotExist(err) {
		err = os.MkdirAll(jobPath, 0600)
		if err != nil {
			logg.Error("CANNOT CREATE JOB FOLDER")
			c.JSON(http.StatusInternalServerError, "CANNOT CREATE JOB FOLDER")
			return
		}
	}

	stdOUT := filepath.Join(jobPath, "output.json")
	stdTrace := filepath.Join(jobPath, "trace.log")
	statusFile := filepath.Join(jobPath, "status.json")
	task := tasks.EzbTasks{}
	task.UUID = taskID
	task.TokenID = tokenID
	task.CreateDate = time.Now()
	task.UpdateDate = time.Now()
	task.Parameters = psParams
	cmd := exec.Command("powershell", "-NoLogo", "-NonInteractive", "-Command", "&{", psscript, " ", psParams, "} 1>", stdOUT, " *>", stdTrace)
	cmd.Start()
	task.PID = cmd.Process.Pid
	task.Status = tasks.TaksStatus(int(tasks.RUNNING))
	c.JSON(http.StatusOK, task)
	go waitTask(cmd, &task, statusFile)
}

func waitTask(cmd *exec.Cmd, task *tasks.EzbTasks, statusFile string) {
	ta, _ := json.Marshal(task)
	ioutil.WriteFile(statusFile, ta, 0600)
	cmd.Wait()
	task.Status = tasks.TaksStatus(int(tasks.FINISH))
	task.UpdateDate = time.Now()
	ta, _ = json.Marshal(task)
	ioutil.WriteFile(statusFile, ta, 0600)
}
