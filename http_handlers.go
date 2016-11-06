package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

// POST /shorten
func createShortURL(c echo.Context) error {
	shorty := new(shorty)

	if err := c.Bind(shorty); err != nil {
		// Got malformed JSON data
		return c.JSON(http.StatusBadRequest, shortyError{Msg: "Bad Request"})
	}

	// Check the url param. Return 400 if missing
	if shorty.URL == nil {
		return c.JSON(http.StatusBadRequest, shortyError{Msg: "url is not present"})
	}

	if shorty.ShortCode != nil {

		// Check the validity of the specified shortcode
		if match, _ := regexp.MatchString("^[0-9a-zA-Z_]{4,}$", *shorty.ShortCode); match != true {

			// Have to specifically return status 422, as http.StatusUnprocessableEntity is not added yet to Golang :-(
			return c.JSON(422, shortyError{Msg: "The shortcode fails to meet the following regexp: ^[0-9a-zA-Z_]{4,}$"})
		}
	} else {

		// Generate new shortcode
		newShortCode := generateRandomString(6)
		shorty.ShortCode = &newShortCode
	}

	// Insert shortcode to DB
	rows, err := db.Query("INSERT INTO urls (url, shortcode) VALUES (?, ?)", shorty.URL, shorty.ShortCode)

	if err != nil {

		// Check for duplicate shortcode. Error 1062: Duplicate key
		if strings.HasPrefix(fmt.Sprint(err), "Error 1062:") {
			return c.JSON(http.StatusConflict, shortyError{Msg: "The the desired shortcode is already in use"})
		}
		return c.JSON(http.StatusInternalServerError, shortyError{Msg: "Internal Server Error"})
	} else {
		defer rows.Close()
	}

	// TODO: NO NO NO
	shorty.URL = nil
	return c.JSON(http.StatusCreated, shorty)
}

// GET /:shortcode
func redirectTo(c echo.Context) error {
	shortcode := c.Param("shortcode")
	shorty := new(shorty)

	query := "SELECT url, shortcode FROM urls WHERE shortcode = ?"
	err := db.QueryRow(query, shortcode).Scan(&shorty.URL, &shorty.ShortCode)

	if err != nil {
		return c.JSON(http.StatusNotFound, shortyError{Msg: "The shortcode cannot be found in the system"})
	}

	// Increment ShortCode redirectCount
	incrementQuery := "UPDATE urls SET redirects = redirects + 1 WHERE shortcode = ?"
	rows, err := db.Query(incrementQuery, shortcode)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, shortyError{Msg: "Internal Server Error"})
	} else {
		defer rows.Close()
	}

	return c.Redirect(http.StatusFound, *shorty.URL)
}

// GET /:shortcode/stats
func getShortcodeStats(c echo.Context) error {
	shortcode := c.Param("shortcode")
	shorty := new(shorty)

	query := "SELECT created_at, redirects, updated_at FROM urls WHERE shortcode = ?"
	err := db.QueryRow(query, shortcode).Scan(&shorty.StartDate, &shorty.RedirectCount, &shorty.LastSeenDate)

	if err != nil {
		return c.JSON(http.StatusNotFound, shortyError{Msg: "The shortcode cannot be found in the system"})
	}

	return c.JSON(http.StatusOK, shorty)
}
