package main

func main() {
	a := App{}

	a.Initialize("./foo.db")
	a.Run(":8080")
}
