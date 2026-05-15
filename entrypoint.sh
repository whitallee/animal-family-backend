#!/bin/sh
set -e

./bin/migrate up
exec ./bin/server
