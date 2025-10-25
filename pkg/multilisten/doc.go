// Package multilisten provides a facility to bundle multiple net.Listener
// instances into a single one. This is useful if you have an external package
// which expects a single listener but you want to listen on several ports at
// once, or on a specific set of interfaces.
//
// Use the Bundle function to bundle up multiple listeners.
package multilisten
