#!/bin/bash

DBSTRING="host=db user=$DBUSER password=$DBPASS dbname=calendar sslmode=disable"

goose --dir /migrations postgres "$DBSTRING" up