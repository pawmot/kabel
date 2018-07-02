package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"flag"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/pkg/errors"
)

var (
	AppName string
	debug = flag.Bool("d", false, "enables the debug mode")
)

func main() {
	flag.Parse()
	astilog.FlagInit()
	if err := bootstrap.Run(bootstrap.Options{
		Asset: Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName: AppName,
		},
		Debug: *debug,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage: "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				BackgroundColor: astilectron.PtrStr("#333"),
				Center: astilectron.PtrBool(true),
				Height: astilectron.PtrInt(600),
				Width:  astilectron.PtrInt(600),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}
}

func handleMessages(_ *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	return nil, nil
}