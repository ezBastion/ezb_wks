package tasks

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/ezbastion/ezb_wks/models"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func GetResult(c *gin.Context) {
	tokenID := c.GetHeader("x-ezb-tokenid")
	xtrack := c.GetHeader("X-Track")
	logg := log.WithFields(log.Fields{
		"controller": "tasks",
		"xtrack":     xtrack,
	})
	logg.Debug("start GetResult")
	uuid := c.Param("UUID")
	conf, _ := c.MustGet("conf").(models.Configuration)
	taskPath := path.Join(strings.Replace(conf.JobPath, "\\", "/", -1), uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:])
	file := path.Join(taskPath, "output.json")
	if !checkTokenID(taskPath, tokenID) {
		logg.Error("log file not found")
		c.AbortWithError(http.StatusBadRequest, errors.New("#I0001"))
		return
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		logg.Error("result file not found")
		c.AbortWithError(http.StatusBadRequest, errors.New("#I0002"))
		return
	}

	raw, _ := ioutil.ReadFile(file)
	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, _ := ioutil.ReadAll(unicodeReader)
	c.Data(http.StatusOK, "application/json", decoded)

}

func GetStatus(c *gin.Context) {
	tokenID := c.GetHeader("x-ezb-tokenid")
	xtrack := c.GetHeader("X-Track")
	logg := log.WithFields(log.Fields{
		"controller": "tasks",
		"xtrack":     xtrack,
	})
	uuid := c.Param("UUID")
	logg.Debug("start GetStatus for uuid: ", uuid)
	conf, _ := c.MustGet("conf").(models.Configuration)
	taskPath := path.Join(strings.Replace(conf.JobPath, "\\", "/", -1), uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:])
	file := path.Join(taskPath, "status.json")
	if !checkTokenID(taskPath, tokenID) {
		logg.Error("log file not found")
		c.AbortWithError(http.StatusBadRequest, errors.New("#I0003"))
		return
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		logg.Error("status file not found")
		c.AbortWithError(http.StatusBadRequest, errors.New("#I0004"))
		return
	}

	raw, err := ioutil.ReadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error reading task status file")
		logg.Error("error reading task status file ", file)
		return
	}
	c.Data(http.StatusOK, "application/json", raw)
}

func GetLog(c *gin.Context) {
	tokenID := c.GetHeader("x-ezb-tokenid")
	xtrack := c.GetHeader("X-Track")
	logg := log.WithFields(log.Fields{
		"controller": "tasks",
		"xtrack":     xtrack,
	})
	logg.Debug("start GetLog")
	uuid := c.Param("UUID")
	conf, _ := c.MustGet("conf").(models.Configuration)
	taskPath := path.Join(strings.Replace(conf.JobPath, "\\", "/", -1), uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:])
	file := path.Join(taskPath, "trace.log")
	if !checkTokenID(taskPath, tokenID) {
		logg.Error("#I0005 log file not found")
		c.String(http.StatusBadRequest, "#I0005 log file not found")
		c.Abort()
		return
	}
	if _, err := os.Stat(file); os.IsNotExist(err) {
		logg.Error("#I0006 log file not found")
		c.String(http.StatusBadRequest, "#I0006 log file not found")
		c.Abort()
		return
	}

	raw, _ := ioutil.ReadFile(file)
	// Make an tranformer that converts MS-Win default to UTF8:
	win16be := unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	// Make a transformer that is like win16be, but abides by BOM:
	utf16bom := unicode.BOMOverride(win16be.NewDecoder())

	// Make a Reader that uses utf16bom:
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16bom)

	// decode and print:
	decoded, _ := ioutil.ReadAll(unicodeReader)
	c.Data(http.StatusOK, "text/plain", decoded)

}
func checkTokenID(taskPath, tokenID string) bool {
	file := path.Join(taskPath, "status.json")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}
	var status EzbTasks
	err = json.Unmarshal(raw, &status)
	if status.TokenID == tokenID {
		return true
	}
	return false
}
