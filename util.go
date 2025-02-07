package resty

import "log"

func FailOnError(err error) {
	if err != nil {
		log.Panicf("Fatal error: %v", err)
	}
}
