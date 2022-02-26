package main

import (
	"os"
	"runtime"
	"sort"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"

	"github.com/bartossh/cryptgo/actions"
)

func main() {
	logoRender()
	if runtime.GOOS == "windows" {
		pterm.Error.Printfln("cannot run on windows yet")
	}
	app := &cli.App{
		Name:      "cryptgo",
		Usage:     "simple and easy file encryption / decryption",
		Copyright: "(c) 2022 Bartossh",
		HelpName:  "cryptgo",
		Version:   "v0.0.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    actions.Input,
				Aliases: []string{"i"},
				Usage:   "path to input file",
			},
			&cli.StringFlag{
				Name:    actions.Output,
				Aliases: []string{"o"},
				Usage:   "path to output file",
			},
			&cli.StringFlag{
				Name:    actions.Passwd,
				Aliases: []string{"p"},
				Usage:   "passphrase for user .ssh rsa key, doesn't work with -generate and -use",
			},
			&cli.StringFlag{
				Name:    actions.Generate,
				Aliases: []string{"g"},
				Usage:   "path where to generate fresh rsa key",
			},
			&cli.StringFlag{
				Name:    actions.Use,
				Aliases: []string{"u"},
				Usage:   "path to specific rsa key to be used",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "encrypt",
				Usage: "encrypts file with provided password",
				Action: func(c *cli.Context) error {
					spnr, err := pterm.DefaultSpinner.Start("Encrypting...")
					if err != nil {
						pterm.Error.Println(err)
					}
					ac, err := actions.NewCommandFactory().SetEncrypter(c)
					if err != nil {
						spnr.Fail("Encryption failed!")
						return err
					}
					if err := ac.Encrypt(); err != nil {
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
					spnr, err := pterm.DefaultSpinner.Start("Decrypting...")
					if err != nil {
						pterm.Error.Println(err)
					}
					ac, err := actions.NewCommandFactory().SetDecrypter(c)
					if err != nil {
						spnr.Fail("Decryption failed!")
						return err
					}
					if err := ac.Decrypt(); err != nil {
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
	err := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle("crypt", pterm.NewStyle(pterm.FgCyan)),
		pterm.NewLettersFromStringWithStyle("go", pterm.NewStyle(pterm.FgLightMagenta))).
		Render()
	if err != nil {
		pterm.Error.Println(err)
	}
}
