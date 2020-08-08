package transform

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/dreadl0ck/netcap/maltego"
)

func openFolder() {
	var (
		lt              = maltego.ParseLocalArguments(os.Args)
		trx             = &maltego.Transform{}
		openCommandName = os.Getenv("NETCAP_MALTEGO_OPEN_FILE_CMD")
		args            []string
	)

	// if no command has been supplied via environment variable
	// then default to:
	// - open for macOS
	// - gio open for linux
	if openCommandName == "" {
		if runtime.GOOS == platformDarwin {
			openCommandName = defaultOpenCommand
		} else { // linux
			openCommandName, args = makeLinuxCommand(defaultOpenCommandLinux, args)
		}
	}

	path := filepath.Dir(lt.Values["location"])
	log.Println("open path:", path)

	log.Println("command for opening path:", openCommandName)
	args = append(args, path)

	out, err := exec.Command(openCommandName, args...).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		log.Fatal(err)
	}
	log.Println(string(out))

	trx.AddUIMessage("completed!", maltego.UIMessageInform)
	fmt.Println(trx.ReturnOutput())
}
