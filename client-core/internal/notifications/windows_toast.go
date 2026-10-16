package notifications

import (
	"bytes"
	"encoding/xml"
	"net/url"
)

const protocolName = "pasteferry-notify"

func actionURI(id, action string) string {
	return protocolName + "://action?" + url.Values{"clip": {id}, "action": {action}}.Encode()
}

// Windows toast tags are at most 16 characters; the full record ID remains in
// every activation URI, so an old toast can never accept a different file.
func toastXML(n Message) string {
	escape := func(s string) string { var b bytes.Buffer; _ = xml.EscapeText(&b, []byte(s)); return b.String() }
	return `<toast activationType="protocol" launch="` + escape(actionURI(n.ID, "open")) + `"><visual><binding template="ToastGeneric"><text>` + escape(n.Title) + `</text><text>` + escape(n.Body) + `</text></binding></visual><actions><action activationType="protocol" content="` + escape(n.ActionLabel) + `" arguments="` + escape(actionURI(n.ID, n.Action)) + `"/></actions><audio silent="true"/></toast>`
}
