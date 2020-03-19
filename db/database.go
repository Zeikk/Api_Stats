package db

import(
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)
func OpenDB() *sql.DB {

	db, err := sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/api_go")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Connexion à la base")

	return db
}