package handler

import (
	"net/http"
	"strconv"

	"github.com/daisuke13/todo-app/server/src/model"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
)

func CreateTask(c echo.Context) error {
	task := new(model.Task)
	if err := c.Bind(task); err != nil {
		return err
	}

	if task.Description == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "invalid to or message fields",
		}
	}

	uid := userIDFromToken(c)
	if user := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}); user.ID == 0 {
		return echo.ErrNotFound
	}

	task.UserRefer = uid
	model.CreateTask(task)

	return c.JSON(http.StatusCreated, task)
}

func GetTasks(c echo.Context) error {
	uid := userIDFromToken(c)

	if user := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}); user.ID == 0 {
		return echo.ErrNotFound
	}

	tasks := model.FindTasks(&model.Task{UserRefer: uid})
	return c.JSON(http.StatusOK, tasks)
}

func UpdateTask(c echo.Context) error {

	uid := userIDFromToken(c)

	if user := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}); user.ID == 0 {
		return echo.ErrNotFound
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.ErrNotFound
	}
	tid := uint(taskID)

	tasks := model.FindTasks(&model.Task{Model: gorm.Model{ID: tid}, UserRefer: uid})
	if len(tasks) == 0 {
		return echo.ErrNotFound
	}
	task := tasks[0]
	task.Completed = !tasks[0].Completed

	if err := model.UpdateTask(&task); err != nil {
		return echo.ErrNotFound
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteTask(c echo.Context) error {
	uid := userIDFromToken(c)

	if user := model.FindUser(&model.User{Model: gorm.Model{ID: uid}}); user.ID == 0 {
		return echo.ErrNotFound
	}

	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return echo.ErrNotFound
	}
	tid := uint(taskID)

	if err := model.DeleteTask(&model.Task{Model: gorm.Model{ID: tid}, UserRefer: uid}); err != nil {
		return echo.ErrNotFound
	}

	return c.NoContent(http.StatusNoContent)
}
