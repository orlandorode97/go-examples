package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {

	file, err := os.Open("/path/to/your/file.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	fileSaleCodes := make(map[string]string)
	for _, record := range records {
		fileSaleCodes[record[2]] = ""

	}

	delete(fileSaleCodes, "9.55E+07")

	delete(fileSaleCodes, "displayName")

	var saleCodes []string
	for key := range fileSaleCodes {
		saleCodes = append(saleCodes, key)
	}

	cfg := mysql.NewConfig()
	cfg.User = "user"
	cfg.Net = "tcp"
	cfg.Addr = "localhost:3306"
	cfg.DBName = "database"

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err)
	}

	builder := sq.Select("*").From("custom_field_values").Where(sq.Eq{"sale_code": saleCodes})

	sql, args, err := builder.ToSql()

	fmt.Println(sql, args)

}
