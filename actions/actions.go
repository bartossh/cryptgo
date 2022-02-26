package actions

import (
	"errors"
	"fmt"
	"io"
	"os/user"
	"path/filepath"

	"github.com/urfave/cli/v2"

	"github.com/bartossh/cryptgo/ciphers"
	"github.com/bartossh/cryptgo/filesbuf"
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
	// Generate flag
	Generate = "generate"
	// Use flag
	Use = "use"
)

type (
	DataPiper interface {
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
		input, output string
		e, d          DataPiper
		rwc           readWriteCloser
	}
	// CommandFactory instance allows to create command
	CommandFactory struct{}
)

// NewCommandFactory creates new factory instance for creating command
func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

// SetEncrypter sets encrypt direction command
func (cf *CommandFactory) SetEncrypter(c *cli.Context) (*Command, error) {
	uh, err := userHome()
	if err != nil {
		return nil, err
	}

	cmd := &Command{}
	cmd.input = c.String(Input)
	if cmd.input == "" {
		return nil, errors.New("input file path is not specified")
	}
	cmd.output = c.String(Output)
	if cmd.output == "" {
		return nil, errors.New("output file path is not specified")
	}

	rwAgent := &filesbuf.Agent{}

	var bufPriv []byte
	if rsaPath := c.String(Generate); rsaPath != "" {
		priv, err := ciphers.GeneratePrivateKey()
		if err != nil {
			return nil, fmt.Errorf("cannot generate rsa new key, %w", err)
		}
		bufPriv = ciphers.EncodePrivateKeyToPEM(priv)
		w, err := rwAgent.GetWriteCloser(rsaPath)
		if err != nil {
			return nil, fmt.Errorf("cannot create file for path %s, %w", rsaPath, err)
		}
		defer w.Close()
		if _, err := w.Write(bufPriv); err != nil {
			return nil, fmt.Errorf("canoot write rsa key of pem fromat to file %s, %w", rsaPath, err)
		}
		if err := cmd.setCmdEncrypt(bufPriv, []byte{}, rwAgent); err != nil {
			return nil, fmt.Errorf("encryptor initializetion failed, %w", err)
		}
		return cmd, nil
	}

	rcPriv, err := rwAgent.GetReadCloser(filepath.Join(uh, unixPrivRSAPath))
	if err != nil {
		return nil, err
	}
	defer rcPriv.Close()
	bufPriv, err = io.ReadAll(rcPriv)
	if err != nil {
		return nil, fmt.Errorf("cannot set cipher, reading private key failed, %w", err)
	}

	var passwd []byte
	passwdStr := c.String(Passwd)
	if passwdStr != "" {
		passwd = []byte(passwdStr)
	}
	if err := cmd.setCmdEncrypt(bufPriv, passwd, rwAgent); err != nil {
		return nil, fmt.Errorf("encryptor initialization failed, %w", err)
	}
	return cmd, nil
}

func (cmd *Command) setCmdEncrypt(bufPriv, passwd []byte, rwAgent readWriteCloser) error {
	e, err := ciphers.NewEncrypt(bufPriv, passwd)
	if err != nil {
		return err
	}
	cmd.e = e
	cmd.rwc = rwAgent
	return nil
}

// SetDecrypter sets decrypt direction command
func (cf *CommandFactory) SetDecrypter(c *cli.Context) (*Command, error) {
	uh, err := userHome()
	if err != nil {
		return nil, err
	}

	cmd := &Command{}
	cmd.input = c.String(Input)
	if cmd.input == "" {
		return nil, errors.New("input file path is not specified")
	}
	cmd.output = c.String(Output)
	if cmd.output == "" {
		return nil, errors.New("output file path is not specified")
	}

	var passwd []byte
	passwdStr := c.String(Passwd)
	if passwdStr != "" {
		passwd = []byte(passwdStr)
	}

	rsaPath := filepath.Join(uh, unixPrivRSAPath)
	if rp := c.String(Use); rp != "" {
		rsaPath = filepath.Join(rp)
		passwd = []byte{}
	}

	rwAgent := &filesbuf.Agent{}
	rcPriv, err := rwAgent.GetReadCloser(rsaPath)
	if err != nil {
		return nil, err
	}
	defer rcPriv.Close()

	bufPriv, err := io.ReadAll(rcPriv)
	if err != nil {
		return nil, fmt.Errorf("cannot set cipher, reading private key failed, %s", err)
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
func (cmd *Command) Decrypt() error {
	rc, err := cmd.rwc.GetReadCloser(cmd.input)
	defer rc.Close()
	if err != nil {
		return err
	}
	wc, err := cmd.rwc.GetWriteCloser(cmd.output)
	defer wc.Close()
	if err != nil {
		return err
	}
	return cmd.d.Pipe(rc, wc)
}

// Encrypt runs decryption
func (cmd *Command) Encrypt() error {
	rc, err := cmd.rwc.GetReadCloser(cmd.input)
	if err != nil {
		return err
	}
	defer rc.Close()
	wc, err := cmd.rwc.GetWriteCloser(cmd.output)
	if err != nil {
		return err
	}
	defer wc.Close()
	return cmd.e.Pipe(rc, wc)
}

func userHome() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}
