package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Asuforce/go-chi-api/internal/server"
	"github.com/go-chi/docgen"
)

func main() {
	var (
		port   = flag.String("port", "8080", "Addr to bind")
		env    = flag.String("env", "develop", "Exec env (local, beta, production")
		gendoc = flag.Bool("gendoc", true, "Generate document")
	)
	flag.Parse()

	s := server.New()
	s.Init(*env)
	s.Middleware()
	s.Router()

	if *gendoc {
		doc := docgen.MarkdownRoutesDoc(s.GetRouter(), docgen.MarkdownOpts{
			ProjectPath: "github.com/Asuforce/go-chi-api",
			Intro:       "generated docs.",
		})
		file, err := os.Create(`doc.md`)
		if err != nil {
			log.Printf("err: %v", err)
		}
		defer file.Close()
		file.Write(([]byte)(doc))
	}
	log.Println("starting app")
	http.ListenAndServe(fmt.Sprint(":", *port), s.GetRouter())
}
