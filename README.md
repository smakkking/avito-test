# как запускать сервис
чтобы запустить сервис
```
make -j 2 build-docker apply-migrations STORAGE=postgres
```

делать запросы к серверу лучше где-то через 30 секунд (потому что нужно время для применения миграций), но если вдруг docker-образы будут билдится дольше 30 секунд, то можно запустить сначала
```
make build-docker
```
Дождаться, пока сервис стартанет, а потом в отдельном терминале применить миграции
```
make apply-migrations
```

# проблемы и вопросы
Для авторизации через токены был выбран путь - захардкодить токены в коде самого сервиса (как самое простое решение)
Чтобы выполнить запрос с правами админа, нужно в хедере указать токен "admin", с правами пользователя - "user"