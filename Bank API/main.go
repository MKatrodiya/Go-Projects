package main

func main() {
	server := APIServer{
		listenPort: ":8080",
	}
	server.Run()
}
