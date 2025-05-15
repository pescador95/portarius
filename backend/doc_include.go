//go:build docs

package main

import (
	_ "portarius/internal/inventory/handler"
	_ "portarius/internal/package/handler"
	_ "portarius/internal/reminder/handler"
	_ "portarius/internal/reservation/handler"
	_ "portarius/internal/resident/handler"
	_ "portarius/internal/user/handler"
	_ "portarius/internal/whatsapp/handler"
)
