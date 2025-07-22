package storage

import "github.com/krshnas/students-api/internal/types"

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetStudents() ([]types.Student, error)
	UpdateStudent(email string, age int, id int64) (int64, error)
	DeleteStudent(id int64) (bool, error)
}
