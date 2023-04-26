package env

import (
	"bufio"
	"bytes"
	"os"

	"github.com/joho/godotenv"
)

func Write(file *os.File, key, value string) {
	// get file size
	stat, _ := file.Stat()

	// make array of bytes with size of file
	bs := make([]byte, stat.Size())

	// update array bs with file data
	bufio.NewReader(file).Read(bs)

	// parse file
	parsedEnv, _ := godotenv.Parse(bytes.NewBuffer(bs))

	// append data to save
	parsedEnv[key] = value

	// update file
	godotenv.Write(parsedEnv, "./.env")
}
