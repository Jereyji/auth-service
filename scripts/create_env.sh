#!/bin/bash

ENV_FILE="deployments/.env"

POSTGRES_HOST="postgres"
POSTGRES_PORT="5432"
POSTGRES_DB="auth-service_db"
POSTGRES_USER="jerey"
POSTGRES_PASSWORD="jereyjerey"
POSTGRES_SSLMODE="disable"

SECRET_KEY=$(openssl rand -hex 32)
EMAIL="soelnstu@gmail.com"
EMAIL_PASSWORD="icto xvum jcom yiny"

REDIS_HOST="redis"
REDIS_PORT="6379"
REDIS_DB="0"
REDIS_PASSWORD="jereyjerey"

echo "POSTGRES_HOST=${POSTGRES_HOST}" > $ENV_FILE
echo "POSTGRES_PORT=${POSTGRES_PORT}" >> $ENV_FILE
echo "POSTGRES_DB=${POSTGRES_DB}" >> $ENV_FILE
echo "POSTGRES_USER=${POSTGRES_USER}" >> $ENV_FILE
echo "POSTGRES_PASSWORD=${POSTGRES_PASSWORD}" >> $ENV_FILE
echo "POSTGRES_SSLMODE=${POSTGRES_SSLMODE}" >> $ENV_FILE

echo "SECRET_KEY=${SECRET_KEY}" >> $ENV_FILE
echo "EMAIL=${EMAIL}" >> $ENV_FILE
echo "EMAIL_PASSWORD=${EMAILPASSWORD}" >> $ENV_FILE

echo "REDIS_HOST=${REDIS_HOST}" >> $ENV_FILE
echo "REDIS_PORT=${REDIS_PORT}" >> $ENV_FILE
echo "REDIS_PASSWORD=${REDIS_PASSWORD}" >> $ENV_FILE
echo "REDIS_DB=${REDIS_DB}" >> $ENV_FILE