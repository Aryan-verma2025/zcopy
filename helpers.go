package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
)

func clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func authenticate(token string, utype string) (string, error) {
	claims, err := verifyToken(token, utype)

	if !err {
		return "", errors.New("invalid token")
	}

	user := fmt.Sprint(claims["username"])

	return user, nil
}

func (app *application) idFromUsername(username string) (int, bool) {

	var id int

	stmt := "SELECT id FROM users WHERE username=?;"

	row := app.db.QueryRow(stmt, username)

	err := row.Scan(&id)

	if err != nil {
		app.errLog.Println(err)
		return 0, false
	}

	return id, true

}

func logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("set-cookie", "usr=")
}

func (app *application) generateURL() string {

	findStmt := "SELECT url FROM url WHERE url=?;"

	retrievedURL := ""
re:

	generatedURL := RandStringBytesRmndr(10)

	row := app.db.QueryRow(findStmt, generatedURL)

	if err := row.Scan(&retrievedURL); err != nil {
		app.errLog.Println(err)
	}

	if generatedURL == retrievedURL {
		goto re
	}

	return generatedURL

}

func RandStringBytesRmndr(n int) string {

	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
