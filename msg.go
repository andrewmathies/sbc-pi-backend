package bg

// header is one of: Hello, StartUpdate, Complete, Fail
type Msg struct {
	ID		string	`json:"id"`
	Header	string	`json:"header"`
	Version	string	`json:"version"`
}