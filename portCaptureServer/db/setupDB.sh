#!/usr/bin/bash
export PGPASSWORD=password
psql -U postgres -h 0.0.0.0 -p 5433   -c '\ir schema.sql'
