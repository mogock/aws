package response

import (
	"encoding/json"
	xml2 "encoding/xml"
	"net/http"
)

var GLOBAL = "Geraldo Mogock"

type Profile struct {
	Name string
	Hobbies []string
}

//Response JSON
func FooJSON (w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Geraldo", []string{"Jugar", "Mirar Serie", "Gimnasio"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//Response XML
func FooXML (w http.ResponseWriter, request *http.Request) {
	profile := Profile{Name: "Jenny", Hobbies: []string{"Hablar, Dream, Trabajar"}}
	xml, err := xml2.MarshalIndent(profile, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/xml")
	w.Write(xml)
}