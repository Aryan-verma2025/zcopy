package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func (app *application) registerUser(w http.ResponseWriter, r *http.Request) {

	stmt := "INSERT into users(password,username,type) values(?,?,?);"

	r.ParseForm()

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	_type := r.PostForm.Get("type")

	//fmt.Println("Got this valuse:", username, password, _type)
	if len(username) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username should be at least 8 charaacters long"))

		return
	} else if len(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password should be at least 8 charaacters long"))

		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)

	if err != nil {
		clientError(w, http.StatusInternalServerError)
		return
	}

	_, err = app.db.Exec(stmt, string(hash), username, _type)

	if err != nil {
		app.errLog.Println(err)
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.name_UNIQUE") {

				//Duplicate email
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("Username is already registered"))

				return
			}
		}
		clientError(w, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {

	utype := ""
	uname := ""
	var hashedPassword []byte

	stmt := "SELECT username,password,type FROM users WHERE username=?;"

	r.ParseForm()

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	//fmt.Println("Got this valuse:", username, password, _type)
	if len(username) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("username should be at least 8 charaacters long"))

		return
	} else if len(password) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password should be at least 8 charaacters long"))

		return
	}

	row := app.db.QueryRow(stmt, username)

	err := row.Scan(&uname, &hashedPassword, &utype)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Username or Password"))
		return
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	if err != nil {
		app.errLog.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Username or Password"))
		return
	}
	token, err := createToken(username, utype)

	if err != nil {
		app.errLog.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid Username or Password"))
		return
	}

	w.Header().Add("set-cookie", fmt.Sprintf("usr=%s;HttpOnly;", token))

}
