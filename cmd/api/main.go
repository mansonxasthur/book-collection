package main

func main() {
	app := NewApp(":8080")

	app.Bootstrap()
	app.Run()
}
