package main

import (
	"fmt"
	"github.com/rmorriso/genstr"
)

func main() {
	fmt.Printf("Password: %s\n", genstr.Simple(12))
}

