package log

import (
	"fmt"
	"time"
	"strings"
)

const RED =    "\033[31m"
const YELLOW = "\033[33m"
const GREEN =  "\033[32m"
const RESET =  "\033[0m"


func Space() {
	fmt.Print("\n\n")
}

func now() string {
	time := time.Now()
	return time.Format("2006-01-02 15:04:05")
}

func message(msg string, args []interface{}) (string, []interface{}) {
	marker := []rune("%")[0]
	nargs := 0
	for _, char := range msg {
		if char == marker {
			nargs++
		}
	}

	msgargs := args[0:nargs]
	args = args[nargs:]
	msg = fmt.Sprintf(msg, msgargs...)
	msg = strings.TrimSpace(msg)
	return msg, args
}

func Log(group, msg string, args ...interface{}) {
	msg, args = message(msg, args)
	line := fmt.Sprintf("%-20s %s%-12s%s %s", now(), GREEN, "["+group+"]", RESET, msg)
	lineargs := append([]interface{}{line},args...)
	fmt.Println(lineargs...)
}

func Warn(group, msg string, args ...interface{}) {
	msg, args = message(msg, args)
	line := fmt.Sprintf("%-20s %s%-12s%s %s", now(), YELLOW, "["+group+"]", RESET, msg)
	lineargs := append([]interface{}{line},args...)
	fmt.Println(lineargs...)
}

func Error(group, msg string, args ...interface{}) {
	msg, args = message(msg, args)
	line := fmt.Sprintf("%-20s %s%-12s%s %s", now(), RED, "["+group+"]", RESET, msg)
	lineargs := append([]interface{}{line},args...)
	fmt.Println(lineargs...)
}