package main

import "fmt"

func pad(text, left, right string) string {
	return fmt.Sprintf("%s%s%s", left, text, right)
}
