package actions

import (
	"errors"
	"fmt"
	"io"
	"os/user"
	"path/filepath"

	"github.com/bartossh/cryptgo/ciphers"
	"github.com/bartossh/cryptgo/filesbuf"
	"github.com/urfave/cli/v2"
)

const (
	unixPrivRSAPath = ".ssh/id_rsa"
)

const (
	// Input flag
	Input = "input"
	// Output flag
	Output = "output"
	// Passwd flag
	Passwd = "passwd"
)

type (
	crypto interface {
		Pipe(rd io.Reader, wr io.Writer) error
	}
	readCloser interface {
		GetReadCloser(path string) (io.ReadCloser, error)
	}
	writeCloser interface {
		GetWriteCloser(path string) (io.WriteCloser, error)
	}
	readWriteCloser interface {
		readCloser
		writeCloser
	}
)

type (
	// Command instance allows to run ActionFunc
	Command struct {
		rwc readWriteCloser
		e   crypto
		d   crypto
	}
	// CommandFactory instance allows to create command
	CommandFactory struct{}
)

// NewCommandFactory creates new factory instance for creating command
func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

// SetEncryptor sets decrypt or encrypt direction
func (cf *CommandFactory) SetEncryptor(c *cli.Context) (*Command, error) {
	uh, err := userHome()
	if err != nil {
		return nil, err
	}
	cmd := &Command{}
	rwAgent := &filesbuf.Agent{}
	rcPriv, err := rwAgent.GetReadCloser(filepath.Join(uh, unixPrivRSAPath))
	if err != nil {
		return nil, err
	}
	defer rcPriv.Close()
	bufPriv, err := io.ReadAll(rcPriv)
	if err != nil {
		return nil, fmt.Errorf("cannot set cipher, reading private key failed, %s", err)
	}
	var passwd []byte
	passwdStr := c.String(Passwd)
	if passwdStr != "" {
		passwd = []byte(passwdStr)
	}
	e, err := ciphers.NewEncrypt(bufPriv, passwd)
	if err != nil {
		return nil, err
	}
	cmd.e = e
	cmd.rwc = rwAgent
	return cmd, nil
}

// SetDecryptor sets decrypt or encrypt direction
func (cf *CommandFactory) SetDecryptor(c *cli.Context) (*Command, error) {
	uh, err := userHome()
	if err != nil {
		return nil, err
	}
	cmd := &Command{}
	rwAgent := &filesbuf.Agent{}
	rcPriv, err := rwAgent.GetReadCloser(filepath.Join(uh, unixPrivRSAPath))
	if err != nil {
		return nil, err
	}
	defer rcPriv.Close()
	bufPriv, err := io.ReadAll(rcPriv)
	if err != nil {
		return nil, fmt.Errorf("cannot set cipher, reading private key failed, %s", err)
	}
	var passwd []byte
	passwdStr := c.String(Passwd)
	if passwdStr != "" {
		passwd = []byte(passwdStr)
	}
	d, err := ciphers.NewDecrypt(bufPriv, passwd)
	if err != nil {
		return nil, err
	}
	cmd.d = d
	cmd.rwc = rwAgent
	return cmd, nil
}

// Decrypt runs decryption
func (cmd *Command) Decrypt(c *cli.Context) error {
	inp := c.String(Input)
	if inp == "" {
		return errors.New("input file path is not specified")
	}
	out := c.String(Output)
	if out == "" {
		return errors.New("output file path is not specified")
	}

	rc, err := cmd.rwc.GetReadCloser(inp)
	defer rc.Close()
	if err != nil {
		return err
	}
	wc, err := cmd.rwc.GetWriteCloser(out)
	defer wc.Close()
	if err != nil {
		return err
	}
	return cmd.d.Pipe(rc, wc)
}

// Encrypt runs decryption
func (cmd *Command) Encrypt(c *cli.Context) error {
	inp := c.String(Input)
	if inp == "" {
		return errors.New("input file path is not specified")
	}
	out := c.String(Output)
	if out == "" {
		return errors.New("output file path is not specified")
	}

	rc, err := cmd.rwc.GetReadCloser(inp)
	defer rc.Close()
	if err != nil {
		return err
	}
	wc, err := cmd.rwc.GetWriteCloser(out)
	defer wc.Close()
	if err != nil {
		return err
	}
	return cmd.e.Pipe(rc, wc)
}

func userHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
