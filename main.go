package main

import (
	"github.com/asticode/go-astilectron"
	"github.com/asticode/go-astilog"
	"flag"
	"github.com/asticode/go-astilectron-bootstrap"
	"github.com/pkg/errors"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/log"
)

var (
	AppName string
	debug   = flag.Bool("d", false, "enables the debug mode")
)

func main() {
	flag.Parse()
	astilog.FlagInit()
	if err := bootstrap.Run(bootstrap.Options{
		Asset:    Asset,
		AssetDir: AssetDir,
		AstilectronOptions: astilectron.Options{
			AppName: AppName,
		},
		OnWait: func(a *astilectron.Astilectron, w []*astilectron.Window, m *astilectron.Menu, t *astilectron.Tray, tm *astilectron.Menu) error {
			w[0].On(astilectron.EventNameWindowEventBlur, func(e astilectron.Event) (deleteListener bool) {
				astilog.Info("BLURRED :O")
				log.Info("BLURRED :D")
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
				BackgroundColor: astilectron.PtrStr("#333"),
				Center:          astilectron.PtrBool(true),
				Height:          astilectron.PtrInt(600),
				Width:           astilectron.PtrInt(600),
			},
		}},
	}); err != nil {
		astilog.Fatal(errors.Wrap(err, "running bootstrap failed"))
	}

}

func handleMessages(w *astilectron.Window, m bootstrap.MessageIn) (payload interface{}, err error) {
	astilog.Info(m.Name)
	w.OpenDevTools()

	if m.Name == "devtools" {
		w.OpenDevTools()
	}

	return nil, nil
}
