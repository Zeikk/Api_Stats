package control

import(
	"fmt"
	"net/http"
	jwt "github.com/dgrijalva/jwt-go"
	db "api_stats/db"
	"log"
)

func getId(w http.ResponseWriter, r *http.Request) string{

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			fmt.Fprintf(w, "\nNécessite Authentification")
			w.WriteHeader(http.StatusUnauthorized)
			return ""
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "\nNécessite Authentification")
		return ""
	}

	tokenStr := c.Value

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("grain_de_sel"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			fmt.Fprintf(w, "\nNécessite Authentification")
			w.WriteHeader(http.StatusUnauthorized)
			return ""
		}
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "\nNécessite Authentification")
		return ""
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "\nNécessite Authentification")
		return ""
	}

	log.Println(claims.IdMedecin)
	return claims.IdMedecin
}	

func GetStatsMaladie(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Pourcentage de Patients par Maladie")
	
	idMedecin := getId(w, r)
	db := db.OpenDB()

	rows, err := db.Query("select count(idPatient) / (select count(*) from etre_malade) * 100 as nbPatient, libelleMaladie from etre_malade join maladie using(idMaladie) join suivre using(idPatient) where idMedecin = ?	group by libelleMaladie", idMedecin)

	defer rows.Close()

	var nb, maladie string
	for rows.Next() {
		err := rows.Scan(&nb, &maladie)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, "\nMaladie: %s NbPatients: %s", maladie, nb)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	db.Close()

}

func GetStatsAge(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprintf(w, "Répartition des Maladies par Tranche d'âge")
	idMedecin := getId(w, r)
	if idMedecin == ""{
		return 
	}
	
	db := db.OpenDB()

	for i := 0; i<99; i+=10 {
		rows, err := db.Query("select count(libelleMaladie) as nbPatient, libelleMaladie from maladie join etre_malade using(idMaladie) join patient using(idPatient) join suivre using(idPatient) join personne using(idPersonne) where idMedecin = ? and ROUND(DATEDIFF(SYSDATE(), dateDeNaissance))/365 BETWEEN ? and ?	group by libelleMaladie", idMedecin, i, i+9);
		defer rows.Close()

		fmt.Fprintf(w, "\nTranche d'âge: %d-%d", i, i+9)
		var nbPatient, maladie string
		for rows.Next() {
			err := rows.Scan(&nbPatient, &maladie)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Fprintf(w, "\nMaladie: %s Pourcentage: %s", maladie, nbPatient)
		}
		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}

	db.Close()
}