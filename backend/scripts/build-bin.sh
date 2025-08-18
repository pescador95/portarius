#!/bin/bash

APP_NAME=$1
APP_DIR=$2
BIN_DIR=$3

if [ -z "$APP_NAME" ] || [ -z "$APP_DIR" ] || [ -z "$BIN_DIR" ]; then
  echo "Uso: $0 <APP_NAME> <APP_DIR> <BIN_DIR>"
  exit 1
fi

mkdir -p "$BIN_DIR"

OUTPUT="$BIN_DIR/$APP_NAME.exe"

echo "Compilando $APP_NAME de $APP_DIR para $OUTPUT..."

GOOS=windows GOARCH=amd64 go build -o "$OUTPUT" "$APP_DIR"

if [ $? -eq 0 ]; then
  echo "✅ Executável gerado em: $OUTPUT"
else
  echo "❌ Erro ao compilar $APP_NAME"
  exit 1
fi
