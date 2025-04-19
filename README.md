**Hello! This is a URL shortener written in Go (Golang).**

**Structure:**

<img width="1447" alt="Структура" src="https://github.com/user-attachments/assets/bc310b41-c67f-4ec7-80b5-c63cf810613a" />

**.env file:**
```
DSN="host=localhost user=postgres password=my_pass dbname=link port=5432 sslmode=disable"
SECRET="/2+XnmJGz1j3ehIVI/5P9kl+CghrE3DcS7rnT+qar5w="
```

**Requests:**

1)  POST http://localhost:8081/auth/register
```
{
    "email": "",
    "password": "",
    "name": ""
}
```
2) POST http://localhost:8081/auth/login
```
{
    "email": "d111@d.ru",
    "password": "1"
}
```
3) POST http://localhost:8081/link
```
Authorization Bearer token
{
    "url": ""
}
```
4) PATCH http://localhost:8081/link/id
```
Authorization Bearer token
{
    "url": ""
}
```
5) DELETE http://localhost:8081/link/id
```
Authorization Bearer token
```
6) GET http://localhost:8081/hash
   
7) GET http://localhost:8081/link?limit=&offset=
```
Authorization Bearer token
```
8) http://localhost:8081/stat?from=yyyy-mm-dd&to=yyyy-mm-dd&by=month
```
Authorization Bearer token
```
