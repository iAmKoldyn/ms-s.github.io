package main

import (
    "encoding/json"
    "log"
    "time"
    "fmt"

    "github.com/gocql/gocql"
    "github.com/streadway/amqp"
)

type User struct {
    URL  string `json:"url"`
    Name string `json:"name"`
    HTML string `json:"html"`
}

const (
    rabbitMQAddr        = "amqp://guest:guest@localhost:5672/"
    cassandraAddr       = "localhost"
    keyspace            = "my_keyspace"
    queueName           = "users_from_habr"
    rabbitMQMaxRetries  = 5
    cassandraMaxRetries = 5
    retryDelay          = 5 * time.Second
)

func main() {
    log.Println("Starting application...")

    conn, err := connectToRabbitMQ()
    if err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %s", err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatalf("Failed to open a channel: %s", err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(queueName, true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to declare a queue: %s", err)
    }

    msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
    if err != nil {
        log.Fatalf("Failed to register a consumer: %s", err)
    }

    session, err := connectToCassandra()
    if err != nil {
        log.Fatalf("Failed to connect to Cassandra: %s", err)
    }
    defer session.Close()

    log.Println(" [*] Waiting for messages. To exit press CTRL+C")
    for d := range msgs {
        var user User
        if err := json.Unmarshal(d.Body, &user); err != nil {
            log.Printf("Error decoding message: %s", err)
            continue
        }

        if err := session.Query(`INSERT INTO my_keyspace.users (id, url, name, html) VALUES (uuid(), ?, ?, ?)`, user.URL, user.Name, user.HTML).Exec(); err != nil {
            log.Printf("Error inserting user: %s", err)
        }
    }
}

func connectToRabbitMQ() (*amqp.Connection, error) {
    for i := 0; i < rabbitMQMaxRetries; i++ {
        conn, err := amqp.Dial(rabbitMQAddr)
        if err == nil {
            return conn, nil
        }
        log.Printf("Failed to connect to RabbitMQ, retrying in %s", retryDelay)
        time.Sleep(retryDelay)
    }
    return nil, fmt.Errorf("failed to connect to RabbitMQ")
}

func connectToCassandra() (*gocql.Session, error) {
    cluster := gocql.NewCluster(cassandraAddr)
    cluster.Keyspace = keyspace
    cluster.Consistency = gocql.Quorum

    for i := 0; i < cassandraMaxRetries; i++ {
        session, err := cluster.CreateSession()
        if err == nil {
            return session, nil
        }
        log.Printf("Failed to connect to Cassandra, retrying in %s", retryDelay)
        time.Sleep(retryDelay)
    }
    return nil, fmt.Errorf("failed to connect to Cassandra")
}
