package client

import (
	"fmt"
	"github.com/helloyi/go-sshclient"
	"log"
	"os"
	"strings"
)

func SetupConnection(host, user, passwd string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}
	pubKey, err := os.ReadFile(homeDir + "/.ssh/id_rsa.pub")
	pubKey = []byte(strings.ReplaceAll(string(pubKey), "\n", ""))
	log.Println(homeDir + "/.ssh/id_rsa.pub")

	client, err := sshclient.DialWithPasswd(host+":22", user, passwd)
	if err != nil {
		log.Fatalln(err)
	}

	addKeyCmd := fmt.Sprintf("echo '%s' >> ~/.ssh/authorized_keys", pubKey)
	log.Printf("Выполнение команды: %s", addKeyCmd)
	addKeyOutput, err := client.Cmd(addKeyCmd).Output()
	if err != nil {
		log.Fatalf("Ошибка при добавлении ключа в authorized_keys: %v\n%s", err, string(addKeyOutput))
	}
	log.Printf("Результат команды add key: %s", string(addKeyOutput))

	defer client.Close()

}
