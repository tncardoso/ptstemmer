Portuguese stemmer for Go
=======================

This package contains a Porter stemmer implementation for Portuguese.

Installing
----------

To install `ptstemmer` run:

    go get github.com/tncardoso/ptstemmer

Usage Example
--------------

    package main

    import (
        "fmt"
        "github.com/tncardoso/ptstemmer"
    )

    func main() {
        // Porter stemmer implements Stemmer interface
        var stemmer ptstemmer.Stemmer = ptstemmer.NewPorterStemmer()

        words := []string {
            "ajudar",
            "ajudei",
            "ajudou",
            "ajuda",
        }   

        for _, w := range words {
            fmt.Printf("word= %s stem= %s\n", w, stemmer.Stem(w))
        }   
    }

