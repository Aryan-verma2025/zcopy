package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) getPrice(w http.ResponseWriter, r *http.Request) {

	stmt := "SELECT one_bw,both_bw,one_cl,both_cl FROM prices WHERE id=?"

	var one_bw int
	var both_bw int
	var one_cl int
	var both_cl int

	token, err := r.Cookie("usr")

	if err != nil {
		clientError(w, http.StatusUnauthorized)
		return
	}

	username, err := authenticate(token.Value, "BUSINESS")

	if err != nil {
		clientError(w, http.StatusUnauthorized)
		return
	}

	var id int64
	row := app.db.QueryRow("SELECT id FROM users WHERE username=?", username)

	err = row.Scan(&id)

	if err != nil {
		app.errLog.Println(err)
		clientError(w, http.StatusBadRequest)
		return
	}

	row = app.db.QueryRow(stmt, int(id))

	err = row.Scan(&one_bw, &both_bw, &one_cl, &both_cl)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Prices not set"))
		return

	}

	w.Write([]byte(fmt.Sprintf("{\"one_bw\":\"%d\",\"both_bw\":\"%d\",\"one_cl\":\"%d\",\"both_cl\":\"%d\"}", one_bw, both_bw, one_cl, both_cl)))

}

func (app *application) setPrice(w http.ResponseWriter, r *http.Request) {

	var intValues [4]int

	names := [4]string{"one_bw", "one_cl", "both_bw", "both_cl"}
	var values [4]string

	retId := "SELECT id FROM users WHERE username=?;"
	var id int
	var checkId int

	findId := "SELECT id FROM prices WHERE id=?;"
	insertStmt := "INSERT INTO prices VALUES (?,?,?,?,?);"
	updateStmt := "UPDATE prices SET one_bw = ? ,one_cl=?,both_bw=?,both_cl=? where id=?;"

	r.ParseForm()

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

	row := app.db.QueryRow(retId, username)

	err = row.Scan(&id)

	if err != nil {
		app.errLog.Println(err)
		clientError(w, http.StatusBadRequest)
		return
	}

	for i := range 4 {
		values[i] = r.PostForm.Get(names[i])
		intValues[i], err = strconv.Atoi(values[i])

		if err != nil {
			app.errLog.Println(err)
			clientError(w, http.StatusBadRequest)
			return
		}
		if !(intValues[i] > 0 && intValues[i] < 15001) {
			clientError(w, http.StatusBadRequest)
			return
		}
	}

	result := app.db.QueryRow(findId, id)

	if err = result.Scan(&checkId); err != nil {
		app.errLog.Println(err)
		//clientError(w, http.StatusInternalServerError)
		//return
	}

	if checkId == id {

		_, err := app.db.Exec(updateStmt, values[0], values[1], values[2], values[3], id)

		if err != nil {
			app.errLog.Println(err)
			clientError(w, http.StatusInternalServerError)
		}
	} else {
		_, err = app.db.Exec(insertStmt, id, values[0], values[1], values[2], values[3])

		if err != nil {
			app.errLog.Println(err)
			clientError(w, http.StatusInternalServerError)
		}
	}

}
