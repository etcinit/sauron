package eye

// SimpleLogger should be able to handle error and info log messages.
type SimpleLogger interface {
	Infoln(v ...interface{})
	Errorln(v ...interface{})
}

// TrailOptions are the different options supported by the Trail object.
type TrailOptions struct {
	// Logger to be used by the trail. Messages about file events and error
	// will be sent to this logger. Use ioutil.Discard to ignore output.
	Logger SimpleLogger
}
