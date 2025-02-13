#!/bin/bash

DB_PATH="db"

if [ ! -d "$DB_PATH" ]; then
  mkdir -p $DB_PATH
fi


./miau-whatspam -media-path $DATA_PATH -history-path $HISTORY_PATH -request-full-sync false -debug false
