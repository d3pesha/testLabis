# Начало работы

Для начала работы выполните следующие шаги:

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/d3pesha/testLabis.git
   ```

2. Создайте файл `.env` в корне вашего проекта с следующим содержимым:

   ```env
   DB_HOST=host.docker.internal
   DB_PORT=5432
   DB_USERNAME=postgres
   DB_PASSWORD=postgres
   DB_NAME=labis
   ```

   Обновите значения в соответствии с конфигурацией вашей базы данных.

3. Создайте базу данных используя файл db.txt из папки database.

4. Соберите Docker-образ:

   ```bash
   docker build -t labis .
   ```

5. Запустите Docker-контейнер:

   ```bash
   docker run -p 8080:8080 labis
   ```

## Использование

- Перейдите в приложение по адресу [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html).

### Создание нового клиента

Создайте нового клиента с указанной информацией.

- **URL:** `/clients`
- **Метод:** `POST`
- **Тело запроса:**
  
  ```json
  {
    "email": "string",
    "fio": "string",
    "password": "string"
  }
  ```

## Договоры

### Создание нового договора

Создайте новый договор, связанный с указанным объектом.

- **URL:** `/contracts`
- **Метод:** `POST`
- **Тело запроса:**
  
  ```json
  {
    "object_id": 0
  }
  ```

## Объекты

### Создание нового объекта

Создайте новый объект, связанный с указанным пользователем.

- **URL:** `/objects`
- **Метод:** `POST`
- **Тело запроса:**
  
  ```json
  {
    "id_user": 0,
    "address": "string",
    "is_visible": true
  }
  ```

--- 
