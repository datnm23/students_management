package main

import (
	"github.com/go-faker/faker/v4"
	"github.com/jmoiron/sqlx"
)

type Student struct {
	Id            int        `db:"id" json:"id"`
	FirstName     string     `db:"first_name" faker:"first_name" json:"first_name"`
	LastName      string     `db:"last_name" faker:"last_name" json:"last_name"`
	Sex           string     `db:"male" faker:"oneof:Male,Female" json:"sex"`
	PhoneNumber   string     `db:"phone_number" faker:"phone_number" json:"phone_number"`
	Falcuty       string     `db:"falcuty" faker:"oneof:Medicine,Law,Business,Engineering,Education,Divinity,Government,Design" json:"falcuty"`
	BirthYear     int        `db:"birth_year" faker:"-" json:"birth_year"`
	Nationalities []*Country `faker:"-" json:"nationalities"`
}

type StudentWithCountry struct {
	StudentId          int     `db:"student_id"`
	StudentFirstName   string  `db:"student_first_name"`
	StudentLastName    string  `db:"student_last_name"`
	StudentSex         string  `db:"student_male"`
	StudentPhoneNumber string  `db:"student_phone_number"`
	StudentFalcuty     string  `db:"student_falcuty"`
	StudentBirthyear   int     `db:"student_birthyear"`
	Code               string  `db:"countries_code"`
	Latitude           float64 `db:"countries_latitude"`
	Longitude          float64 `db:"countries_longitude"`
	NationalitiesId    int     `db:"student_nationalities_id"`
	Nationalities      string  `db:"student_nationalities"`
}

type StudentWithCountryV1 struct {
	StudentId          int    `db:"student_id"`
	StudentFirstName   string `db:"student_first_name"`
	StudentLastName    string `db:"student_last_name"`
	StudentSex         string `db:"student_male"`
	StudentPhoneNumber string `db:"student_phone_number"`
	StudentFalcuty     string `db:"student_falcuty"`
	StudentBirthyear   int    `db:"student_birthyear"`
	NationalitiesId    int    `db:"student_nationalities_id"`
}

var createTableStudents = `
	CREATE TABLE IF NOT EXISTS students (
		id INT PRIMARY KEY AUTO_INCREMENT,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		male VARCHAR(20),
		phone_number VARCHAR(20),
		falcuty VARCHAR(50),
		birth_year INT
	);
`
var createTableStudentsWithCountry = `
	CREATE TABLE IF NOT EXISTS students_with_country (
		student_id INT,
		student_nationalities_id INT,
		FOREIGN KEY (student_id) REFERENCES students(id),
		FOREIGN KEY (student_nationalities_id) REFERENCES countries(id)
	);
`

func SeedStudents(n int) ([]*Student, error) {
	students := []*Student{}

	for i := 1; i <= n; i++ {
		student := &Student{}
		err := faker.FakeData(&student)

		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func GetStudentWithNationalities(db *sqlx.DB, id int) (*Student, error) {
	studentData := []*StudentWithCountry{}
	student := &Student{}
	err := db.Select(&studentData,
		`SELECT students.id as student_id ,
		students.first_name as student_first_name,
		students.last_name as student_last_name,
		students.male as student_male,
		students.phone_number as student_phone_number,
		students.falcuty as student_falcuty,
		students.birth_year as student_birthyear,
		students_with_country.student_nationalities_id as student_nationalities_id,
		countries.name as student_nationalities,
		countries.code as countries_code,
		countries.latitude as countries_latitude,
		countries.longitude as countries_longitude
		FROM students
		LEFT JOIN students_with_country 
		ON students.id=students_with_country.student_id
		LEFT JOIN countries
		ON students_with_country.student_nationalities_id=countries.id
		WHERE (students.id =?);`,
		id)

	if err != nil {
		return nil, err
	}

	for _, studentWithCountry := range studentData {
		countries := []*Country{}

		for _, allOfStudentCountries := range studentData {
			if studentWithCountry.StudentId == allOfStudentCountries.StudentId {
				country := &Country{
					ID:        allOfStudentCountries.NationalitiesId,
					Code:      allOfStudentCountries.Code,
					Name:      allOfStudentCountries.Nationalities,
					Latitude:  allOfStudentCountries.Latitude,
					Longitude: allOfStudentCountries.Longitude,
				}
				countries = append(countries, country)
			}
		}

		student = &Student{
			Id:            studentWithCountry.StudentId,
			FirstName:     studentWithCountry.StudentFirstName,
			LastName:      studentWithCountry.StudentLastName,
			Sex:           studentWithCountry.StudentSex,
			PhoneNumber:   studentWithCountry.StudentPhoneNumber,
			Falcuty:       studentWithCountry.StudentFalcuty,
			BirthYear:     studentWithCountry.StudentBirthyear,
			Nationalities: countries,
		}
	}

	return student, nil
}
