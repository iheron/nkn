package password

import (
	"errors"
	"fmt"
	"github.com/nknorg/nkn/util/config"
	"os"

	"github.com/howeyc/gopass"
)

var (
	Passwd string
)

// GetPassword gets password from user input
func GetPassword() ([]byte, error) {
	fmt.Printf("Password:")
	passwd, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	return passwd, nil
}

// GetConfirmedPassword gets double confirmed password from user input
func GetConfirmedPassword() ([]byte, error) {
	fmt.Printf("Password:")
	first, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}

	if len(first) == 0 {
		fmt.Println("password is invalid.")
		os.Exit(1)
	}

	fmt.Printf("Re-enter Password:")
	second, err := gopass.GetPasswd()
	if err != nil {
		return nil, err
	}
	if len(first) != len(second) {
		fmt.Println("Unmatched Password")
		os.Exit(1)
	}
	for i, v := range first {
		if v != second[i] {
			fmt.Println("Unmatched Password")
			os.Exit(1)
		}
	}
	return first, nil
}

// GetAccountPassword gets the node's wallet password from the command line,
// the NKN_WALLET_PASSWORD environment variable, or user input, in this order
func GetAccountPassword() ([]byte, error) {
	if Passwd == "" {
		Passwd = os.Getenv("NKN_WALLET_PASSWORD")
	}
	if !config.Parameters.WebGuiCreateWallet {
		if Passwd == "" {
			return GetPassword()
		}
	}
	if Passwd == "" {
		return []byte(Passwd), errors.New("wait for set password")
	}
	return []byte(Passwd), nil
}
