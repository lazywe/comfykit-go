package comfykit

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "[ComfyKit] ", log.LstdFlags)

func Info(format string, v ...interface{}) {
	logger.Printf("[INFO] "+format, v...)
}

func Warn(format string, v ...interface{}) {
	logger.Printf("[WARN] "+format, v...)
}

func Error(format string, v ...interface{}) {
	logger.Printf("[ERROR] "+format, v...)
}

func Debug(format string, v ...interface{}) {
	logger.Printf("[DEBUG] "+format, v...)
}
