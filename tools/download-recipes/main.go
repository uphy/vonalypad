package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/gommon/log"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

type recipe struct {
	title string
}

func run() error {
	urls, err := readURLs()
	if err != nil {
		return err
	}
	for _, u := range urls {
		dir := u[1:]
		log.Printf("Download recipe: url=%s", u)
		os.MkdirAll(dir, 0777)
		_, err := downloadRecipe("https://cookpad.com/"+u, dir)
		if err != nil {
			log.Printf("Failed to download: url=%s, err=%v", u, err)
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func downloadRecipe(u string, dir string) (*recipe, error) {
	resp, err := http.DefaultClient.Get(u)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "index.html"), b, 0666); err != nil {
		return nil, err
	}
	return nil, nil
}

func readURLs() ([]string, error) {
	b, err := ioutil.ReadFile("dendoiri.json")
	if err != nil {
		return nil, err
	}
	var s []string
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}
