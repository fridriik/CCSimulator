package main

import (
	"database/sql"
//	"fmt"
	_"github.com/lib/pq"
	"log"
)
type cliente struct {
	nro_cliente int
	nombre, apellido, domicilio string
	telefono [12] rune
}

func createDatabase(){

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err !=nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err =db.Exec(`create database tp`)
	if err != nil {
		log.Fatal(err)
	}
		}
		
func deleteDatabase(){
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err !=nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err =db.Exec(`drop database if exists tp`)
	if err != nil {
		log.Fatal(err)
	}
}
	
func main(){

    deleteDatabase()
	createDatabase()

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=tp sslmode=disable")
	if err !=nil {
		log.Fatal(err)
	}
	defer db.Close()

	 _, err = db.Exec(`create table cliente(
    nrocliente int,
    nombre text,
    apellido text,
    domicilio text,
    telefono char(12)
    )`)
    
    if err !=nil {
    	log.Fatal(err)
    }
    
}
