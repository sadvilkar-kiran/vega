package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/sadvilkar-kiran/vega"
)

func doNew(appName string) {
	appName = strings.ToLower(appName)
	appURL = appName

	// sanitize the application name (convert url to single word)
	if strings.Contains(appName, "/") {
		exploded := strings.SplitAfter(appName, "/")
		appName = exploded[(len(exploded) - 1)]
	}

	log.Println("App name is", appName)

	// create directory structure
	color.Green("\tCreating application structure...")
	err := os.MkdirAll("./"+appName, 0755)
	if err != nil {
		exitGracefully(err)
	}

	// create subdirectories
	dirs := []string{"handlers", "migrations", "views", "views/layouts", "mail", "data", "public", "tmp", "logs", "middleware"}
	for _, dir := range dirs {
		err = os.MkdirAll(fmt.Sprintf("./%s/%s", appName, dir), 0755)
		if err != nil {
			exitGracefully(err)
		}
	}

	// create .env file
	color.Yellow("\tCreating .env file...")
	data, err := templateFS.ReadFile("templates/env.txt")
	if err != nil {
		exitGracefully(err)
	}

	env := string(data)
	env = strings.ReplaceAll(env, "${APP_NAME}", appName)
	// Create a temporary Vega instance to generate random string
	tempVega := &vega.Vega{}
	env = strings.ReplaceAll(env, "${KEY}", tempVega.RandomString(32))

	err = copyDataToFile([]byte(env), fmt.Sprintf("./%s/.env", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create go.mod file
	color.Yellow("\tCreating go.mod file...")
	data, err = templateFS.ReadFile("templates/go.mod.txt")
	if err != nil {
		exitGracefully(err)
	}

	mod := string(data)
	mod = strings.ReplaceAll(mod, "${APP_NAME}", appURL)

	err = copyDataToFile([]byte(mod), fmt.Sprintf("./%s/go.mod", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create main.go
	color.Yellow("\tCreating source files...")
	data, err = templateFS.ReadFile("templates/main.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	mainContent := strings.ReplaceAll(string(data), "${APP_NAME}", appURL)
	err = copyDataToFile([]byte(mainContent), fmt.Sprintf("./%s/main.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create init-vega.go
	data, err = templateFS.ReadFile("templates/init-vega.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	initContent := strings.ReplaceAll(string(data), "${APP_NAME}", appURL)
	err = copyDataToFile([]byte(initContent), fmt.Sprintf("./%s/init-vega.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create routes.go
	data, err = templateFS.ReadFile("templates/routes.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	err = copyDataToFile([]byte(string(data)), fmt.Sprintf("./%s/routes.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create convenience.go
	data, err = templateFS.ReadFile("templates/convenience.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	err = copyDataToFile([]byte(string(data)), fmt.Sprintf("./%s/convenience.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create handlers/handlers.go
	data, err = templateFS.ReadFile("templates/handlers/handlers.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	handlerContent := strings.ReplaceAll(string(data), "${APP_NAME}", appURL)
	err = copyDataToFile([]byte(handlerContent), fmt.Sprintf("./%s/handlers/handlers.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create data/models.go
	data, err = templateFS.ReadFile("templates/data/models.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	modelContent := strings.ReplaceAll(string(data), "${APP_NAME}", appURL)
	err = copyDataToFile([]byte(modelContent), fmt.Sprintf("./%s/data/models.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// create middleware/middleware.go
	data, err = templateFS.ReadFile("templates/middleware/middleware.go.txt")
	if err != nil {
		exitGracefully(err)
	}
	middlewareContent := strings.ReplaceAll(string(data), "${APP_NAME}", appURL)
	err = copyDataToFile([]byte(middlewareContent), fmt.Sprintf("./%s/middleware/middleware.go", appName))
	if err != nil {
		exitGracefully(err)
	}

	// copy view templates
	color.Yellow("\tCopying view templates...")
	viewFiles := []string{
		"views/home.jet",
		"views/home.page.tmpl",
		"views/login.jet",
		"views/forgot.jet",
		"views/reset-password.jet",
		"views/form.jet",
		"views/cache.jet",
		"views/sessions.jet",
		"views/jet-template.jet",
		"views/layouts/base.jet",
	}

	for _, viewFile := range viewFiles {
		data, err = templateFS.ReadFile("templates/" + viewFile)
		if err != nil {
			// skip if file doesn't exist
			continue
		}
		targetPath := fmt.Sprintf("./%s/%s", appName, viewFile)
		err = os.MkdirAll(filepath.Dir(targetPath), 0755)
		if err != nil {
			exitGracefully(err)
		}
		err = copyDataToFile(data, targetPath)
		if err != nil {
			exitGracefully(err)
		}
	}

	// update source files with correct app name
	color.Yellow("\tUpdating source files...")
	os.Chdir("./" + appName)
	updateSource()

	// run go mod tidy in the project directory
	color.Yellow("\tRunning go mod tidy...")
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	if err != nil {
		color.Yellow("Warning: go mod tidy failed. You may need to run it manually.")
	}

	color.Green("Done building " + appURL)
	color.Green("Go build something awesome!")
	color.Yellow("\nNext steps:")
	color.Yellow("  cd %s", appName)
	color.Yellow("  go run .")
}

