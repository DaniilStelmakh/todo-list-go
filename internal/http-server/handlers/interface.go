package handlers

import apinext "github.com/DaniilStelmakh/go_final_project_main/apinext"

// Интерфейс для реализации обработчиков
type ServerJob interface {
	Add(date, title, comment, repeat string) (int, error)
	Delete(id string) error
	Done(id string) error
	Get(id string) (*apinext.Task, error)
	GetAll(search string) ([]apinext.Task, error)
	Update(*apinext.Task) error
}

// Структура для ответа в JSON для ближайщих задач
type Response struct {
	Tasks []apinext.Task `json:"tasks"`
}

// Структура для ответа в JSON при добавлении задачи
type ResponseID struct {
	Id int `json:"id"`
}

// Структура для ошибки в JSON формате
type ErrorResponse struct {
	Error string `json:"error"`
}
