package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	infoLog *log.Logger
	errLog  *log.Logger
	db      *sql.DB
}

func main() {

	server_database_dsn := "aryan_zcopy:xyab*!2j6imr@tcp(mysql-aryan.alwaysdata.net:3306)/aryan_zcopy"
	port := ":8100"

	//local_database_dsn := "zcpy:pass@/zcopy?parseTime=true"
	//port := ":8080"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := sql.Open("mysql", server_database_dsn)

	if err != nil {
		errLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		infoLog: infoLog,
		errLog:  errLog,
		db:      db,
	}

	srv := &http.Server{
		Addr:     port,
		ErrorLog: errLog,
		Handler:  app.routes(),
	}

	infoLog.Println("Server started on port 8080")
	log.Fatal(srv.ListenAndServe())

}

//****************************************************************************************************************************************************************
// connStr := "postgresql://mypsql_0xxs_user:xvEcApjQxkPm1Ny8gSc0rZHWl5vRWLZT@dpg-cqhonqo8fa8c73c08kng-a.singapore-postgres.render.com/mypsql_0xxs"
// db, err := sql.Open("postgres", connStr)
// defer db.Close()

// if err != nil {
// 	log.Fatal(err)
// }

// row := db.QueryRow("SELECT * FROM users WHERE id=$1", 1)

// u := &user{}

// err = row.Scan(&u.id, &u.name)

// if err == sql.ErrNoRows {
// 	log.Fatal("no records found")
// } else if err != nil {
// 	log.Fatal(err)
// }

// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Printf("%v", u)
// type businessUser struct {
// 	username string `json:"Username"`
// 	password string `json:"Password"`
// }
