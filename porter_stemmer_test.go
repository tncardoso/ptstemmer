// ptstemmer - Portuguese stemmer for Go
// 
// Copyright (c) 2013 - Thiago Cardoso <thiagoncc@gmail.com>
// 
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met: 
// 
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer. 
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution. 
// 
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package ptstemmer

import (
    "testing"
    "io"
    "os"
    "bufio"
    "strings"
)

// Test if vowels with diacritics are being correctly idenfied as
// vowels.
func TestVowels (t *testing.T) {
    ps := NewPorterStemmer()
    vowels := "aeiouáéíóúâêô"
    notVowels := "nlpqrxzcvbnm"

    for _, v := range []rune(vowels) {
        if !ps.isVowel(v) {
            t.Errorf("'%s' should be a vowel\n", v)
        }
    }
    for _, v := range []rune(notVowels) {
        if ps.isVowel(v) {
            t.Errorf("'%s' should not be a vowel\n", v)
        }
    }
}

// TestR checks if the remainder regions R1 and R2 are correctly
// identified.
func TestR (t *testing.T) {
    var cases = []struct {
        word string
        r string
    }{
        {"animadversion", "imadversion"},
        {"imadversion", "adversion"},
        {"sprinkled", "kled"},
        {"eucharist", "harist"},
        {"harist", "ist"},
        {"kled", ""},
        {"beau", ""},
        {"", ""},
        {"beauty", "y"},
        {"y", ""},
        {"beautiful", "iful"},
        {"iful", "ul"},
    }

    ps := NewPorterStemmer()

    for _, c := range cases {
        r := ps.r(c.word)
        if r != c.r {
            t.Errorf("Error finding R. expectd= '%s' actual= '%s'\n", c.r, r)
        }
    }
}

// TestRV checks if the portuguese description of RV is correctly
// implemented.
func TestRV (t *testing.T) {
    // Some of these test cases are in spanish since that the RV
    // algorithm is the same.
    var cases = []struct {
        word string
        r string
    }{
        {"macho", "ho"},
        {"trabajo", "bajo"},
        {"áureo", "eo"},
        {"oliva", "va"},
        {"ôôiii", ""},
    }

    ps := NewPorterStemmer()

    for _, c := range cases {
        r := ps.rv(c.word)
        if r != c.r {
            t.Errorf("Error finding RV. expectd= '%s' actual= '%s'\n", c.r, r)
        }
    }
}

// TestStemmer checks if some words are being correctly stemmed. Most of
// this words fall in specific corner cases.
func TestStemmer (t *testing.T) {
    var stemCases = []struct {
        word string
        stem string
    }{
        { "á", "á" },
        { "ajuda", "ajud" },
        { "ajudá", "ajud" },
        { "ajudado", "ajud" },
        { "ajudou", "ajud" },
        { "abafaram", "abaf" },
        { "abaixa", "abaix" },
        { "abraçada", "abrac" },
        { "adequadamente", "adequ" },
        { "aérea", "aér" },
        { "anatomicamente", "anatom" },
        { "cheira", "cheir" },
        { "ôôiii", "ôôiii" },
    }

    ps := NewPorterStemmer()

    for _, c := range stemCases {
        r := ps.Stem(c.word)
        if r != c.stem {
            t.Errorf("Invalid stem. word= %s expected= %s actual= %s",
            c.word, c.stem, r)
        }
    }
}

// TestFile checks if the stemming is working correctly for the snowball
// test cases. The test file have one test case per line in the
// following format:
//
//      [original_word] [expected_stem]
func TestFile (t *testing.T) {
    ip, err := os.Open("testdata/ptstems.txt")
    if err != nil {
        t.Errorf("Could not open test file: testdata/ptstems.txt")
        return
    }
    defer ip.Close()

    ps := NewPorterStemmer()
    r := bufio.NewReader(ip)
    for {
        l, err := r.ReadString('\n')
        if err != nil {
            if err == io.EOF {
                break
            } else {
                t.Errorf("Error reading file")
            }
        }
        spt := strings.SplitN(strings.Trim(l, "\n"), " ", 2)
        word := strings.Trim(spt[0], " ")
        stem := strings.Trim(spt[1], " ")
        res := ps.Stem(word)
        if res != stem {
            t.Errorf("Invalid stem. word= %s expected= %s actual= %s",
            word, stem, res)
            break
        }
    }
}
