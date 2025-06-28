package pkg

import "os"

func main() {
	os.DirFS("sss")
	os.Exit(1)
}

func f() {
	os.Exit(1)
}
