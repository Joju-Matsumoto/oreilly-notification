/*
Copyright © 2024 jj_ <jj_@team-jj.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Joju-Matsumoto/oreilly-notification/internal/adapter/discordnotifier"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/adapter/jsonrepository"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/adapter/oreilly"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/adapter/printnotifier"
	"github.com/Joju-Matsumoto/oreilly-notification/internal/domain/notifier"
	notifyrecentlyaddedbooks "github.com/Joju-Matsumoto/oreilly-notification/internal/usecase/notify_recently_added_books"
	updaterepository "github.com/Joju-Matsumoto/oreilly-notification/internal/usecase/update_repository"
	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, args []string) error {
	// positional arguments

	jsonPath := args[0]

	// flag arguments

	token, _ := cmd.Flags().GetString("token")
	channelID, _ := cmd.Flags().GetString("channel")
	notify, _ := cmd.Flags().GetBool("notify")

	// init updateRepository Usecase

	bookWebAPI := oreilly.New()

	bookRepository := jsonrepository.New(jsonPath)
	if err := bookRepository.Open(); err != nil {
		return fmt.Errorf("error Open: %w", err)
	}
	defer func() {
		if err := bookRepository.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	updateRepositoryUsecase := updaterepository.NewUsecase(bookWebAPI, bookRepository)

	// init notifyRecentlyAddedBooks Usecases

	var bookNotifier notifier.BookNotifier
	if !notify {
		bookNotifier = printnotifier.New(os.Stdout)
	} else {
		discordBookNotifier, err := discordnotifier.New(discordnotifier.Config{
			Token: token,
			TargetChannelIDs: []string{
				channelID,
			},
		})
		if err != nil {
			return err
		}
		if err := discordBookNotifier.Open(); err != nil {
			return err
		}
		defer discordBookNotifier.Close()

		bookNotifier = discordBookNotifier
	}
	notifyRecentlyAddedBooksUsecase := notifyrecentlyaddedbooks.NewUsecase(updateRepositoryUsecase, bookNotifier)

	// execute usecase
	if err := notifyRecentlyAddedBooksUsecase.NotifyRecentlyAddedBooks(context.Background()); err != nil {
		return err
	}
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oreilly-notification",
	Short: "oreilly learning recently added books notifier",
	Long:  `oreilly learning recently added books notifier`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(cmd, args); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.oreilly-notification.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("token", "t", "", "discord bot token")
	rootCmd.Flags().StringP("channel", "c", "", "discord channel id")
	rootCmd.Flags().BoolP("notify", "n", false, "execute discord notification")
}
