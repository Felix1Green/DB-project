# DB-project
Реализация проекта "Форумы" на курсе по базам данных.

### Запуск
В корне проекта:
```bash
$ sudo docker build -t prod_db_image .
$ sudo docker run -p 5000:5000 -p 5432:5432 prod_db_image
```
