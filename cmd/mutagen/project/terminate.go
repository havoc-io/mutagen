package project

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/mutagen-io/mutagen/cmd/mutagen/forward"
	"github.com/mutagen-io/mutagen/cmd/mutagen/sync"
	projectcfg "github.com/mutagen-io/mutagen/pkg/configuration/project"
	"github.com/mutagen-io/mutagen/pkg/filesystem/locking"
	"github.com/mutagen-io/mutagen/pkg/identifier"
	"github.com/mutagen-io/mutagen/pkg/project"
)

func terminateMain(command *cobra.Command, arguments []string) error {
	// Validate arguments.
	if len(arguments) > 0 {
		return errors.New("unexpected arguments provided")
	}

	// Compute the name of the configuration file and ensure that our working
	// directory is that in which the file resides. This is required for
	// relative paths (including relative synchronization paths and relative
	// Unix Domain Socket paths) to be resolved relative to the project
	// configuration file.
	configurationFileName := project.ConfigurationFileName()
	if terminateConfiguration.projectFile != "" {
		var directory string
		directory, configurationFileName = filepath.Split(terminateConfiguration.projectFile)
		if directory != "" {
			if err := os.Chdir(directory); err != nil {
				return errors.Wrap(err, "unable to switch to target directory")
			}
		}
	}

	// Compute the lock path.
	lockPath := configurationFileName + project.LockFileExtension

	// Track whether or not we should remove the lock file on return.
	var removeLockFileOnReturn bool

	// Create a locker and defer its closure and potential removal. On Windows
	// systems, we have to handle this removal after the file is closed.
	locker, err := locking.NewLocker(lockPath, 0600)
	if err != nil {
		return errors.Wrap(err, "unable to create project locker")
	}
	defer func() {
		locker.Close()
		if removeLockFileOnReturn && runtime.GOOS == "windows" {
			os.Remove(lockPath)
		}
	}()

	// Acquire the project lock and defer its release and potential removal. On
	// Windows systems, we can't remove the lock file if it's locked or even
	// just opened, so we handle removal for Windows systems after we close the
	// lock file (see above). In this case, we truncate the lock file before
	// releasing it to ensure that any other process that opens or acquires the
	// lock file before we manage to remove it will simply see an empty lock
	// file, which it will ignore or attempt to remove.
	if err := locker.Lock(true); err != nil {
		return errors.Wrap(err, "unable to acquire project lock")
	}
	defer func() {
		if removeLockFileOnReturn {
			if runtime.GOOS == "windows" {
				locker.Truncate(0)
			} else {
				os.Remove(lockPath)
			}
		}
		locker.Unlock()
	}()

	// Read the project identifier from the lock file. If the lock file is
	// empty, then we can assume that we created it when we created the lock and
	// just remove it.
	buffer := &bytes.Buffer{}
	if length, err := buffer.ReadFrom(locker); err != nil {
		return errors.Wrap(err, "unable to read project lock")
	} else if length == 0 {
		removeLockFileOnReturn = true
		return errors.New("project not running")
	}
	projectIdentifier := buffer.String()

	// Ensure that the project identifier is valid.
	if !identifier.IsValid(projectIdentifier) {
		return errors.New("invalid project identifier found in project lock")
	}

	// Load the configuration file.
	configuration, err := projectcfg.LoadConfiguration(configurationFileName)
	if err != nil {
		return errors.Wrap(err, "unable to load configuration file")
	}

	// Perform pre-termination commands.
	for _, command := range configuration.BeforeTerminate {
		fmt.Println(">", command)
		if err := runInShell(command); err != nil {
			return errors.Wrap(err, "pre-terminate command failed")
		}
	}

	// Compute the label selector that we're going to use to terminate sessions.
	labelSelector := fmt.Sprintf("%s=%s", project.LabelKey, projectIdentifier)

	// Terminate forwarding sessions.
	if err := forward.TerminateWithLabelSelector(labelSelector); err != nil {
		return errors.Wrap(err, "unable to terminate forwarding session(s)")
	}

	// Terminate synchronization sessions.
	if err := sync.TerminateWithLabelSelector(labelSelector); err != nil {
		return errors.Wrap(err, "unable to terminate synchronization session(s)")
	}

	// Perform post-termination commands.
	for _, command := range configuration.AfterTerminate {
		fmt.Println(">", command)
		if err := runInShell(command); err != nil {
			return errors.Wrap(err, "post-terminate command failed")
		}
	}

	// Schedule the project lock for removal.
	removeLockFileOnReturn = true

	// Success.
	return nil
}

var terminateCommand = &cobra.Command{
	Use:          "terminate",
	Short:        "Terminate project sessions",
	RunE:         terminateMain,
	SilenceUsage: true,
}

var terminateConfiguration struct {
	// help indicates whether or not to show help information and exit.
	help bool
	// projectFile is the path to the project file, if non-default.
	projectFile string
}

func init() {
	// Grab a handle for the command line flags.
	flags := terminateCommand.Flags()

	// Disable alphabetical sorting of flags in help output.
	flags.SortFlags = false

	// Manually add a help flag to override the default message. Cobra will
	// still implement its logic automatically.
	flags.BoolVarP(&terminateConfiguration.help, "help", "h", false, "Show help information")

	// Wire up project file flags.
	flags.StringVarP(&terminateConfiguration.projectFile, "project-file", "f", "", "Specify project file")
}
