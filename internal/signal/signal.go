// Модуль управления сигналами завершения программы
package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// Signal управляет сигналом для завершением программы
func Signal() chan os.Signal {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	return termChan
}
