package utilis

import "log"

func CheckError(err error) {
	// If an error is returned, print it to the console
	// and exit
	if err != nil {
		log.Fatal(err)
	}
}
