package main

import (
	"github.com/playwright-community/playwright-go"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	/*err := playwright.Install()
	if err != nil {
		panic(err)
	}*/

	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	defer pw.Stop()
	browser, err := pw.Firefox.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	defer browser.Close()

	context, err := browser.NewContext()
	if err != nil {
		panic(err)
	}

	page, err := context.NewPage()
	if err != nil {
		panic(err)
	}
	_, err = page.Goto("https://cromer.cl")
	if err != nil {
		panic(err)
	}

	var imageUrl []string
	entries, err := page.QuerySelectorAll("img")
	if err != nil {
		log.Fatalf("could not get entries: %v", err)
	}

	for _, entry := range entries {
		att, err := entry.GetAttribute("src")
		if err != nil {
			panic(err)
		}
		imageUrl = append(imageUrl, att)
	}

	for i, url := range imageUrl {
		err = downloadImage(url, strconv.Itoa(i)+".png")
		if err != nil {
			panic(err)
		}
	}
}

func downloadImage(src string, filename string) (err error) {
	out, err := os.Create("downloads/" + filename)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err = out.Close()
		if err != nil {
			return
		}
	}(out)
	resp, err := http.Get(src)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
