package data

import (
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"embed"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

//go:embed migrations/*.sql
var migrations embed.FS

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("sqlite3", "chitchat.db")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Fatal("failed to open db")
	}

	runMigrations()
}

func runMigrations() {
	dbDriver, err := sqlite3.WithInstance(Db, &sqlite3.Config{})
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("failed to create migration driver")
		os.Exit(1)
	}

	// https://github.com/golang-migrate/migrate/issues/471
	// temp solution would update to use embed source once merged
	srcDriver, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("failed to create source driver")
		os.Exit(1)
	}

	m, err := migrate.NewWithInstance("migrations", srcDriver, "sqlite3", dbDriver)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("failed to create migrate instance")
		os.Exit(1)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.WithFields(log.Fields{
			"Error": err.Error(),
		}).Fatal("failed to run database migration")
		os.Exit(1)
	}
}

// create a random UUID with from RFC 4122
// adapted from http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("Cannot generate UUID", err)
	}

	// 0x40 is reserved variant from RFC 4122
	u[8] = (u[8] | 0x40) & 0x7F
	// Set the four most significant bits (bits 12 through 15) of the
	// time_hi_and_version field to the 4-bit version number.
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// hash plaintext with SHA-1
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
