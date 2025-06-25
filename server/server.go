package server

import (
	"fmt"
	"log"
	"net/http"

	tg "github.com/navigator-systems/jrx/server/templateGit"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func getTemplatesRepo(repo, branch, filename string) tg.TemplateFile {
	// This function can be used to fetch the templates repository
	err := tg.CloneTemplate(repo, branch)
	if err != nil {
		log.Fatalln("Failed to clone template repository:", err)
	}
	templateFile, err := tg.GetTemplateCtrl(filename)
	if err != nil {
		log.Fatalln("Failed to get template control:", err)
	}
	log.Println("Templates loaded successfully!")
	return templateFile
}

func StartServer() {
	serverConfiguration, err := GetConfig()
	if err != nil {
		log.Fatalf("Failed to load server configuration: %v", err)
	}

	info := getTemplatesRepo(serverConfiguration.TemplateRepo,
		serverConfiguration.TemplateBranch,
		serverConfiguration.TemplateFile)
	fmt.Println("Template Information:")
	fmt.Println(info)
	http.HandleFunc("/", indexHandler)

	port := fmt.Sprintf(":%s", serverConfiguration.Port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Println("Server started successfully!")
}
