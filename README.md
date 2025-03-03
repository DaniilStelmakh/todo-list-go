# todo_list_go

## Описание проекта

Проект представляет собой Go веб-сервер, который реализует функциональность простейшего планировщика задач. Приложение позволяет создавать, редактировать, удалять и отмечать задачи как выполненные. В проекте используется REST API, все данные хранятся в базе данных SQLite. 

## Настройка tests/settings.go
```

var Port = 7540                     // порт, указанный в переменной окружения .env
var DBFile = "./scheduler.db"       // путь где храниться база данных
var FullNextDate = false            // флаг для проверки повторения задач в указанные дни 
var Search = false                  // флаг для проверки поисковая система по дате, названию задачи
var Token = ``                      // токен для авторизации 

```


## Структура проекта

```
go_final_project/
├── todo_app/
|   └── main.go              # Точка входа в приложение 
├── apinext/
│   └── apinext.go           # Модуль для проерки валидности запросов и формата времени  
├── internal/
│   └── http-server/
|   |     └── handlers.go    # Модуль в котором описаны все обработчики 
|   |     └── interface.go   
|   └── service/ 
|         └── task.go        # Модуль для обновления даты выполенения задач 
├── storage/            
│   ├── database/           
│   │    └── sqlite.go       # Модуль для работы с базой данных
├──  tests/                  # Тесты
│   └── settings.go          # Настройки для тестов
├── web/                     # Фронтенд
├── Dockerfile               # Файл для сборки Docker образа
├── go.mod                   # Модуль Go
├── go.sum                   # Хеши для зависимостей Go
├── README.md                # Документация
└── scheduler.db             # База данных SQLite  
```

## Инструкция по запуску кода локально

### Предварительные условия

- Установлен Go (версии 1.23.2 и выше)
- Установлен SQLite  или фреймворк, поддерживающий работу БД данного типа в вашей ОС

### Запуск приложения

1. Установите зависимости командой:

```sh
go mod download
```

2. Сборка приложения:

```sh
go build -o todo_app todo_app/main.go
```
3. Запуск приложения:

```sh
todo_app/main
```

Приложение будет доступно по адресу http://localhost:7540 (Либо на другом порте, который вы определите в переменной окружения)

### Запуск тестов

Для запуска тестов выполните, после сборки исполняемого файла (в новом терминале):

```sh
go test ./tests
```

## Cборка и запуск проекта через Docker

Выполните сборку Docker (в корне проекта) образа командой:

```sh
docker build -t todo_app .
```

### Запуск Docker контейнера

Запустите Docker контейнер командой:

```sh
docker run -d --name todo_server -p 7540:7540 todo_app
```

При сборке и запуске Docker образа вместо todo_app можете использовать своё название проекта.



## Спасибо за внимание
