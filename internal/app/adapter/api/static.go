package api

import (
	"fmt"
	"net/http"
)

func (r *Router) static() {
	mainPath := "/"
	r.mux.Handle(mainPath,
		http.StripPrefix(mainPath,
			http.FileServer(http.Dir("./static/")),
		),
	)

	docPath := fmt.Sprintf("%sdoc/", r.path)
	r.mux.Handle(docPath,
		http.StripPrefix(docPath,
			http.FileServer(http.Dir("./docs/static/")),
		),
	)
}