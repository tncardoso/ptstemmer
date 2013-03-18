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
)

// Checks if fluent pattern is working correctly for the Add function.
func TestFluent(t *testing.T) {
    st := newSuffixTree()

    st.Add("horse", 0).Add("banana", 1).Add("dog", 2)

    if !st.Contains("horse") {
        t.Errorf("Missing word: horse\n")
    }
    if !st.Contains("banana") {
        t.Errorf("Missing word: horse\n")
    }
    if !st.Contains("dog") {
        t.Errorf("Missing word: dog\n")
    }
}

// Checks if words are being correctly inserted and retrieved from the
// suffix tree.
func TestContains(t *testing.T) {
    addedWords := []string{
        "horse",
        "banana",
        "ana",
        "ban",
        "dog"}
    notAddedWords := []string{
        "hor",
        "bana",
        "hors",
        "do"}

    st := newSuffixTree()
    for _, w := range addedWords {
        st.Add(w, 0)
    }

    // Check for words that should be present.
    for _, w := range addedWords {
        if !st.Contains(w) {
            t.Errorf("Missing word: %s\n", w)
        }
    }

    // Check for words that were not added.
    for _, w := range notAddedWords {
        if st.Contains(w) {
            t.Errorf("False positive word: %s\n", w)
        }
    }
}

// Checks if the longest suffix is being correctly retrieved from the
// suffix tree.
func TestLongestSuffix(t *testing.T) {
    addedWords := []string{
        "ismos",
        "a",
        "ma",
        "dog",
        "ia"}

    var cases = []struct {
        word   string
        suffix string
        group  int
    }{
        {"algoritmos", "", -1},
        {"algorismos", "ismos", 1},
        {"laia", "ia", 1},
        {"lama", "ma", 1},
        {"dog", "dog", 1},
        {"abaixa", "a", 1},
    }

    st := newSuffixTree()
    for _, w := range addedWords {
        st.Add(w, 1)
    }

    for _, w := range addedWords {
        if !st.Contains(w) {
            t.Errorf("Word should be in tree: %s\n", w)
        }
    }

    for _, c := range cases {
        r, g := st.LongestSuffix(c.word)
        if r != c.suffix {
            t.Errorf("Wrong suffix. word= %s expected= %s returned= %s\n",
                c.word, c.suffix, r)
        }
        if g != c.group {
            t.Errorf("Wrong group. expected= %d returned= %d\n",
                c.group, g)
        }
    }
}
