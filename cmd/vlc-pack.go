package cmd

import (
	"archivator/lib/vlc"
	"errors"
	"github.com/spf13/cobra"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var vlcPackCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack file using variable-length code",
	Run:   pack,
}

const vlcPackedExtension = ".yu"

var ErrEmptyPath = "please specify a file"

func pack(_ *cobra.Command, args []string) {
	if len(args) < 1 {
		handleError(errors.New(ErrEmptyPath))
	}

	filePath := args[0]

	r, err := os.Open(filePath)
	if err != nil {
		handleError(err)
	}
	defer func() {
		err := r.Close()
		if err != nil {
			handleError(err)
		}
	}()

	data, err := io.ReadAll(r)
	if err != nil {
		handleError(err)
	}

	packed := vlc.Encode(string(data))
	if err := os.WriteFile(packedFileName(filePath), packed, 0644); err != nil {
		handleError(err)
	}
}

func packedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + vlcPackedExtension
}

func init() {
	rootCmd.AddCommand(vlcPackCmd)
}
