package models

import (
	"database/sql"
	// "fmt"
	// "log"
	// "os"
	"time"

	"server/types"
)

var db *sql.DB
const dbTimeout = 5 * time.Second

type Models struct {
	Users User
	Questions Question
	JsonResponse types.JsonResponse
}

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{}
}
