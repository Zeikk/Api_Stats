package control

import(
	"fmt"
	"log"
	"net/http"
	db "api_stats/db"
	"encoding/json"
	"time"
	"github.com/dgrijalva/jwt-go"
)

type User struct {
	IdPersonne string `json:"id"`
	Password string `json:"password"`
}

type Claims struct {
	IdMedecin string `json:"id"`
	jwt.StandardClaims
}

func LoginMedecin(w http.ResponseWriter, r *http.Request) {

	//fmt.Fprintln(w, "Login")
	var medecin User
	err := json.NewDecoder(r.Body).Decode(&medecin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	
	database := db.OpenDB()
	rows, err := database.Query("select passwordMedecin, idMedecin from medecin where idPersonne = ?", medecin.IdPersonne)

	var password, id string
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&password, &id)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	

	if password != medecin.Password {
		fmt.Fprintf(w, "\nIdentifiant / Mot de passe invalide")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	claims := &Claims{
		IdMedecin: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	tokenString, err := token.SignedString([]byte("grain_de_sel"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{Name: "token", Value: tokenString, HttpOnly: true, MaxAge: 50000, Path: "/"}
	
	
	http.SetCookie(w, &cookie)
	database.Close()
}

func LogoutMedecin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}