package omap

// This interface is a soap box.. I would like more packages in go to give programatic way of denoting
// if a instance is thread safe.  If people do, that is great!
type IsThreadSafe interface {
	// Should return true if the instance is thread safe, fakse if not.
	ThreadSafe() bool
}
