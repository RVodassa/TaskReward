# Простой HTTP сервер для управления пользователями на Go 

## Описание проекта

Этот проект представляет собой простой HTTP сервер, реализованный на языке Go,
для управления пользователями. Сервер предоставляет следующие возможности:
- Создание пользователей с указанием реферального кода при регистрации.
- Отслеживание статуса пользователей.
- Вознаграждение пользователей за выполнение заданий.
- Получение информации о лидерах по балансу.

#### Проект использует современный стэк технологий:
- Docker/Docker Compose — для контейнеризации и управления сервисами.
- Go-Chi — легкий и мощный роутер для создания API.
- PostgreSQL — для хранения информации.
- Golang-migrate для управления миграциями.
- Go 1.23.3 — последняя версия языка Go для разработки сервера.
- Swagger — для автоматической генерации документации API.

## Быстрый старт
#### 1. Клонируйте репозиторий

Для начала работы склонируйте репозиторий на ваш компьютер:

```bash
git clone https://github.com/RVodassa/TaskReward.git
cd TaskReward
```

#### 2. Настройка окружения

Создайте файл .env в корневой директории проекта для настройки переменных окружения. Пример содержимого файла:
```env
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=appdb
DB_SSL=disable
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRATION=1h
SERVER_PORT:8080
```
Убедитесь что на вашем хостинге свободен порт указанный в SERVER_PORT

#### 3. Запуск приложения

Запустите приложение с помощью Docker Compose:

```bash
sudo docker-compose up --build
```

Если вам нужно использовать другой файл окружения (например, для production), укажите его с помощью флага --env-file:

```bash
docker-compose --env-file .env.prod up
```

#### 4. Документация API

После запуска сервера документация API будет доступна по адресу:
http://localhost:8080/swagger

#### Как пользоваться документацией ?
1. Создайте пользователя в RegisterUser
2. Получите JWToken выполнив аутентификацию в login
3. Вставьте токен в окно Authorization с приставкой 'Bearer', пример: 'Bearer your_token_jwt'.
4. Теперь вам доступен полный список ф-ций в течении жизни токена (укажите свое время в .env 'JWT_EXPIRATION') 
5. При запуске приложении, в базу данных добавилось 10 задач, попробуйте поиграть с ними.

#### Что продумано:
Основное: 
- Чистая архитектура, ООП, принципы SOLID и чистый код.
- Документация swagger.
- Передача и работа с контекстом
- Graceful shutdown для мягкого завершения работы сервера.
- Хеширование пароля с использованием bcrypt.
- healthcheck в docker-compose, чтобы наш сервис дождался базу данных.
- Трассировка ошибок, логирование в случае внутренней ошибки.
- Транзакции при работе с хранилищем.


## Что можно улучшить?
Основное:
- Refresh tokens для обновления access tokens
- Unit-tests/integration-tests конечно да!
- Поддержка файлов конфигурации для более гибкой настройки приложения
- Структурное логирование
- Гибкость ф-ций(например пагинация для списков)

### Если у вас есть вопросы, свяжитесь:
Email: assadov.spb@bk.ru
Telegram: @LilDragxn