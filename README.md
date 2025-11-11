# QA API Service

API-сервис для вопросов и ответов

Сервис поддерживает:
- Создание и получение вопросов.
- Добавление ответов к вопросам.
- Каскадное удаление: при удалении вопроса все ответы к нему удаляются автоматически.

## Запуск

```bash
git clone https://github.com/EugeneKrivoshein/qa_api_service.git
cd qa_api_service
go mod tidy
docker-compose up --build

go test ./tests -count=1
```
Сервер будет доступен по адресу: 

```bash
http://localhost:8080
```

Swagger документация 

```bash
http://localhost:8080/swagger/index.html
```
