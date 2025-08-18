APP_NAME := portarius
BIN_DIR := C:\\apps
EXE := $(BIN_DIR)/$(APP_NAME).exe

APP_DIR ?= ../

APP_PATH := C:\\apps\\$(APP_NAME).exe
APP_DISPLAY := Portarius

# ===== Targets =====

.PHONY: all build startup clean

all: build startup

build:
	@echo "==> Compilando para Windows..."
	@bash build-bin.sh $(APP_NAME) $(APP_DIR) $(BIN_DIR)

startup:
	@echo "==> Adicionando app à inicialização automática do Windows..."
	@powershell -ExecutionPolicy Bypass -File add-startup.ps1 -AppPath "$(APP_PATH)" -AppName "$(APP_DISPLAY)"
	@echo "==> App adicionado à inicialização automática."

clean:
	@echo "==> Limpando executáveis..."
	@rm -f $(EXE)
	@powershell -ExecutionPolicy Bypass -File remove-startup.ps1 -AppName "$(APP_DISPLAY)"
	@echo "==> Limpeza concluída."