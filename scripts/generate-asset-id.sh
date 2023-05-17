#!/bin/sh

# Generates a random 7 character long string.
head /dev/urandom | LC_ALL=C tr -dc A-Za-z0-9 | head -c 7
