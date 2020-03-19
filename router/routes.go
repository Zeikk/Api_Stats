package router

import(
	"net/http"
	control "api_stats/control"
	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Methode      string
	Lien     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{

	Route{
		Name: "getStatsMaladie",
		Methode: "GET",
		Lien: "/stats/maladie",
		HandlerFunc: control.GetStatsMaladie,
	},
	Route{
		Name: "getStatsAge",
		Methode: "GET",
		Lien: "/stats/age",
		HandlerFunc: control.GetStatsAge,
	},
	Route{
		Name: "loginMedecin",
		Methode: "GET",
		Lien: "/user/login",
		HandlerFunc: control.LoginMedecin,
	},
	Route{
		Name: "logoutMedecin",
		Methode: "GET",
		Lien: "/user/logout",
		HandlerFunc: control.LogoutMedecin,
	},
}


func Router() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {

		router.
			Methods(route.Methode).
			Path(route.Lien).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}