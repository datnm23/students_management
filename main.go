package main

import (
	"fmt"
	"strconv"
	"strings"
	"github.com/davecgh/go-spew/spew"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/middleware/etag"
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
	//Recover
	app.Use(recover.New())
	app.Get("/", func(c *fiber.Ctx) error {
   	 panic("I'm an error")
	})

	//Logger
	app.Use(logger.New())
	app.Use(logger.New(logger.Config{
	    Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
	    Format: "${pid} ${locals:requestid} ${status} - ${method} ${path}​\n",
	}))

	app.Use(logger.New(logger.Config{
	    Format:     "${pid} ${status} - ${method} ${path}\n",
	    TimeFormat: "02-Jan-2006",
	    TimeZone:   "America/New_York",
	}))
	
	file, err := os.OpenFile("./123.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	app.Use(logger.New(logger.Config{
	    Output: file,
	}))
	
	// Add Custom Tags
	app.Use(logger.New(logger.Config{
	    CustomTags: map[string]logger.LogFunc{
	        "custom_tag": func(output logger.Buffer, c *fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
	            return output.WriteString("it is a custom tag")
	        },
	    },
	}))
	
	// Callback after log is written
	app.Use(logger.New(logger.Config{
	    TimeFormat: time.RFC3339Nano,
	    TimeZone:   "Asia/Shanghai",
	    Done: func(c *fiber.Ctx, logString []byte) {
	        if c.Response().StatusCode() != fiber.StatusOK {
	            reporter.SendToSlack(logString) 
	        }
	    },
	}))
	
	// Disable colors when outputting to default format
	app.Use(logger.New(logger.Config{
	    DisableColors: true,
	}))
	//Compress
	app.Use(compress.New())
	
	app.Use(compress.New(compress.Config{
	    Level: compress.LevelBestSpeed, // 1
	}))

	app.Use(compress.New(compress.Config{
	  Next:  func(c *fiber.Ctx) bool {
	    return c.Path() == "/dont_compress"
	  },
	  Level: compress.LevelBestSpeed, // 1
	}))

	//Cache
	app.Use(cache.New())
	
	app.Use(cache.New(cache.Config{
	    Next: func(c *fiber.Ctx) bool {
	        return c.Query("noCache") == "true"
	    },
	    Expiration: 30 * time.Minute,
	    CacheControl: true,
	}))

	//RequestID
	app.Use(requestid.New())
	
	app.Use(requestid.New(requestid.Config{
	    Header:    "X-Custom-Header",
	    Generator: func() string {
	        return "static-id"
	    },
	}))
	
	//etag
	app.Use(etag.New())
	
	app.Get("/", func(c *fiber.Ctx) error {
	    return c.SendString("Hello, World!")
	})
	
	app.Use(etag.New(etag.Config{
	    Weak: true,
	}))
	
	app.Get("/", func(c *fiber.Ctx) error {
	    return c.SendString("Hello, World!")
	})
	
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
