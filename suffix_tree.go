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
    "unicode/utf8"
)

// A node in the suffix tree. It stores the children of this node along
// with the word, if existent, that finishes in this node.
type node struct {
    children map[rune]*node // Edges leaving this node.
    word     string         // Word completed in this node.
    group    int            // Group of this word
}

// A suffix tree used to identify the longest known suffix in a given
// word. Along with each suffix, an identifier is stored. This
// identifier is used to choose which action should be taken in the
// stemming process. 
type suffixTree struct {
    root *node // Root node of suffix tree
}

// Create a new tree node with default values.
func newNode() *node {
    n := new(node)
    n.children = make(map[rune]*node)
    n.word = ""
    n.group = -1
    return n
}

// Create a new suffix tree with the root node.
func newSuffixTree() *suffixTree {
    t := new(suffixTree)
    t.root = newNode()
    return t
}

// Add a new suffix to the tree. The word is inserted in reverse order
// to make it easier to match suffixes. The group value is used to
// identify the category of the suffix and take the necessary actions.
func (st *suffixTree) Add(word string, group int) *suffixTree {
    cnode := st.root
    runes := []rune(word)

    for i := len(runes) - 1; i >= 0; i-- {
        n, ok := cnode.children[runes[i]]
        if ok {
            cnode = n
        } else {
            t := newNode()
            cnode.children[runes[i]] = t
            cnode = t
        }
    }

    cnode.word = word
    cnode.group = group
    return st
}

// Returns true if a given word is already stored in the suffix tree.
func (st *suffixTree) Contains(word string) bool {
    cnode := st.root
    runes := []rune(word)

    for i := len(runes) - 1; i >= 0; i-- {
        n, ok := cnode.children[runes[i]]
        if ok {
            cnode = n
        } else {
            return false
        }
    }

    if cnode.word != "" && cnode.word == word {
        return true
    }

    return false
}

// Returns the longest known suffix that matches the given word. If no
// suffix is found, empty string "" and group id -1 are returned. If a known
// suffix matches the word, it is returned along with its category id.
func (st *suffixTree) LongestSuffix(word string) (string, int) {
    cnode := st.root
    runes := []rune(word)

    currentSuffix := ""
    currentSuffixSize := -1
    currentSuffixGroup := -1

    for i := len(runes) - 1; i >= 0; i-- {
        n, ok := cnode.children[runes[i]]
        if ok {
            cnode = n

            // check if a word finishes in this node
            if cnode.word != "" {
                sz := utf8.RuneCountInString(cnode.word)
                if sz > currentSuffixSize {
                    currentSuffix = cnode.word
                    currentSuffixSize = sz
                    currentSuffixGroup = cnode.group
                }
            }
        } else {
            break
        }
    }

    return currentSuffix, currentSuffixGroup
}
