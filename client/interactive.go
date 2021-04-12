package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/term"

	"bouchaud.org/legion/kubernetes/k8s-ldap-auth/types"
)

func readData(readLine func(screen io.ReadWriter) (string, error)) (string, error) {
	if !term.IsTerminal(0) || !term.IsTerminal(1) {
		return "", fmt.Errorf("stdin should be terminal")
	}

	oldState, err := term.MakeRaw(0)
	if err != nil {
		return "", err
	}
	defer term.Restore(0, oldState)

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}

	line, err := readLine(screen)
	if err != nil {
		return "", err
	}

	return line, nil

}

func username(screen io.ReadWriter) (string, error) {
	terminal := term.NewTerminal(screen, "username: ")

	line, err := terminal.ReadLine()
	if err == io.EOF || line == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return line, err
}

func password(screen io.ReadWriter) (string, error) {
	terminal := term.NewTerminal(screen, "")

	line, err := terminal.ReadPassword("password: ")
	if err == io.EOF || line == "" {
		return "", fmt.Errorf("password cannot be empty")
	}

	return line, err
}

func performAuth(addr, user, pass string) ([]byte, error) {
	var (
		err error
		res *http.Response
	)

	if user == "" {
		user, err = readData(username)
		if err != nil {
			return nil, err
		}
	}

	if pass == "" {
		pass, err = readData(password)
		if err != nil {
			return nil, err
		}
	}

	cred := types.Credentials{
		Username: user,
		Password: pass,
	}
	data, err := json.Marshal(cred)
	if err != nil {
		return nil, err
	}

	res, err = http.Post(addr, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
