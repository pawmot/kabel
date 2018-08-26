# Kabel

Simplify examining of network traffic into, out of, and between docker containers.

## Name

`Kabel` is a polish word meaning `wire`.

## Usage

Just download a released version and run it.

## Development

This app uses a library called astilectron that allows devs to write electron applications
with backends in arbitrary languages (here it's GO). Because of this the build is a little more
involved than just running `go build`:

During development:
* in one terminal window run `make frontend-dist-watch`. This will start a `ng build --watch` that in turn will watch and 
compile any changes in the frontend part of the app.
* in another terminal window run `make`. This runs the `bundle-dev` target which will copy the `dist` folder
from the angular frontend app and bundle everything.
* go to `output/linux-amd64` folder and run `Kabel`.

Release a "production" artifact:
* run `make bundle` to compile a production distribution of the angular app and bundle it with the go backend.
* go to `output/linux-amd64` where you'll find the `Kabel` executable.f
 