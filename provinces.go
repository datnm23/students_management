package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Province struct {
	ID        int      `db:"id"`
	Name      string   `db:"name"`
	CountryID int      `db:"country_id"`
	Country   *Country `db:"country"`
}

type ProvinceWithCountry struct {
	Id          int    `db:"provinces_id"`
	Name        string `db:"provinces_name"`
	CountryId   int    `db:"countries_id"`
	CountryName string `db:"countries_name"`
}

var createTableProvinces = `
CREATE TABLE IF NOT EXISTS provinces (
	id INT PRIMARY KEY AUTO_INCREMENT,
	name VARCHAR(100),
	country_id INT,
	FOREIGN KEY (country_id) REFERENCES countries(id) 
);`

func GetAllProvinces(db *sqlx.DB) ([]*Province, error) {
	provinces := []*Province{}
	err := db.Select(&provinces, "SELECT * FROM provinces")
	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func GetOneProvince(db *sqlx.DB, id int) (*Province, error) {
	province := &Province{}
	err := db.Get(province, "SELECT * FROM provinces WHERE (id =?)", id)

	if err != nil {
		return nil, err
	}

	return province, nil
}

func GetProvincesByCountryID(db *sqlx.DB, countryID int) ([]*Province, error) {
	provinces := []*Province{}
	err := db.Select(&provinces, "SELECT * FROM provinces LIKE (country_id= ? )", countryID)

	if err != nil {
		return nil, err
	}

	return provinces, nil
}

func SearchProvinceByName(db *sqlx.DB, name string) ([]*Province, error) {
	provinces := []*Province{}
	nameSeach := fmt.Sprintf("%%%v%%", name)
	countryProvinces := []*ProvinceWithCountry{}
	err := db.Select(
		&countryProvinces,
		`SELECT
			provinces.id as provinces_id,
			provinces.name as provinces_name,
			countries.id as countries_id,
			countries.name as countries_name
		FROM provinces
		LEFT JOIN countries
		ON provinces.country_id=countries.id
		WHERE provinces.name LIKE ?;`, nameSeach,
	)

	if err != nil {
		return nil, err
	}

	countryIDs := []interface{}{}

	for _, countryID := range countryProvinces {
		countryIDs = append(countryIDs, countryID.CountryId)
	}

	countries := []*Country{}
	for _, countryID := range countryIDs {
		err = db.Select(&countries, "SELECT * FROM countries WHERE id IN (?);", countryID)
		if err != nil {
			return nil, err
		}
	}

	for _, countryProvince := range countryProvinces {
		for _, country := range countries {
			if countryProvince.CountryId == country.ID {
				province := &Province{
					ID:        countryProvince.Id,
					Name:      countryProvince.Name,
					CountryID: countryProvince.CountryId,
					Country:   country,
				}

				provinces = append(provinces, province)
			}
		}
	}

	return provinces, nil
}
