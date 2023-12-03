## go-api-docker-example

### Docker commands

```bash
docker compose build
docker compose up
```

### Env file

```bash
PORT=:5000
DATABASE_URL=<user>:<password>@tcp(<mysql-container>:3306)/<database>?parseTime=true
DB_DRIVER=mysql
MYSQL_RANDOM_ROOT_PASSWORD=<password>
MYSQL_DATABASE=<database>
MYSQL_USER=<user>
MYSQL_PASSWORD=<password>
```

### API commands

1. Fetching todos from DB

```bash
curl -X GET 'http://localhost:5000/todos'
```

2. Inserting todos in DB

```bash
curl -X POST 'http://localhost:5000/todos' -d '{"task": "Learn Go", "completed": false}'
```
