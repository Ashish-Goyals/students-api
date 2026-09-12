package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/ashish-goyals/students-api/internal/config"
	"github.com/ashish-goyals/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB

}



func New(cfg *config.Config)(*Sqlite,error){
   db, err := sql.Open("sqlite3",cfg.StoragePath)
	if err != nil{
		return nil,err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
    email TEXT,
	age INTEGER
	)`)

	if err !=nil {
		return nil, err

	}

return &Sqlite{
	Db: db,

},nil
}

func (s *Sqlite) CreateStudent (name string,email string,age int)(int64,error){
   stmt, err :=	s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")

   if err != nil{
	return 0,err
   }
   defer stmt.Close()

  result,err := stmt.Exec(name,email,age)

  if err != nil{
	return 0,err
  }

  lastId,err := result.LastInsertId()

  if err != nil{
	return 0,err
  }

  return lastId, nil

}

func (s *Sqlite) GetStudentbyId(id int64) (types.Student, error){
   stmt , err := s.Db.Prepare("SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1")

   if err != nil{
	return types.Student{},err
   }

   defer stmt.Close()

   var student types.Student

   err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)



   if err != nil{
	if err == sql.ErrNoRows {
		return types.Student{},fmt.Errorf("No Student fond with id %s", fmt.Sprint(id))
	}
	return types.Student{},fmt.Errorf("Query Error : %w",err)
   }

   return student, nil

}