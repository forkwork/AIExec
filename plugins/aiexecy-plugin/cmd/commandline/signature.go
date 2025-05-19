package main

import (
	"os"

	"aiexec-plugin/cmd/commandline/signature"
	"github.com/spf13/cobra"
)

var (
	signatureGenerateCommand = &cobra.Command{
		Use:   "generate",
		Short: "Generate a key pair",
		Long:  "Generate a key pair",
		Args:  cobra.ExactArgs(0),
		Run: func(c *cobra.Command, args []string) {
			keyPairName := c.Flag("filename").Value.String()
			if keyPairName == "" {
				keyPairName = "aiexec_plugin_signing_key"
			}
			err := signature.GenerateKeyPair(keyPairName)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	signatureSignCommand = &cobra.Command{
		Use:   "sign [aiexecpkg_path]",
		Short: "Sign a aiexecpkg file",
		Long:  "Sign a aiexecpkg file with the specified private key",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			aiexecpkgPath := args[0]
			privateKeyPath := c.Flag("private_key").Value.String()
			err := signature.Sign(aiexecpkgPath, privateKeyPath)
			if err != nil {
				os.Exit(1)
			}
		},
	}

	signatureVerifyCommand = &cobra.Command{
		Use:   "verify [aiexecpkg_path]",
		Short: "Verify a aiexecpkg file",
		Long:  "Verify a aiexecpkg file with the specified public key. If no public key is provided, the official public key will be used",
		Args:  cobra.ExactArgs(1),
		Run: func(c *cobra.Command, args []string) {
			aiexecpkgPath := args[0]
			publicKeyPath := c.Flag("public_key").Value.String()
			err := signature.Verify(aiexecpkgPath, publicKeyPath)
			if err != nil {
				os.Exit(1)
			}
		},
	}
)

func init() {
	signatureCommand.AddCommand(signatureGenerateCommand)
	signatureCommand.AddCommand(signatureSignCommand)
	signatureCommand.AddCommand(signatureVerifyCommand)

	signatureGenerateCommand.Flags().StringP("filename", "f", "", "filename of the key pair")

	signatureSignCommand.Flags().StringP("private_key", "p", "", "private key file")
	signatureSignCommand.MarkFlagRequired("private_key")

	signatureVerifyCommand.Flags().StringP("public_key", "p", "", "public key file")
}
