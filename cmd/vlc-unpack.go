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

var vlcUnpackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Pack file using variable-length code",
	Run:   unpack,
}

const vlcUnpackedExtension = ".txt"

func unpack(_ *cobra.Command, args []string) {
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

	packed := vlc.Decode(data)
	if err := os.WriteFile(unpackedFileName(filePath), []byte(packed), 0644); err != nil {
		handleError(err)
	}
}

func unpackedFileName(path string) string {
	fileName := filepath.Base(path)

	return strings.TrimSuffix(fileName, filepath.Ext(fileName)) + vlcUnpackedExtension
}

func init() {
	rootCmd.AddCommand(vlcUnpackCmd)
}
