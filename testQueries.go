package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strconv"

)

type City struct {
	ID          int    `json:"id,omitempty"  db:"ID"`
	Name        string `json:"name,omitempty"  db:"Name"`
	CountryCode string `json:"countryCode,omitempty"  db:"CountryCode"`
	District    string `json:"district,omitempty"  db:"District"`
	Population  int    `json:"population,omitempty"  db:"Population"`
}

type Country struct {
	Code           string  `json:"code,omitempty"  db:"Code"`
	Name           string  `json:"name,omitempty"  db:"Name"`
	Continent      string  `json:"continent,omitempty"  db:"Continent"`
	Region         string  `json:"region,omitempty"  db:"Region"`
	SurfaceArea    float64 `json:"surfaceArea,omitempty"  db:"SurfaceArea"`
	IndepYear      int     `json:"indepyear,omitempty"  db:"IndepYear"`
	Population     int     `json:"population,omitempty"  db:"Population"`
	LifeExpectancy float64 `json:"lifeExpectancy,omitempty"  db:"LifeExpectancy"`
	GNP            float64 `json:"gnp,omitempty"  db:"GNP"`
	GNPOld         float64 `json:"gnpOld,omitempty"  db:"GNPOld"`
	LocalName      string  `json:"localName,omitempty"  db:"LocalName"`
	GovernmentForm string  `json:"governmentForm,omitempty"  db:"GovernmentForm"`
	HeadOfState    string  `json:"headofState,omitempty"  db:"HeadOfState"`
	Capital        int     `json:"capital,omitempty"  db:"Capital"`
	Code2          string  `json:"code2,omitempty"  db:"Code2"`
}

func main() {
	db, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USERNAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOSTNAME"), os.Getenv("DB_PORT"), os.Getenv("DB_DATABASE")))
	if err != nil {
		log.Fatalf("Cannot Connect to Database: %s", err)
	}

	fmt.Println("Connected!")
	// country := []Country{}
	// city := City{}
	cities := []City{}

	//色々なクエリ
	// db.Get(&city, "SELECT * FROM city WHERE Name='"+os.Args[1]+"'")
	// db.Get(&country, "select * from country where Name = 'Japan'")
	// db.Get(&country, "select * from country where Code ='" + city.CountryCode + "'")
	// fmt.Printf("%sの人口はその国の%.2g%%です\n", os.Args[1], float64(city.Population)/float64(country.Population)*100)

	//複数のデータの取得
	// db.Select(&cities,"select * from city where CountryCode = 'JPN'")
	// fmt.Println("日本の年一覧")
	// for _, city := range cities {
	// 	fmt.Printf("都市名: %s, 人口: %d人\n", city.Name,city.Population)
	// }

	//都市の追加
	// people,err := strconv.Atoi(os.Args[4])
	// db.Exec(`insert into city (Name, CountryCode, District,Population) values (?,?,?,?)`,os.Args[1],os.Args[2],os.Args[3],people)
	// db.Select(&cities,"select * from city where CountryCode = 'JPN'")
	// for _,city := range cities {
	// 	fmt.Printf("都市名: %s, 人口: %d人\n", city.Name,city.Population)
	// }
}
