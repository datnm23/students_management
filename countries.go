package main

import (
	"github.com/jmoiron/sqlx"
)

type Country struct {
	ID        int
	Code      string
	Name      string
	Latitude  float64
	Longitude float64
	Provinces []*Province
}

var createTableCountries = `
CREATE TABLE IF NOT EXISTS countries (
	id INT PRIMARY KEY AUTO_INCREMENT,
	code VARCHAR(100),
	name VARCHAR(100),
	latitude FLOAT(24) NULL,
	longitude FLOAT(24) NULL,
	UNIQUE (code),
	UNIQUE (name)
);`

func GetAllCountries(db *sqlx.DB) ([]*Country, error) {
	countries := []*Country{}
	err := db.Select(&countries, "SELECT * FROM countries")
	if err != nil {
		return nil, err
	}

	return countries, nil

}

func GetOneCountry(db *sqlx.DB, id int) (*Country, error) {
	country := &Country{}
	err := db.Get(country, "SELECT * FROM countries WHERE (id=?);", id)

	if err != nil {
		return nil, err
	}
	return country, nil
}

func GetCountryByName(db *sqlx.DB, name string) (*Country, error) {
	country := &Country{}
	err := db.Get(country, "SELECT id FROM countries WHERE (name=?) LIMIT 1;", name)
	if err != nil {
		return nil, err
	}

	return country, nil
}

func GetCountryWithProvincesByID(db *sqlx.DB, id int) (*Country, error) {

	country, err := GetOneCountry(db, id)
	if err != nil {
		return nil, err
	}

	country.Provinces, err = GetProvincesByCountryID(db, id)
	if err != nil {
		return nil, err
	}
	return country, nil

}

func GetAllCountriesWithProvincesV1(db *sqlx.DB) ([]*Country, error) {
	countries, err := GetAllCountries(db)

	if err != nil {
		return nil, err
	}

	for i, country := range countries {
		id := countries[i].ID
		country.Provinces, err = GetProvincesByCountryID(db, id)

		if err != nil {
			return nil, err
		}
	}

	return countries, nil
}

func GetAllCountriesWithProvincesV2(db *sqlx.DB) ([]*Country, error) {
	countries, err := GetAllCountries(db)

	if err != nil {
		return nil, err
	}
	provinces, err := GetAllProvinces(db)

	if err != nil {
		return nil, err
	}

	for _, country := range countries {
		for _, province := range provinces {
			if country.ID == province.CountryID {
				country.Provinces = append(country.Provinces, province)
			}
		}
	}

	return countries, nil
}
func GetAllCountriesWithProvincesV3(db *sqlx.DB) ([]*Country, error) {
	country := []*Country{}
	err := db.Select(&country, `SELECT provinces.id ,provinces.name as provinces_name, countries.name as countries_name FROM countries LEFT JOIN provinces ON provinces.country_id=countries.id;`)

	if err != nil {
		return nil, err
	}
	return country, nil
}
