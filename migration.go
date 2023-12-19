package main

import (
	"strings"

	"github.com/datnm/student_management/pkg/readfile"
	"github.com/jmoiron/sqlx"
)

func Migrate(*sqlx.DB) error {
	db, err := sqlxConnect()
	if err != nil {
		return err
	}

	tx1 := db.MustBegin()
	if _, err = tx1.Exec("SET FOREIGN_KEY_CHECKS=0;"); err != nil {
		return err
	}

	if _, err = tx1.Exec("DROP TABLE IF EXISTS provinces; "); err != nil {
		return err
	}

	if _, err = tx1.Exec("DROP TABLE IF EXISTS countries; "); err != nil {
		return err
	}

	// if _, err = tx1.Exec("DROP TABLE IF EXISTS students "); err != nil {
	// 	return err
	// }

	// if _, err = tx1.Exec("DROP TABLE IF EXISTS students_with_country "); err != nil {
	// 	return err
	// }

	if _, err = tx1.Exec("SET FOREIGN_KEY_CHECKS=1;"); err != nil {
		return err
	}

	if _, err = db.Exec(createTableCountries); err != nil {
		return err
	}

	if _, err = db.Exec(createTableProvinces); err != nil {
		return err
	}

	if _, err = db.Exec(createTableStudents); err != nil {
		return err
	}

	if _, err = db.Exec(createTableStudentsWithCountry); err != nil {
		return err
	}

	fileCountries, _ := readfile.ReadFile("data/countries.csv")
	rowCountries := strings.Split(fileCountries, "\n")

	for i, textDataCountries := range rowCountries {
		if i == 0 {
			continue
		}
		result := strings.Split(textDataCountries, ",")

		if result[3] == "" {
			result[3] = "0"
		}

		if result[4] == "" {
			result[4] = "0"
		}

		if _, err = tx1.Exec(
			`INSERT INTO countries (id, code, name, latitude, longitude) VALUES
			(?, ?, ?, ?, ?);`,
			result[0], result[1], result[2], result[3], result[4],
		); err != nil {
			return err
		}
	}

	fileProvinces, _ := readfile.ReadFile("data/provinces.csv")
	rowProvinces := strings.Split(fileProvinces, "\n")
	country, err := GetCountryByName(db, "VietNam")

	if err != nil {
		return err
	}

	for _, textDataProvinces := range rowProvinces {
		result := strings.Split(textDataProvinces, ",")
		if _, err = tx1.Exec(`INSERT INTO provinces(id, name, country_id) VALUES (?, ?, ?);`,
			result[0], result[1], country.ID,
		); err != nil {
			return err
		}
	}

	// numberOfStudents := 1000000
	// students, err := SeedStudents(numberOfStudents)

	// if err != nil {
	// 	panic(err)
	// }

	// chunkstudentsnumber := chunk.Chunk[*Student](students, 1000)
	// for _, students := range chunkstudentsnumber {
	// 	arguments := make([]any, 0)
	// 	var queries []string
	// 	for _, student := range students {
	// 		query := "(?, ?, ?, ?, ?, ?)"
	// 		birthYear := rand.Intn(2010-1980) + 1980
	// 		queries = append(queries, query)
	// 		arguments = append(arguments, student.FirstName, student.LastName, student.Sex, student.PhoneNumber, student.Falcuty, birthYear)
	// 	}
	// 	data := strings.Join(queries, ",\n")
	// 	dataqueries := fmt.Sprintf(
	// 		`INSERT INTO students(
	// 		first_name,
	// 		last_name,
	// 		male,
	// 		phone_number,
	// 		falcuty,
	// 		birth_year)
	// 		VALUES
	// 		%s;`, data)
	// 	if _, err = tx1.Exec(dataqueries,
	// 		arguments...); err != nil {
	// 		return err
	// 	}
	// }

	// var listOfStudents []*StudentWithCountry
	// for i := 1; i <= numberOfStudents; i++ {
	// 	recordStudentCountries := rand.Intn(5-1) + 1
	// 	for j := 1; j <= recordStudentCountries; j++ {
	// 		studentCountryID := rand.Intn(245-1) + 1
	// 		studentWithNationalities := &StudentWithCountry{
	// 			StudentId:       i,
	// 			NationalitiesId: studentCountryID,
	// 		}

	// 		listOfStudents = append(listOfStudents, studentWithNationalities)
	// 	}
	// }

	// chunkstudentWithCountry := chunk.Chunk[*StudentWithCountry](listOfStudents, 1000)
	// for _, studentData := range chunkstudentWithCountry {
	// 	arguments := make([]any, 0)
	// 	var queries []string
	// 	for _, student := range studentData {
	// 		query := "(?, ?)"
	// 		queries = append(queries, query)
	// 		arguments = append(arguments, student.StudentId, student.NationalitiesId)
	// 	}

	// 	data := strings.Join(queries, ",\n")
	// 	dataqueries := fmt.Sprintf(`INSERT INTO students_with_country(
	// 		student_id,
	// 		student_nationalities_id)
	// 		VALUES
	// 		%s;`, data)
	// 	if _, err = tx1.Exec(dataqueries,
	// 		arguments...); err != nil {
	// 		return err
	// 	}
	// }
	tx1.Commit()
	return nil
}
