package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
)

func main() {
	db, err := sqlxConnect()

	if err != nil {
		panic(err)
	}

	err = Migrate(db)

	if err != nil {
		panic(err)
	}

	app := fiber.New()

	app.Post("/api/students", func(c *fiber.Ctx) error {
		p := new(Student)
		if err := c.BodyParser(p); err != nil {
			return err
		}


		spew.Dump(p)

		if err != nil {
			return err
		}

		tx := db.MustBegin()
		result, err := tx.Exec(
			`INSERT INTO students(
				first_name,
				last_name,
				male,
				phone_number,
				falcuty,
				birth_year)
				VALUES (?, ?, ?, ?, ?, ?);`,
			p.FirstName, p.LastName, p.Sex, p.PhoneNumber, p.Falcuty, p.BirthYear,
		)

		if err != nil {
			return err
		}

		studentID, err := result.LastInsertId()
		if err != nil {
			return err
		}

		for _, id := range p.Nationalities {
			if _, err = tx.Exec(
				`INSERT INTO students_with_country(
					student_id,
					student_nationalities_id)
					VALUES (?, ?);`,
				studentID, id.ID,
			); err != nil {
				return err
			}
		}

		tx.Commit()
		return nil
	})

	app.Get("/api/countries", func(c *fiber.Ctx) error {
		countries := []*Country{}
		err := db.Select(&countries, "SELECT * FROM countries")

		if err != nil {
			return err
		}

		provinces := []*Province{}
		err = db.Select(&provinces, `SELECT * FROM provinces`)

		if err != nil {
			return err
		}

		for _, country := range countries {
			for _, province := range provinces {
				if province.CountryID == country.ID {
					country.Provinces = append(country.Provinces, province)
				}
			}
		}
		return c.JSON(countries)
	})

	app.Get("/api/students", func(c *fiber.Ctx) error {
		students := []*Student{}
		limit := 100
		pagenumber, err := strconv.Atoi(c.Query("page"))

		if err != nil {
			return err
		}

		if err := db.Select(&students, `SELECT * FROM students LIMIT ? OFFSET ?;`, limit, (pagenumber-1)*limit); err != nil {
			return err
		}

		var result string
		var listID []string

		for _, student := range students {
			result = fmt.Sprintf("%d", student.Id)
			listID = append(listID, result)
		}

		studentWithCountries := []*StudentWithCountryV1{}
		query := strings.Join(listID, ",")
		idqueries := fmt.Sprintf(
			`SELECT
			students.id as student_id,
			students.first_name as student_first_name,
			students.last_name as student_last_name,
			students.male as student_male,
			students.phone_number as student_phone_number,
			students.falcuty as student_falcuty,
			students.birth_year as student_birthyear,
			students_with_country.student_nationalities_id student_nationalities_id
			FROM students
			LEFT JOIN students_with_country
			ON students.id= students_with_country.student_id
			WHERE id IN (%s);`, query)

		if err = db.Select(&studentWithCountries,
			idqueries); err != nil {
			return err
		}

		countries := []*Country{}

		if err = db.Select(&countries, `SELECT DISTINCT * FROM countries`); err != nil {
			return err
		}

		for _, student := range students {
			studentNationalities := []*Country{}
			for _, studentData := range studentWithCountries {
				for _, country := range countries {
					if country.ID == studentData.NationalitiesId {
						studentNationalities = append(studentNationalities, country)
					}
				}
				if student.Id == studentData.StudentId {
					student.Nationalities = studentNationalities
				}
			}
			students = append(students, student)
		}

		pagination, err := Pagination(limit)

		if err != nil {
			return err
		}

		return c.JSON(map[string]any{
			"data": students,
			"key":  pagination,
		})
	})

	app.Put("/api/students/:id", func(c *fiber.Ctx) error {
		p := new(Student)

		if err := c.BodyParser(p); err != nil {
			return err
		}

		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return err
		}

		tx := db.MustBegin()

		if _, err := tx.Exec(
			`DELETE FROM students_with_country 
			WHERE student_id = ?;`,
			id); err != nil {
			return err
		}

		if _, err = tx.Exec(
			`UPDATE students
			SET
				first_name = ?,
				last_name = ?,
				male = ?,
				phone_number = ?,
				falcuty =?,
				birth_year=? 
			WHERE id = ?;`,
			p.FirstName, p.LastName, p.Sex, p.PhoneNumber, p.Falcuty, p.BirthYear, id); err != nil {
			return err
		}

		for _, studentcountries := range p.Nationalities {
			if _, err := tx.Exec(
				`INSERT INTO students_with_country(
					student_id,
					student_nationalities_id)
				VALUES (?, ?);`, id, studentcountries.ID); err != nil {
				return err
			}
		}

		tx.Commit()
		return c.JSON(p)
	})

	app.Delete("/api/students/:id", func(c *fiber.Ctx) error {
		id, err := strconv.Atoi(c.Params("id"))

		if err != nil {
			return err
		}

		tx := db.MustBegin()

		if _, err := tx.Exec(
			`DELETE FROM students_with_country 
			WHERE student_id = ?;`,
			id); err != nil {
			return err
		}

		if _, err := tx.Exec(
			`DELETE FROM students 
			WHERE id = ?;`,
			id); err != nil {
			return err
		}

		tx.Commit()
		return nil
	})

	app.Listen(":3000")
}
