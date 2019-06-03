package model

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Name     string `sql:"not null"`
	Password string `sql:"not null"`
	Tasks    []Task `gorm:"foreignkey:UserRefer"`
}

type Task struct {
	gorm.Model
	Description string `sql:"not null"`
	Completed   bool   `sql:"not null"`
	UserRefer   uint   `sql:"not null"`
}

type Tasks []Task

var db *gorm.DB

func init() {
	var err error
	var datasource string
	if os.Getenv("DATABASE_URL") != "" {
		// Heroku
		datasource = "b83548ac2c0113:9c9db49c@tcp(us-cdbr-iron-east-02.cleardb.net:3306)/heroku_84f9a6ef922e736?parseTime=true"
	} else {
		// local
		datasource = "root@tcp(127.0.0.1:3306)/todoapp?charset=utf8&parseTime=True&loc=Local"
	}

	db, err = gorm.Open("mysql", datasource)
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(User{})
	db.AutoMigrate(Task{})
}

func CreateUser(user *User) {
	db.Create(user)
}

func FindUser(u *User) User {
	var user User
	db.Where(u).First(&user)
	return user
}

func CreateTask(task *Task) {
	db.Create(task)
}

func FindTasks(t *Task) Tasks {
	var tasks Tasks
	db.Where(t).Find(&tasks)
	return tasks
}

func UpdateTask(t *Task) error {
	rows := db.Model(t).Update(map[string]interface{}{
		"completed": t.Completed,
	}).RowsAffected
	if rows == 0 {
		return fmt.Errorf("Could not find Task (%v) to update", t)
	}
	return nil
}

func DeleteTask(t *Task) error {
	if rows := db.Where(t).Delete(&Task{}).RowsAffected; rows == 0 {
		return fmt.Errorf("Could not find Task (%v) to delete", t)
	}
	return nil
}
