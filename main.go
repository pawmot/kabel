package main

import (
	"flag"
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

var (
	AppName string
	debug   = flag.Bool("d", false, "enables the debug mode")
	devtoolsOpen = false
)

func main() {
	flag.Parse()
	astilog.FlagInit()
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName:            AppName,
			AppIconDefaultPath: "resources/cable-512.png",
		},
		OnWait: func(a *astilectron.Astilectron, w []*astilectron.Window, m *astilectron.Menu, t *astilectron.Tray, tm *astilectron.Menu) error {
			if *debug {
				w[0].OpenDevTools()
			}

			// TODO the following handlers don't work - find out how to listen to those events. A PR may be necessary here.
			w[0].On("devtools-opened", func(e astilectron.Event) (deleteListener bool) {
				astilog.Debug("Devtools opened")
				devtoolsOpen = true
				return false
			})

			w[0].On("devtools-closed", func(e astilectron.Event) (deleteListener bool) {
				astilog.Debug("Devtools closed")
				devtoolsOpen = false
				return false
			})

			return nil
		},
		Debug:         *debug,
		RestoreAssets: RestoreAssets,
		Windows: []*bootstrap.Window{{
			Homepage:       "index.html",
			MessageHandler: handleMessages,
			Options: &astilectron.WindowOptions{
				Center:    astilectron.PtrBool(true),
				Height:    astilectron.PtrInt(600),
				Width:     astilectron.PtrInt(600),
				MinHeight: astilectron.PtrInt(400),
				MinWidth:  astilectron.PtrInt(400),
				Title:     &AppName,
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}

}

func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	switch m.Name {
	case "connection_spec":
		astilog.Info(string(m.Payload))
	case "devtools":
		if *debug {
			if devtoolsOpen {
				w.CloseDevTools()
			} else {
				w.OpenDevTools()
			}
		}
	}

	return nil, nil
}
