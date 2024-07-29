package main

import "net/http"

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api", workingOnIt)          //apiDoc
	mux.HandleFunc("POST /login", app.loginUser)     //----------------------------------DONE
	mux.HandleFunc("POST /signup", app.registerUser) //----------------------------------DONE

	mux.HandleFunc("POST /upload/{id}", workingOnIt) //upload)
	mux.HandleFunc("GET /business/orders/{id}", workingOnIt)
	mux.HandleFunc("GET /prints/{id}", workingOnIt)
	mux.HandleFunc("POST /business/price", app.setPrice) //---------------------------------DONE
	mux.HandleFunc("GET /business/price", app.getPrice)  //---------------------------------DONE
	mux.HandleFunc("GET /business/details", workingOnIt)
	mux.HandleFunc("GET /get-url", app.getUrl) //------------------------------DONE
	mux.HandleFunc("GET /logout", logout)      //-----------------------------------DONE

	return mux
}
