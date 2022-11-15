package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ddworken/hishtory/client/hctx"
	"github.com/ddworken/hishtory/client/lib"
	"github.com/spf13/cobra"
)

var offline *bool

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Re-initialize hiSHtory with a specified secret key",
	GroupID: GROUP_ID_CONFIG,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Flag("offline").Value.String()
		db, err := hctx.OpenLocalSqliteDb()
		lib.CheckFatalError(err)
		data, err := lib.Search(nil, db, "", 10)
		lib.CheckFatalError(err)
		if len(data) > 0 {
			fmt.Printf("Your current hishtory profile has saved history entries, are you sure you want to run `init` and reset?\nNote: This won't clear any imported history entries from your existing shell\n[y/N]")
			reader := bufio.NewReader(os.Stdin)
			resp, err := reader.ReadString('\n')
			lib.CheckFatalError(err)
			if strings.TrimSpace(resp) != "y" {
				fmt.Printf("Aborting init per user response of %#v\n", strings.TrimSpace(resp))
				return
			}
		}
		secretKey := ""
		if len(args) > 0 {
			secretKey = args[0]
		}
		lib.CheckFatalError(lib.Setup(secretKey, *offline))
		if os.Getenv("HISHTORY_SKIP_INIT_IMPORT") == "" {
			fmt.Println("Importing existing shell history...")
			ctx := hctx.MakeContext()
			numImported, err := lib.ImportHistory(ctx, false, false)
			lib.CheckFatalError(err)
			if numImported > 0 {
				fmt.Printf("Imported %v history entries from your existing shell history\n", numImported)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	offline = initCmd.Flags().Bool("offline", false, "Install hiSHtory in offline mode wiht all syncing capabilities disabled")
}
