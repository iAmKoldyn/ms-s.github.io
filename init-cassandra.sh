#!/bin/bash

echo "Creating keyspace..."
cqlsh -e "CREATE KEYSPACE IF NOT EXISTS my_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};"
echo "Keyspace created successfully."

echo "Creating 'users' table in my_keyspace..."
cqlsh -e "
CREATE TABLE IF NOT EXISTS my_keyspace.users (
    id UUID PRIMARY KEY,
    url TEXT,
    name TEXT,
    html TEXT
);"
echo "'users' table created successfully."
