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

package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"ezb_wks/Middleware"
	"ezb_wks/models"
	"ezb_wks/models/exec"
	"ezb_wks/models/healthCheck"
	"ezb_wks/models/wkslog"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"time"

	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

var conf models.Configuration
var exPath string

func mainGin(serverchan *chan bool) {
	ex, _ := os.Executable()
	exPath = filepath.Dir(ex)

	err := gonfig.GetConf(path.Join(exPath, "/conf/config.json"), &conf)
	if err != nil {
		panic(err)
	}


	/* log */
	outlog := true
	gin.DisableConsoleColor()
	log.SetFormatter(&log.JSONFormatter{})
	switch conf.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "warning":
		log.SetLevel(log.WarnLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "critical":
		log.SetLevel(log.FatalLevel)
		break
	default:
		outlog = false
	}
	if outlog {
		if _, err := os.Stat(path.Join(exPath, "log")); os.IsNotExist(err) {
			err = os.MkdirAll(path.Join(exPath, "log"), 0600)
			if err != nil {
				log.Println(err)
			}
		}

		ti := time.NewTicker(1 * time.Minute)
		defer ti.Stop()
		go func() {
			for range ti.C {
				t := time.Now().UTC()
				l := fmt.Sprintf("log/ezb_wks-%d%d.log", t.Year(), t.YearDay())
				f, _ := os.OpenFile(path.Join(exPath, l), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				defer f.Close()
				log.SetOutput(io.MultiWriter(f))
			}
		}()
	}
	/* log */

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(ginrus.Ginrus(log.StandardLogger(), time.RFC3339, true))
	r.Use(Middleware.ConfigMiddleware(conf))
	r.Use(Middleware.Limit)

	healthCheck.Routes(r)
	wkslog.Routes(r)
	exec.Routes(r)
	caCert, err := ioutil.ReadFile(path.Join(exPath, conf.CaCert))
	if err != nil {
		log.Fatal(err)

	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	tlsConfigPKI := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS12,
	}
	tlsConfigPKI.BuildNameToCertificate()
	srv := &http.Server{
		Addr:      conf.Listen,
		TLSConfig: tlsConfigPKI,
		Handler:   r,
	}

	go func() {
		if err := srv.ListenAndServeTLS(path.Join(exPath, conf.PublicCert), path.Join(exPath, conf.PrivateKey)); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
