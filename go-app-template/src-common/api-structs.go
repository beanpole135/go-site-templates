package main
// =================
// Sample "common" file which both the server
// and WASM builds can use. 
// These structs correspond to the "body" of the API requests/replies
// =================
// NOTE: Make sure to include the "json:" field-format tags
// =================
import ()

//Example Request: Requesting a "list" of info from an API
type APIRequest struct {
	ListName string	`json:"list_name"`
}

//Example Reply: Parse out the "list" of data returned from the API
type APIReply struct {
	List []string	`json:"list"`
}
