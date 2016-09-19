# godocdoc

`godocdoc` is a dumb wrapper around Godoc that does two things:

1) Starts godoc
2) Waits for godoc to listen on a port, then opens the doc URL in a web browser

If the working directory is in the GOPATH, `godocdoc` will
open the browser to the location of the current package, i.e.
`http://localhost:6060/pkg/github.com/kevinburke/godoc`.

### Configuration

You can change the port by passing `--port=<port>` to `godocdoc`. The `godoc`
binary has a number of formatting options, which are all currently ignored.
