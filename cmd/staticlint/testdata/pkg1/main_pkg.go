package main

import "os"

func main() {
	os.DirFS("sss")
	os.Exit(1) // want "os.Exit call error"
}

func f() {
	os.Exit(1)
}
