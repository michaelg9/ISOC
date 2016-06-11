package models

// DataIn is the XML struct for the data transfer
// from client to server
type DataIn struct {
	Meta    Meta      `xml:"metadata"`
	Battery []Battery `xml:"battery"`
}

/* IDEA: Refactor this as follows
 * type DataIn struct {
 *   Meta
 *   Device
 * }
 * Pro's: Intercompability with output data
 * Problem: Meta tag redundant?
 */

// Meta is the struct for the Metadata
type Meta struct {
	Device int `xml:"device"` // Device id to identify the sending device
}
