package dataaccess

import (
	"database/sql"
	"fmt"
	"../types"
	_"github.com/lib/pq"
	"time"
	//log "github.com/sirupsen/logrus"
)

type connection struct {
	DB *sql.DB
}

// GetConnection opens a new db connection
func GetConnection(connectionString string) (*connection, error) {
	var connection connection
	var err error

	connection.DB, err = sql.Open("postgres", connectionString)

	return &connection, err
}

// GetServerList reads all servers with scheduled backups for current date
func (connection connection) GetServerList() ([]string, error) {
	date := fmt.Sprintf("%s", time.Now().Format("2006-01-02"))
        //log.Info("Date: ", date)
	results, err := connection.DB.Query("SELECT DISTINCT Name FROM job WHERE SchedTime::timestamp::date = $1", date)

	if err != nil {
		return nil, err
	}

	var servers []string

	for results.Next() {
		var server string
		err = results.Scan(&server)
		servers = append(servers, server)
	}

	return servers, err
}

// TotalBytes returns total bytes saved for a server since the very first backup
func (connection connection) TotalBytes(server string) (*types.TotalBytes, error) {
	results, err := connection.DB.Query("SELECT SUM(JobBytes) FROM job WHERE Name=$1", server)

	if err != nil {
		return nil, err
	}

	var totalBytes types.TotalBytes
	if results.Next() {
		err = results.Scan(&totalBytes.Bytes)
		results.Close()
	}

	return &totalBytes, err
}

// TotalFiles returns total files saved for a server since the very first backup
func (connection connection) TotalFiles(server string) (*types.TotalFiles, error) {
	results, err := connection.DB.Query("SELECT SUM(JobFiles) FROM job WHERE Name=$1", server)

	if err != nil {
		return nil, err
	}

	var totalFiles types.TotalFiles
	if results.Next() {
		err = results.Scan(&totalFiles.Files)
		results.Close()
	}

	return &totalFiles, err
}

// LastJob returns metrics for latest executed server backup
func (connection connection) LastJob(server string) (*types.LastJob, error) {
	results, err := connection.DB.Query("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE $1 ORDER BY StartTime DESC LIMIT 1", server)

	if err != nil {
		return nil, err
	}

	var lastJob types.LastJob
	if results.Next() {
		err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)
		results.Close()
	}

	return &lastJob, err
}

// LastJob returns metrics for latest executed server backup with Level F
func (connection connection) LastFullJob(server string) (*types.LastJob, error) {
	results, err := connection.DB.Query("SELECT Level,JobBytes,JobFiles,JobErrors,DATE(StartTime) AS JobDate FROM job WHERE Name LIKE $1 AND Level = 'F' ORDER BY StartTime DESC LIMIT 1", server)

	if err != nil {
		return nil, err
	}

	var lastJob types.LastJob
	if results.Next() {
		err = results.Scan(&lastJob.Level, &lastJob.JobBytes, &lastJob.JobFiles, &lastJob.JobErrors, &lastJob.JobDate)
		results.Close()
	}

	return &lastJob, err
}

// ScheduledTime returns amount of scheduled jobs
func (connection connection) ScheduledJobs(server string) (*types.ScheduledJob, error) {
	date := fmt.Sprintf("%s%%", time.Now().Format("2006-01-02"))
	results, err := connection.DB.Query("SELECT COUNT(DATE(SchedTime)) AS JobsScheduled FROM job WHERE Name LIKE $1 AND SchedTime >= $2", server, date)

	if err != nil {
		return nil, err
	}

	var schedJob types.ScheduledJob
	if results.Next() {
		err = results.Scan(&schedJob.ScheduledJobs)
		results.Close()
	}

	return &schedJob, err
}
