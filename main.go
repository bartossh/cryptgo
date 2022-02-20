package main

import (
	"os"
	"runtime"
	"sort"

	"examples.com/cryptgo/actions"
	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

func main() {
	logoRender()
	if runtime.GOOS == "windows" {
		pterm.Error.Printfln("cannot run on windows yet")
	}
	app := &cli.App{
		Name:      "cryptgo",
		Usage:     "simple and easy file encryption",
		Copyright: "(c) 2022 Bartossh",
		HelpName:  "cryptgo",
		Version:   "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    actions.Input,
				Aliases: []string{"i"},
				Usage:   "path to input name",
			},
			&cli.StringFlag{
				Name:    actions.Output,
				Aliases: []string{"o"},
				Usage:   "path to output name",
			},
			&cli.StringFlag{
				Name:    actions.Passwd,
				Aliases: []string{"p"},
				Usage:   "passphrase for rsa key",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "encrypt",
				Usage: "encrypts file with provided password",
				Action: func(c *cli.Context) error {
					spnr, _ := pterm.DefaultSpinner.Start("Encrypting...")
					ac, err := actions.NewCommandFactory().SetEncryptor(c)
					if err != nil {
						spnr.Fail("Encryption failed!")
						return err
					}
					if err := ac.Encrypt(c); err != nil {
						spnr.Fail("Encryption failed!")
						return err
					}
					spnr.Success("File encrypted.")
					return nil
				},
			},
			{
				Name:  "decrypt",
				Usage: "decrypts file with provided password",
				Action: func(c *cli.Context) error {
					spnr, _ := pterm.DefaultSpinner.Start("Decrypting...")
					ac, err := actions.NewCommandFactory().SetDecryptor(c)
					if err != nil {
						spnr.Fail("Decryption failed!")
						return err
					}
					if err := ac.Decrypt(c); err != nil {
						spnr.Fail("Decryption failed!")
						return err
					}
					spnr.Success("File decrypted")
					return nil
				},
			},
		},
	}
	app.EnableBashCompletion = true

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		pterm.Error.Printf("%s", err)
		return
	}

}

func logoRender() {
	pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("crypt", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("go", pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
}
