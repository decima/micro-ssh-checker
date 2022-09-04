package main

import (
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
)

const RsaPath = "./id_rsa"
const logPath = "./users.log"

var port = flag.Int("p", 2222, "help message for flag n")

var logFile *os.File

//go:embed motd.txt
var motd string

func init() {

	_, err := os.Stat(RsaPath)
	log.Printf("checking for %v file\n", RsaPath)

	if errors.Is(err, os.ErrNotExist) {
		log.Printf("%v not found\n", RsaPath)
		bitSize := 4096

		privateKey, err := generatePrivateKey(bitSize)
		if err != nil {
			log.Fatal(err.Error())
		}

		privateKeyBytes := encodePrivateKeyToPEM(privateKey)
		err = writeKeyToFile(privateKeyBytes, RsaPath)
		if err != nil {
			log.Fatal(err.Error())
		}

	}

	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Could not open example.txt")
		return
	}
}

func main() {
	defer logFile.Close()

	flag.Parse()

	ssh.Handle(func(s ssh.Session) {
		if s.User() == "decima" {
			_, _ = io.WriteString(s, fmt.Sprintf("Users :\n"))
			content, _ := ioutil.ReadFile(logPath)
			_, _ = io.WriteString(s, string(content))
			return
		}
		motdGenerate(s)

		newUser := fmt.Sprintf(" * %v (%v) – %v\n", s.User(), strings.Join(s.Command(), " "), s.RemoteAddr().String())
		_, _ = logFile.WriteString(newUser)
	})

	addr := fmt.Sprintf(":%v", *port)

	log.Printf("Starting on address %v", addr)
	log.Fatal(ssh.ListenAndServe(addr, nil, ssh.HostKeyFile(RsaPath)))
}

func motdGenerate(s ssh.Session) {

	tmpl, err := template.New("test").Parse(motd)
	if err != nil {
		log.Println(err)
		return
	}

	err = tmpl.Execute(s, s)
	if err != nil {
		log.Println(err)
	}

}
