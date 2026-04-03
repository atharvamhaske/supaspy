package models

import "time"

// Here Operation means the SQL command type
type Operation string

const (
	OpSelect Operation = "SELECT"
	OpInsert Operation = "INSERT"
	OpUpdate Operation = "UPDATE"
	OpDelete Operation = "DELETE"
)

// QueryEvent represents a single DB query execution captured from the system. we are treating every query as a query event with following fields
type QueryEvent struct {
	ID        string
	Query     string
	Duration  int
	Error     string
	Timestamp time.Time
	Table     string
	Operation Operation
}

// Severity represents how crictical an alert is
type Severity string

const (
	SeverityInfo     Severity = "INFO"
	SeverityWarning  Severity = "WARNING"
	SeverityCritical Severity = "CRITICAL"
)

// Alert represents and anomaly detetced by an analyzer
type Alert struct {
	Title     string
	Message   string
	Severity  Severity
	Event     QueryEvent
	Timestamp time.Time
}

// this is to be done from tomorrow so for reminder i am commenting this pls bear lol
