1. **Запуск Docker-сервисов**:
   Из корневой директории запустите Docker-контейнеры для Cassandra и RabbitMQ.
   ```bash
   docker-compose up -d
   ```

2. **Ожидание Загрузки Cassandra и RabbitMQ**:
   Убедитесь, что обе службы полностью функционируют, прежде чем продолжить.

3. **Развертывание Скрипта Создания Keyspace**:
   Скопируйте скрипт `init-cassandra.sh` в запущенный контейнер Cassandra и выполните его.
   ```bash
   docker cp init-cassandra.sh parser_hw-cassandra-1:/
   docker exec -it parser_hw-cassandra-1 /bin/bash -c "/init-cassandra.sh"
   ```

4. **Подготовка Сервиса Базы Данных Go**:
   Перейдите в директорию `go_db_service`.
   - Установите необходимые пакеты:
     ```bash
     go mod tidy
     ```
   - Соберите и запустите `main.go`:
     ```bash
     go build
     go run .
     ```

5. **Подготовка и Запуск Rust Parser**:
   Перейдите в директорию `rust_parser`.
   - Установите зависимости Rust:
     ```bash
     cargo install
     ```
   - Соберите Rust приложение:
     ```bash
     cargo build
     ```
   - Запустите Rust приложение:
     ```bash
     cargo run
     ```

6. **Доступ к Базе Данных Cassandra**:
   Подключитесь к базе данных Cassandra в Docker и проверьте данные.
   ```bash
   docker exec -it parser_hw-cassandra-1 cqlsh
   ```
   Внутри CQL shell используйте ваш keyspace и просмотрите содержимое таблицы `users`:
   ```sql
   USE my_keyspace;
   SELECT * FROM users;
   ```

7. **Интерфейс Управления RabbitMQ**:
   Вы можете мониторить очереди и сообщения RabbitMQ через веб-интерфейс, доступный по адресу `http://localhost:15672/#/`.

8. **Почему RabbitMQ**:
   <a href="https://github.com/iAmKoldyn/ms-s.github.io/wiki/Why-RabbitMq" target="_blank">Why RabbitMq</a>
