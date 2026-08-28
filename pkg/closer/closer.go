package closer

import (
	"io"
	"log"
)

func Close(closer io.Closer) {
	if closer == nil {
		return
	}
	err := closer.Close()
	if err != nil {
		log.Printf("Error closing item: %v", err)
	}
}
