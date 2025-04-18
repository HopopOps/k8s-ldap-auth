package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog/log"
	"github.com/zalando/go-keyring"
	"golang.org/x/term"

	"k8s-ldap-auth/types"
)

func readData(readLine func(screen io.ReadWriter) (string, error)) (string, error) {
	if !isatty.IsTerminal(os.Stdin.Fd()) && !isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		return "", fmt.Errorf("stdin should be terminal")
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
	terminal := term.NewTerminal(screen, "")

	print("username: ")

	line, err := terminal.ReadLine()
	if err == io.EOF || line == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	return line, err
}

func password(screen io.ReadWriter) (string, error) {
	terminal := term.NewTerminal(screen, "")

	print("password: ")

	line, err := terminal.ReadPassword("")

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

	interactiveMode := false

	if user == "" {
		log.Info().Msg("Username was not provided, asking for input")
		user, err = readData(username)
		print("\n")
		if err != nil {
			return nil, err
		}
	}
	log.Info().Str("username", user).Msg("Username exists.")

	if pass == "" {
		pass, err = keyring.Get(addr, user)
		if err != nil {
			log.Error().Err(err).Msg("Error while fetching credentials from store.")
		}
	}

	if pass == "" {
		interactiveMode = true
		log.Info().Msg("Password was not provided, asking for input")
		pass, err = readData(password)
		print("\n")
		if err != nil {
			return nil, err
		}
	}
	log.Info().Msg("Password exists.")

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

	if res.StatusCode != http.StatusOK {
		if err := keyring.Delete(addr, user); err != nil {
			log.Error().Err(err).Msg("Error while removing credentials from store.")
		}
		return nil, fmt.Errorf(http.StatusText(res.StatusCode))
	}

	var body []byte
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if interactiveMode {
		if err = keyring.Set(addr, user, pass); err != nil {
			log.Error().Err(err).Msg("Error while registering credentials into store.")
		}
	}

	return body, nil
}
