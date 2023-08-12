package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"math/big"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/itchyny/base58-go"
	"github.com/labstack/echo/v4"
)

type ConnectionDetails struct {
	host     string // The hostname (default "localhost")
	port     string // The port to connect on (default "5432")
	user     string // The username (default "root")
	password string // The password (default "root")
	db_name  string // The database name (default "mysql")
}

type URLEntry struct {
	ID           string
	Redirect_URL string
	Created_At   time.Time
	Updated_At   time.Time
}

type urlRequestBody struct {
	URL string `json:"url"`
}

type successResponse struct {
	TS  int64  `json:"ts"`
	URL string `json:"url"`
}

/*
Establishes a connection with the MySQL instance.
*/
func NewConnection(conn_details ConnectionDetails) (*sql.DB, error) {
	// Connect to MySQL database
	db_info := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		conn_details.user, conn_details.password, conn_details.host, conn_details.port, conn_details.db_name)

	db, err := sql.Open("mysql", db_info)
	if err != nil {
		return nil, err
	}

	return db, nil
}
func InsertURL(db *sql.DB, c context.Context, entry URLEntry) error {
	insert_query := `
		INSERT INTO url_map (id, redirect_url, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	select_query := `
			SELECT * 
			FROM url_map
			WHERE id = ?
	`

	update_timestamp := `
			UPDATE url_map
			SET updated_at = ?
			WHERE id = ?
	`

	_, err := db.ExecContext(c, `CREATE TABLE IF NOT EXISTS url_map (
		id char(8) NOT NULL,
		redirect_url varchar(150) NOT NULL UNIQUE,
		created_at timestamp NOT NULL,
		updated_at timestamp NOT NULL,
		PRIMARY KEY(id)
	);`)
	if err != nil {
		return err
	}

	res, err := db.QueryContext(c, select_query, entry.ID) // See if the URL already exists
	if err != nil {
		return err
	}
	defer res.Close()
	if res.Next() { // If we get rows back, that means we have a duplicate URL
		var createdAt, updatedAt []byte
		err := res.Scan(&entry.ID, &entry.Redirect_URL, &createdAt, &updatedAt)
		if err != nil {
			return err
		}

		entry.Created_At, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
		if err != nil {
			return err
		}
		entry.Updated_At, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
		if err != nil {
			return err
		}

		_, err = db.ExecContext(c, update_timestamp, time.Now(), entry.ID)
		if err != nil {
			return err
		}

		return nil
	}

	_, err = db.ExecContext(c, insert_query, entry.ID, entry.Redirect_URL, entry.Created_At, entry.Updated_At)
	if err != nil {
		return err
	}

	return nil
}

func GetURL(c echo.Context) error {
	Database, _ := NewConnection(ConnectionDetails{
		host:     "localhost",
		port:     "3306",
		user:     "user",
		password: "password",
		db_name:  "shorturl_db",
	})
	defer Database.Close()
	err := Database.PingContext(c.Request().Context())
	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not connect to MySQL.")
	}
	id := c.Param("param")

	find_url_query := `
		SELECT * 
		FROM url_map
		WHERE id = ?
	`
	res, err := Database.QueryContext(c.Request().Context(), find_url_query, id) // Attempt to find a matching URL

	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error encountered while attempting to lookup URL.")
	}
	defer res.Close()
	if !res.Next() {
		return echo.NewHTTPError(http.StatusNotFound, "Invalid URL ID.")
	}

	var entry URLEntry
	var createdAt, updatedAt []uint8
	err = res.Scan(&entry.ID, &entry.Redirect_URL, &createdAt, &updatedAt)
	if err != nil {
		Logger.Fatal(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not successfully scan retrieved DB entry.")
	}
	entry.Created_At, err = time.Parse("2006-01-02 15:04:05", string(createdAt))
	if err != nil {
		return err
	}
	entry.Updated_At, err = time.Parse("2006-01-02 15:04:05", string(updatedAt))
	if err != nil {
		return err
	}

	fmt.Println("Redirecting to URL: ", entry) // Will Print the URL if GET call is successful
	c.Redirect(http.StatusPermanentRedirect, entry.Redirect_URL)
	return nil
}

func UpdateURL(c echo.Context) error {
	Database, _ := NewConnection(ConnectionDetails{
		host:     "localhost",
		port:     "3306",
		user:     "user",
		password: "password",
		db_name:  "shorturl_db",
	})
	defer Database.Close()

	req_body := new(urlRequestBody)

	err := c.Bind(req_body)
	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request.")
	}

	u := req_body.URL

	if u == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing URL parameter")
	}

	id := c.Param("param")
	entry := URLEntry{
		ID:           id,
		Updated_At:   time.Now(),
		Redirect_URL: u,
	}

	update_query := `
		UPDATE url_map
		SET redirect_url = ?, updated_at = ?
		WHERE id = ?
	`

	_, err = Database.ExecContext(c.Request().Context(), update_query, entry.Redirect_URL, entry.Updated_At, entry.ID)
	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not update URL.")
	}

	return c.NoContent(http.StatusOK)
}

func DeleteURL(c echo.Context) error {
	Database, _ := NewConnection(ConnectionDetails{
		host:     "localhost",
		port:     "3306",
		user:     "user",
		password: "password",
		db_name:  "shorturl_db",
	})
	defer Database.Close()
	err := Database.PingContext(c.Request().Context())
	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not connect to MySQL.")
	}
	id := c.Param("param")

	sqlStatement := `
	DELETE FROM url_map
	WHERE id = ?;`
	res, err := Database.ExecContext(c.Request().Context(), sqlStatement, id)
	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "Error encountered while attempting to delete URL.")
	}

	rowsAffected, _ := res.RowsAffected()
	return c.JSON(http.StatusOK, map[string]int64{
		"rows_affected": rowsAffected,
	})

	// Rest of your function
}

// GenerateShortLink, sha256Of, and base58Encoded functions remain the same.
func PutURL(c echo.Context) error {

	Database, _ = NewConnection(ConnectionDetails{
		host:     "localhost",
		port:     "3306",
		user:     "user",
		password: "password",
		db_name:  "shorturl_db",
	})

	fmt.Println("connection established")
	// err := Database.PingContext(c.Request().Context())
	// if err != nil {
	// 	Logger.Error(err.Error())
	// 	return echo.NewHTTPError(http.StatusInternalServerError, "Could not connect to MySQL.")
	// }

	// Database.

	req_body := new(urlRequestBody)

	err := c.Bind(req_body)
	if err != nil {
		fmt.Println(req_body)
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Failed to decode request.")
	}
	u := req_body.URL

	if u == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing URL parameter")
	}

	id := GenerateShortLink(u)
	t := time.Now()

	err = InsertURL(Database, c.Request().Context(), URLEntry{
		ID:           id,
		Created_At:   t,
		Updated_At:   t,
		Redirect_URL: u,
	})

	if err != nil {
		Logger.Error(err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Unable to shorten URL: %s", err.Error()))
	}

	return c.JSON(http.StatusOK, &successResponse{
		TS:  t.UnixNano(),
		URL: "http://localhost:8082/" + id,
	})
}

func GenerateShortLink(initialLink string) string {
	urlHashBytes := sha256Of(initialLink)
	generatedNumber := new(big.Int).SetBytes(urlHashBytes).Uint64()
	finalString := base58Encoded([]byte(fmt.Sprintf("%d", generatedNumber)))
	return finalString[:8]
}

func sha256Of(input string) []byte {
	algorithm := sha256.New()
	algorithm.Write([]byte(input))
	return algorithm.Sum(nil)
}

func base58Encoded(bytes []byte) string {
	encoding := base58.BitcoinEncoding
	encoded, _ := encoding.Encode(bytes)
	return string(encoded)
}
