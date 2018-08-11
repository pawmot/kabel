let index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        asticode.loader.show();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            document.getElementById('hw').innerText = 'READY!';

            asticode.loader.hide();

            astilectron.sendMessage({"name": "My Message", "paylod": "Wut!!!"});

            document.onkeypress(ev => {

                if (ev.key === "F12") {
                    astilectron.sendMessage({"name": "devtools", "paylod": "Wut!!!"});
                }
            })
        })
    },
};