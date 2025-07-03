package utils

import (
	"log"
	"time"
)

func TimeCallback(callbackName string, fn func()) {
	start := time.Now()
	fn()
	duration := time.Since(start)
	log.Printf("Callback %s took %v", callbackName, duration)
}

func TimeCallbackWithReturn[T any](callbackName string, fn func() T) T {
	start := time.Now()
	result := fn()
	duration := time.Since(start)
	log.Printf("Callback %s took %v", callbackName, duration)
	return result
}