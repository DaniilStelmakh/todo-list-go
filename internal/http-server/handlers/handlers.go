package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	apinext "github.com/DaniilStelmakh/go_final_project_main/apinext"
)

// Post запрос, ф-ия обработчик, добавляет задачу в БД
func AddTask(jobServ ServerJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req apinext.Task
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Некорректный формат запроса"})
			return
		}

		if req.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Неуказан заголовок задачи"})
			return
		}

		if req.Date == "" {
			req.Date = time.Now().Format(apinext.DateFormat)
		} else {
			_, err = time.Parse(apinext.DateFormat, req.Date)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат времени"})
				return
			}
		}

		if req.Repeat != "" {
			_, err = apinext.NextDate(time.Now(), req.Date, req.Repeat)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(ErrorResponse{Error: "Неверный формат повторений"})
				return
			}
		}

		id, err := jobServ.Add(req.Date, req.Title, req.Comment, req.Repeat)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseID{Id: id})
	}
}

// Delete запрос, ф-ия обработчик, удаляет задачу из БД
func DeleteTask(jobServ ServerJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсуетствует id"})
			return
		}
		if id, err := strconv.Atoi(id); err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "id должен быть положительным числом"})
			return
		}

		err := jobServ.Delete(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

// Post запрос, ф-ия обработчик, для обновления удаленной задачи в БД
func DoneTask(jobServ ServerJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсуетствует id"})
			return
		}
		if id, err := strconv.Atoi(id); err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "id должен быть положительным числом"})
			return
		}

		err := jobServ.Done(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}

// Get запрос, ф-ия обработчик, возвращает задачу из БД
func GetTask(jobServ ServerJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		id := r.URL.Query().Get("id")
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсуетствует id"})
			return
		}
		if id, err := strconv.Atoi(id); err != nil || id < 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "id должен быть положительным числом"})
			return
		}

		res, err := jobServ.Get(id)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}
		if res == nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Отсутвует task"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(*res)
	}
}

// Get запрос, ф-ия обработчик, для получения ближайщих задач
func GetTasks(jobServ ServerJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		var err error

		search := r.URL.Query().Get("search")

		tasks, err := jobServ.GetAll(search)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Response{Tasks: tasks})
	}
}

// Get запрос, ф-ия обаботчик, указывает на дату когда должна быть выполнена задача
func NextTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		now := r.URL.Query().Get("now")
		date := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")

		n, err := time.Parse(apinext.DateFormat, now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = time.Parse(apinext.DateFormat, date)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res, err := apinext.NextDate(n, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		w.Write([]byte(res))
	}
}

// Put запрос, ф-ия обработчик, исправляет задачу в БД
func UpdateTask(jobServ ServerJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req apinext.Task
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Некорректный формат запроса"})
			return
		}

		if ok, err := req.Valid(); !ok {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		err = ServerJob.Update(jobServ, &req)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{}{})
	}
}
