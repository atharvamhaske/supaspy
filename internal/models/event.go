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
