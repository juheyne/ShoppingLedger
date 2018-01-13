package main

func main() {
	a := App{}

	a.Initialize("./ledger.db")
	a.Run(":8080")
}
