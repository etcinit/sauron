package eye

// SimpleLogger should be able to handle fatal, panic, and info log messages.
// It is compatible with the standard log package and any other library that
// mimics its interface.
type SimpleLogger interface {
	Fatalln(v ...interface{})
	Panicln(v ...interface{})
	Println(v ...interface{})
}

// TrailOptions are the different options supported by the Trail object.
type TrailOptions struct {
	// Logger to be used by the trail. Messages about file events and error
	// will be sent to this logger. Use ioutil.Discard to ignore output.
	Logger SimpleLogger
}
