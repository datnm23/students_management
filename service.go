package main

import (
	"math"
)

type Service struct {
	PageNumber  int
	TotalRecord int
}

func Pagination(i int) (*Service, error) {
	db, err := sqlxConnect()

	if err != nil {
		return nil, err
	}
	var totalrecord int
	// service := &Service{}
	err = db.QueryRow(`SELECT COUNT(*) FROM students`).Scan(&totalrecord)
	if err != nil {
		return nil, err
	}

	pagenumber := int(math.Round(float64(totalrecord / i)))
	pagination := &Service{
		PageNumber:  pagenumber,
		TotalRecord: totalrecord,
	}
	return pagination, nil
}
