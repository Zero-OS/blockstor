package gridapi

import (
	"log"
	"net/http"

	"github.com/g8os/blockstor/nbdserver/stubs/gridapi/goraml"

	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
)

func main() {
	// input validator
	validator.SetValidationFunc("multipleOf", goraml.MultipleOf)

	r := mux.NewRouter()

	// home page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// apidocs
	r.PathPrefix("/apidocs/").Handler(http.StripPrefix("/apidocs/", http.FileServer(http.Dir("./apidocs/"))))

	NodesInterfaceRoutes(r, NodesAPI{})

	StorageclustersInterfaceRoutes(r, StorageclustersAPI{})

	VolumesInterfaceRoutes(r, VolumesAPI{})

	log.Println("starting server")
	http.ListenAndServe(":5000", r)
}