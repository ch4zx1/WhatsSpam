#!/bin/bash

go get && go build 

DB_PATH="db"

if [ ! -d "$DB_PATH" ]; then
  mkdir -p $DB_PATH
fi
