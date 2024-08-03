# Student Management API

This project is a simple Student Management API written in Golang using the [Fiber](https://github.com/gofiber/fiber) web framework. The API allows for managing students and retrieving country information.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Contributing](#contributing)
- [License](#license)

## Installation

To get started with this project, clone the repository and install the necessary dependencies:

```bash
git clone https://github.com/datnm23/students_management.git
cd students_management
go get
```

## Usage

To run the API server, use the following command:

`go run main.go`

By default, the server will start on `http://localhost:3000`.

## API Endpoints

### Create a Student

- **URL:** `/api/students`
- **Method:** `POST`
- **Description:** Adds a new student.
- **Request Body:**

  `{
  "name": "John Doe",
  "age": 21,
  "country": "USA"
}`

### Get All Countries

- **URL:** `/api/countries`
- **Method:** `GET`
- **Description:** Retrieves a list of all countries.

### Get All Students

- **URL:** `/api/students`
- **Method:** `GET`
- **Description:** Retrieves a list of all students.

### Update a Student

- **URL:** `/api/students/:id`
- **Method:** `PUT`
- **Description:** Updates the information of an existing student.
- **Request Body:**

  `{
  "name": "Jane Doe",
  "age": 22,
  "country": "Canada"
}`

### Delete a Student

- **URL:** `/api/students/:id`
- **Method:** `DELETE`
- **Description:** Deletes a student by ID.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any changes.

1.  Fork the repository.
2.  Create a new branch (`git checkout -b feature-branch`).
3.  Commit your changes (`git commit -m 'Add some feature'`).
4.  Push to the branch (`git push origin feature-branch`).
5.  Open a pull request.
