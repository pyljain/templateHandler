package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

func main() {
	componentName := os.Args[1]
	componentDir := fmt.Sprintf("./%s", componentName)

	// Clone Repository
	repoName := os.Args[2]
	cloneRepository(componentDir, repoName)

	// Read Template.yaml
	params := getQuestions(componentDir)

	// Apply variables
	apply(componentDir, params)
}

func cloneRepository(dir string, name string) {
	_, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL: name,
	})

	if err != nil {
		log.Fatal(err)
	}
}

func getQuestions(dir string) map[string]string {
	templateFileLocation := fmt.Sprintf("%s/template.yaml", dir)
	data, err := os.ReadFile(templateFileLocation)
	if err != nil {
		log.Fatalf("Could not read template file %s", err)
	}

	template := templateFormat{}
	err = yaml.Unmarshal(data, &template)
	if err != nil {
		log.Fatalf("Could not read template file %s", err)
	}

	result := make(map[string]string)
	for _, parameter := range template.Parameters {
		fmt.Println(parameter.Question)
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)
		result[parameter.Value] = text
	}

	return result
}

type templateFormat struct {
	Parameters []templateParameters `yaml:"parameters"`
}

type templateParameters struct {
	Question string `yaml:"question"`
	Value    string `yaml:"value"`
}

func apply(rootDir string, params map[string]string) {
	err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, ".git") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		fileContents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		tmp := template.New("simple")
		tmp, err = tmp.Parse(string(fileContents))
		if err != nil {
			return err
		}

		f, err := os.Create(path)
		defer f.Close()

		err = tmp.Execute(f, params)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}
