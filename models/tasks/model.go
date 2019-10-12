package tasks

import (
	"time"
)

type EzbTasks struct {
	UUID       string    `json:"uuid"`
	CreateDate time.Time `json:"createddate"`
	UpdateDate time.Time `json:"updatedate"`
	Status     string    `json:"status"`
	TokenID    string    `json:"tokenid"`
	PID        int       `json:"pid"`
	Parameters string    `json:"parameters"`
	// StatusURL  string    `json:"statusurl"`
	// LogURL     string    `json:"logurl"`
	// ResultURL  string    `json:"resulturl"`
}

type taksStatus int

const (
	// PENDING: the job is created but not started
	PENDING taksStatus = iota
	RUNNING
	FAILED
	FINISH
)

func TaksStatus(i int) string {
	p := [4]string{"PENDING", "RUNNING", "FAILED", "FINISH"}
	return p[i]
}
