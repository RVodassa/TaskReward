## Настройка переменных окружения

Создайте файл `.env` в корне проекта и укажите следующие переменные:

```env
DB_HOST=db
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=appdb
JWT_SECRET=your_jwt_secret_key
```

Если вы хотите использовать другой файл (например, `.env.prod`), вы можете указать его с помощью флага `--env-file`:
```bash
docker-compose --env-file .env.prod up
```
Запустите приложение с помощью Docker Compose:
```bash
docker-compose up --build
```

Чтобы подергать ручки и посмотреть документацию перейдите по адресу:
```http://localhost:8080/swagger``` 

