package main

import (
	"errors"
	"fmt"
	"hukaapp/todo/internal/data"
	"hukaapp/todo/internal/validator"
	"net/http"
	"time"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content string `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	task := &data.Task{
		ID:        1,
		CreatedAt: time.Now(),
		Content:   input.Content,
	}

	err = app.model.Task.Create(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/task/%d", task.ID))

	err = app.writeJSON(w, http.StatusCreated, envelop{"task": task}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) readTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task, err := app.model.Task.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) readListTasksHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Content     string
		IsCompleted bool
		data.Filter
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Content = app.readString(qs, "content", "")
	input.IsCompleted = app.readBool(qs, "is-completed", false)

	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page-size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafelist = []string{"id", "content", "is-completed", "created_at", "-id", "-content", "-is-completed", "-created_at"}

	if data.ValidateFilter(v, input.Filter); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	tasks, metadata, err := app.model.Task.ReadAll(input.Content, input.IsCompleted, input.Filter)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"tasks": tasks, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	task, err := app.model.Task.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Content     *string `json:"content,omitempty"`
		IsCompleted *bool   `json:"is-completed,omitempty"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Content != nil {
		task.Content = *input.Content
	}

	if input.IsCompleted != nil {
		task.IsCompleted = *input.IsCompleted
	}

	err = app.model.Task.Update(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.model.Task.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelop{"message": "task successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
