package main

import (
	"fmt"
	"net/http"
)

func (app *application) getUrl(w http.ResponseWriter, r *http.Request) {

	findId := "SELECT id FROM url WHERE id=?;"
	retrieveURL := "SELECT url FROM url WHERE id=?;"
	insertVal := "INSERT INTO url VALUES (?,?);"

	checkId := 0

	retURL := ""

	token, err := r.Cookie("usr")

	if err != nil {
		app.errLog.Println(err)
		clientError(w, http.StatusUnauthorized)
		return
	}

	username, err := authenticate(token.Value, "BUSINESS")

	if err != nil {
		app.errLog.Println(err)
		clientError(w, http.StatusUnauthorized)
		return
	}

	id, ok := app.idFromUsername(username)

	if !ok {
		clientError(w, http.StatusInternalServerError)
		return
	}

	result := app.db.QueryRow(findId, id)

	if err = result.Scan(&checkId); err != nil {
		app.errLog.Println(err)
	}

	if id == checkId {
		row := app.db.QueryRow(retrieveURL, id)
		err = row.Scan(&retURL)

		if err != nil {
			app.errLog.Println(err)
			clientError(w, http.StatusInternalServerError)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"url\":\"%s\"}", retURL)))
	} else {
		//generate unique url and insert into table
		retURL = app.generateURL()
		_, err = app.db.Exec(insertVal, id, retURL)

		if err != nil {
			app.errLog.Println(err)
			clientError(w, http.StatusInternalServerError)
			return
		}

		w.Write([]byte(fmt.Sprintf("{\"url\":\"%s\"}", retURL)))
	}

}

func (app *application) getDetails(w http.ResponseWriter, r *http.Request) {

}
