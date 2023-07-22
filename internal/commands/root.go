package commands

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type RootCommand struct {
	logger *zap.Logger
	root   *cobra.Command
	input  string
}

func NewRootCommand() (*RootCommand, error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.DisableStacktrace = true
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, errors.Wrap(err, "Error initializing logger.")
	}
	command := &RootCommand{
		logger: logger,
	}
	command.root = &cobra.Command{
		Use:           "ipman",
		Short:         "An opinionated IP address management tool, providing integration with DHCP and DNS servers.",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
	command.root.PersistentFlags().StringVar(&command.input, "input", ".", "The path to the input files.")
	registerCheckCommand(command)
	registerDiffCommand(command)
	return command, nil
}

func (command *RootCommand) Run() {
	err := command.root.Execute()
	if err != nil {
		command.logger.Sugar().Fatalf("%+v", err)
	}
}
