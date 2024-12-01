package utilis

import (
	"bufio"
	"errors"
	"esercizioSDCC/configuration"
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"os"
)

func GenerateRandomIntsFIle() (string, error) {
	file, err := os.Create(configuration.FILE_NAME)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return "", errors.New("error creating file")
	}
	defer func(file *os.File) {
		err = file.Close()
	}(file)

	for i := 0; i < configuration.FILE_SIZE; i++ {
		stringValue := fmt.Sprintf("%d", rand.IntN(configuration.MAX_RANDOM_NUMBER))
		if i != configuration.FILE_SIZE-1 {
			stringValue = stringValue + " "
		}
		_, err = file.WriteString(stringValue)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return "", err
		}
	}
	return configuration.FILE_NAME, err
}

func ReadInts(buffer *bufio.Reader, amount int) (string, error) {
	read := 0
	testo := ""
	for {
		char, _, err := buffer.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				read++
				if read < amount {
					return "", errors.New("Not enough entries to read")
				}
				return testo, nil
			}

		}
		if char == ' ' {
			read++
		}
		testo += string(char)
		if read == amount {
			return testo, nil
		}
	}
}

func GetStrings(nome string) []string {
	valuesForHost := int(math.Round(configuration.FILE_SIZE / configuration.HOSTS_NUMBER))
	var err error
	var arrayStrings []string
	file, err := os.Open(nome)
	CheckError(err)
	defer func(file *os.File) {
		err = file.Close()
	}(file)
	remaining := configuration.FILE_SIZE
	reader := bufio.NewReader(file)
	for i := 0; i < configuration.HOSTS_NUMBER; i++ {
		var testo string
		if remaining-2*valuesForHost < 0 {
			valuesForHost = remaining
		}
		testo, err = ReadInts(reader, valuesForHost)
		CheckError(err)
		remaining -= valuesForHost
		arrayStrings = append(arrayStrings, testo)
	}
	if remaining != 0 {
		log.Fatal("mismatch on numbers read")
	}
	CheckError(err)
	return arrayStrings
}
