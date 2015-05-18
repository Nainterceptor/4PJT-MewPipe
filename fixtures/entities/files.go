package entities

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func download(link string) (string, error) {

	fmt.Println("Downloading file...")

	file, err := ioutil.TempFile("", "fixtures_")

	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	check := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	resp, err := check.Get(link) // add a filter to check redirect

	if err != nil {
		fmt.Println(err)
		return "", err

	}
	defer resp.Body.Close()
	fmt.Println(resp.Status)

	_, err = io.Copy(file, resp.Body)

	if err != nil {
		fmt.Println(err)
		return "", err

	}

	thePath, err := filepath.Abs(filepath.Dir(file.Name()))

	if err != nil {
		fmt.Println(err)
		return "", err

	}
	return thePath, nil
}
