package main

func main() {
	go connectToAtmotube()
	go startHTTPServer()
	select {}
}
