package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type Resultado struct {
	URL      string
	Estado   int
	Duracion time.Duration
	Err      error
}

func verificar(url string) Resultado {
	inicio := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return Resultado{URL: url, Duracion: time.Since(inicio), Err: err}
	}
	defer resp.Body.Close()
	return Resultado{URL: url, Estado: resp.StatusCode, Duracion: time.Since(inicio)}
}

func leerURLs(ruta string) ([]string, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var urls []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		if linea := strings.TrimSpace(sc.Text()); linea != "" {
			urls = append(urls, linea)
		}
	}
	return urls, sc.Err()
}

func main() {
	urls, err := leerURLs(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	inicio := time.Now()
	for _, url := range urls {
		r := verificar(url)
		fmt.Printf("%d  %v  %s\n", r.Estado, r.Duracion.Round(time.Millisecond), r.URL)
	}
	fmt.Printf("\n%d sitios en %v\n", len(urls), time.Since(inicio).Round(time.Millisecond))
}
