package grades

import (
	"fmt"
	"sync"
)

var (
	students     Students
	studentMutex sync.Mutex
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}

type GradeType string

const (
	GradeTest     = GradeType("Test")
	GradeHomework = GradeType("HomeWork")
	GradeQuiz     = GradeType("Quiz")
)

func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}

	return result / float32(len(s.Grades))
}

type Students []Student

func (s Students) GetByID(id int) (*Student, error) {
	for i := range s {
		if id == s[i].ID {
			return &s[i], nil
		}
	}

	return nil, fmt.Errorf("student with id %v not found ", id)
}
