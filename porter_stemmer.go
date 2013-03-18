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

import "strings"

// PorterStemmer implements the Porter stemming algorithm for the
// portuguese language.
// The implementation was based in the following implementation:
// http://snowball.tartarus.org/algorithms/portuguese/stemmer.html
type PorterStemmer struct {
    vowels  map[rune]bool       // Runes that should be considered vowels
    step1SuffixTree *suffixTree // Suffixes checked in step1
    step2SuffixTree *suffixTree // Suffixes checked in step2
    step4SuffixTree *suffixTree // Suffixes checked in step4
    step5SuffixTree *suffixTree // Suffixes checked in step5
}

// Create Porter stemmer struct. Vowels and necessary suffixes for the
// algorithm are also loaded in this step.
func NewPorterStemmer () (*PorterStemmer) {
    ps := new(PorterStemmer)

    // Load portuguese vowels.
    ps.vowels = make(map[rune]bool)
    vowels := "aeiouáéíóúâêô"
    vowelsRunes := []rune(vowels)
    for _, rn := range vowelsRunes {
        ps.vowels[rn] = true
    }

    // Load suffixes that are checked in Step 1.
    ps.step1SuffixTree = newSuffixTree()
    ps.step1SuffixTree.Add("eza",       0).Add("ezas",      0)
    ps.step1SuffixTree.Add("ico",       0).Add("ica",       0)
    ps.step1SuffixTree.Add("icos",      0).Add("icas",      0)
    ps.step1SuffixTree.Add("ismo",      0).Add("ismos",     0)
    ps.step1SuffixTree.Add("ável",      0).Add("ível",      0)
    ps.step1SuffixTree.Add("ista",      0).Add("istas",     0)
    ps.step1SuffixTree.Add("oso",       0).Add("osa",       0)
    ps.step1SuffixTree.Add("osos",      0).Add("osas",      0)
    ps.step1SuffixTree.Add("amento",    0).Add("amentos",   0)
    ps.step1SuffixTree.Add("imento",    0).Add("imentos",   0)
    ps.step1SuffixTree.Add("adora",     0).Add("ador",      0)
    ps.step1SuffixTree.Add("aça~o",     0).Add("adoras",    0)
    ps.step1SuffixTree.Add("adores",    0).Add("aço~es",    0)
    ps.step1SuffixTree.Add("ante",      0).Add("antes",     0)
    ps.step1SuffixTree.Add("ância",     0)
    ps.step1SuffixTree.Add("logía",     1).Add("logías",    1)
    ps.step1SuffixTree.Add("ución",     2).Add("uciones",   2)
    ps.step1SuffixTree.Add("ência",     3).Add("ências",    3)
    ps.step1SuffixTree.Add("amente",    4)
    ps.step1SuffixTree.Add("mente",     5)
    ps.step1SuffixTree.Add("idade",     6).Add("idades",    6)
    ps.step1SuffixTree.Add("iva",       7).Add("ivo",       7)
    ps.step1SuffixTree.Add("ivas",      7).Add("ivos",      7)
    ps.step1SuffixTree.Add("ira",       8).Add("iras",      8)

    // Load suffixes that are checked in Step 2.
    ps.step2SuffixTree = newSuffixTree()
    ps.step2SuffixTree.Add("ada",       0).Add("ida",       0)
    ps.step2SuffixTree.Add("ia",        0).Add("aria",      0)
    ps.step2SuffixTree.Add("eria",      0).Add("iria",      0)
    ps.step2SuffixTree.Add("ará",       0).Add("ara",       0)
    ps.step2SuffixTree.Add("erá",       0).Add("era",       0)
    ps.step2SuffixTree.Add("irá",       0).Add("ava",       0)
    ps.step2SuffixTree.Add("asse",      0).Add("esse",      0)
    ps.step2SuffixTree.Add("isse",      0).Add("aste",      0)
    ps.step2SuffixTree.Add("este",      0).Add("iste",      0)
    ps.step2SuffixTree.Add("ei",        0).Add("arei",      0)
    ps.step2SuffixTree.Add("erei",      0).Add("irei",      0)
    ps.step2SuffixTree.Add("am",        0).Add("iam",       0)
    ps.step2SuffixTree.Add("ariam",     0).Add("eriam",     0)
    ps.step2SuffixTree.Add("iriam",     0).Add("aram",      0)
    ps.step2SuffixTree.Add("eram",      0).Add("iram",      0)
    ps.step2SuffixTree.Add("avam",      0).Add("em",        0)
    ps.step2SuffixTree.Add("arem",      0).Add("erem",      0)
    ps.step2SuffixTree.Add("irem",      0).Add("assem",     0)
    ps.step2SuffixTree.Add("essem",     0).Add("issem",     0)
    ps.step2SuffixTree.Add("ado",       0).Add("ido",       0)
    ps.step2SuffixTree.Add("ando",      0).Add("endo",      0)
    ps.step2SuffixTree.Add("indo",      0).Add("ara~o",     0)
    ps.step2SuffixTree.Add("era~o",     0).Add("ira~o",     0)
    ps.step2SuffixTree.Add("ar",        0).Add("er",        0)
    ps.step2SuffixTree.Add("ir",        0).Add("as",        0)
    ps.step2SuffixTree.Add("adas",      0).Add("idas",      0)
    ps.step2SuffixTree.Add("ias",       0).Add("arias",     0)
    ps.step2SuffixTree.Add("erias",     0).Add("irias",     0)
    ps.step2SuffixTree.Add("arás",      0).Add("aras",      0)
    ps.step2SuffixTree.Add("erás",      0).Add("eras",      0)
    ps.step2SuffixTree.Add("irás",      0).Add("avas",      0)
    ps.step2SuffixTree.Add("es",        0).Add("ardes",     0)
    ps.step2SuffixTree.Add("erdes",     0).Add("irdes",     0)
    ps.step2SuffixTree.Add("ares",      0).Add("eres",      0)
    ps.step2SuffixTree.Add("ires",      0).Add("asses",     0)
    ps.step2SuffixTree.Add("esses",     0).Add("isses",     0)
    ps.step2SuffixTree.Add("astes",     0).Add("estes",     0)
    ps.step2SuffixTree.Add("istes",     0).Add("is",        0)
    ps.step2SuffixTree.Add("ais",       0).Add("eis",       0)
    ps.step2SuffixTree.Add("íeis",      0).Add("aríeis",    0)
    ps.step2SuffixTree.Add("eríeis",    0).Add("iríeis",    0)
    ps.step2SuffixTree.Add("áreis",     0).Add("areis",     0)
    ps.step2SuffixTree.Add("éreis",     0).Add("ereis",     0)
    ps.step2SuffixTree.Add("íreis",     0).Add("ireis",     0)
    ps.step2SuffixTree.Add("ásseis",    0).Add("ésseis",    0)
    ps.step2SuffixTree.Add("ísseis",    0).Add("áveis",     0)
    ps.step2SuffixTree.Add("ados",      0).Add("idos",      0)
    ps.step2SuffixTree.Add("ámos",      0).Add("amos",      0)
    ps.step2SuffixTree.Add("íamos",     0).Add("aríamos",   0)
    ps.step2SuffixTree.Add("eríamos",   0).Add("iríamos",   0)
    ps.step2SuffixTree.Add("áramos",    0).Add("éramos",    0)
    ps.step2SuffixTree.Add("íramos",    0).Add("ávamos",    0)
    ps.step2SuffixTree.Add("emos",      0).Add("aremos",    0)
    ps.step2SuffixTree.Add("eremos",    0).Add("iremos",    0)
    ps.step2SuffixTree.Add("ássemos",   0).Add("êssemos",   0)
    ps.step2SuffixTree.Add("íssemos",   0).Add("imos",      0)
    ps.step2SuffixTree.Add("armos",     0).Add("ermos",     0)
    ps.step2SuffixTree.Add("irmos",     0).Add("eu",        0)
    ps.step2SuffixTree.Add("iu",        0).Add("ou",        0)
    ps.step2SuffixTree.Add("ira",       0).Add("iras",      0)

    // Load suffixes that are checked in Step 4.
    ps.step4SuffixTree = newSuffixTree()
    ps.step4SuffixTree.Add("os", 0).Add("a",0).Add("i", 0)
    ps.step4SuffixTree.Add("o", 0).Add("á",0).Add("í", 0)
    ps.step4SuffixTree.Add("ó", 0)

    // Load suffixes that are checked in Step 5.
    ps.step5SuffixTree = newSuffixTree()
    ps.step5SuffixTree.Add("e", 0).Add("é",0).Add("ê", 0)

    return ps
}

// Return true if letter is a vowel. Otherwise it should be treated
// as a consonant.
func (ps *PorterStemmer) isVowel (r rune) bool {
    _, ok := ps.vowels[r]
    return ok
}

// Expand nasalised vowels. 'ã' should be expanded to 'a~', with '~' being
// treated as a regular consonant.
func (ps *PorterStemmer) expandNasalisedVowels (word string) string {
    word = strings.Replace(word, "ã", "a~", -1)
    word = strings.Replace(word, "õ", "o~", -1)
    return word
}

// Contract nasalised vowels. 'a~' should be contracted to 'ã'.
func (ps *PorterStemmer) contractNasalisedVowels (word string) string {
    word = strings.Replace(word, "a~", "ã", -1)
    word = strings.Replace(word, "o~", "õ", -1)
    return word
}


// Find the remainder of the word after the first vowel, non-vowel
// sequence. This remainder is then returned as a string.
func (ps *PorterStemmer) r (word string) string {
    runes := []rune(word)
    for i := 0; i < len(runes)-1; i++ {
        if ps.isVowel(runes[i]) &&
           !ps.isVowel(runes[i+1]) {
            return string(runes[i+2:])
        }
    }
    return ""
}

// If the second letter is a consonant, RV is the region after the next
// following vowel, or if the first two letters are vowels, RV is the
// region after the next consonant, and otherwise (consonant-vowel case) RV
// is the region after the third letter. But RV is the end of the word if
// these positions cannot be found. 
func (ps *PorterStemmer) rv (word string) string {
    runes := []rune(word)
    if len(runes) < 3 {
        return ""
    }


    if !ps.isVowel(runes[1]) {
        for i := 2; i < len(runes); i++ {
            if ps.isVowel(runes[i]) {
                return string(runes[i+1:])
            }
        }
    } else if ps.isVowel(runes[0]) &&
        ps.isVowel(runes[1]) {
        for i := 2; i < len(runes); i++ {
            if !ps.isVowel(runes[i]) {
                return string(runes[i+1:])
            }
        }

        // If didnt return than RV is empty
        return ""
    }
    return string(runes[3:])
}

// This function executes the first step in the stemming algorithm. It
// checks and removes standard suffixes. This function returns the
// resultant word along with a boolean which is 'true' if the word was
// modified.
func (ps *PorterStemmer) step1 (word, r1, r2, rv string) (string, bool) {
    // Search for the longest among the known suffixes, perform the
    // action suitable to suffixe's group.
    suffix, group := ps.step1SuffixTree.LongestSuffix(word)

    if suffix == "" {
        return word, false
    }

    switch group {
    case 0:
        // eza   ezas   ico   ica   icos   icas   ismo   ismos   ável
        // ível   ista   istas   oso   osa   osos   osas   amento   amentos
        // imento   imentos   adora   ador   aça~o   adoras   adores
        // aço~es   ante   antes   ância
        //
        // Delete if in R2
        if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid], true
        }

    case 1:
        // logía   logías
        //
        // Replace with 'log' if in R2 
        if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid]+"log", true
        }

    case 2:
        // ución   uciones
        //
        // Replace with 'u' if in R2 
        if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid]+"u", true
        }

    case 3:
        // ência   ências
        //
        // Replace with 'ente' if in R2 
        if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid]+"ente", true
        }

    case 4:
        // amente
        //
        // Delete if in R1
        // If preceded by 'iv', delete if in R2 (and if further preceded by
        // 'at', delete if in R2), otherwise,
        // If preceded by 'os', 'ic' or 'ad', delete if in R2 
        res := word
        mod := false
        if strings.HasSuffix(r1, suffix) {
            lid := strings.LastIndex(word, suffix)
            res = word[:lid]
            mod = true
        }

        if strings.HasSuffix(r2, "iv"+suffix) {
            lid := strings.LastIndex(res, "iv")
            res = res[:lid]
            if strings.HasSuffix(r2, "ativ"+suffix) {
                lid := strings.LastIndex(res, "at")
                res = res[:lid]
            }
        } else if strings.HasSuffix(r2, "os"+suffix) {
            lid := strings.LastIndex(res, "os")
            res = res[:lid]
        } else if strings.HasSuffix(r2, "ic"+suffix) {
            lid := strings.LastIndex(res, "ic")
            res = res[:lid]
        } else if strings.HasSuffix(r2, "ad"+suffix) {
             lid := strings.LastIndex(res, "ad")
            res = res[:lid]
        }
        return res, mod

    case 5:
        // mente
        //
        // Delete if in R2
        // If preceded by 'ante', 'avel' or 'ível', delete if in R2
        if strings.HasSuffix(r2, "ante" +suffix) {
            lid := strings.LastIndex(word, "ante"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, "avel" +suffix) {
            lid := strings.LastIndex(word, "avel"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, "ível" +suffix) {
            lid := strings.LastIndex(word, "ível"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid], true
        }

    case 6:
        // idade   idades
        //
        // Delete if in R2
        // If preceded by 'abil', 'ic' or 'iv', delete if in R2
        if strings.HasSuffix(r2, "abil" +suffix) {
            lid := strings.LastIndex(word, "abil"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, "ic" +suffix) {
            lid := strings.LastIndex(word, "ic"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, "iv" +suffix) {
            lid := strings.LastIndex(word, "iv"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid], true
        }

    case 7:
        // iva   ivo   ivas   ivos
        // Delete if in R2
        // If preceded by 'at', delete if in R2 
        if strings.HasSuffix(r2, "at" +suffix) {
            lid := strings.LastIndex(word, "at"+suffix)
            return word[:lid], true
        } else if strings.HasSuffix(r2, suffix) {
            lid := strings.LastIndex(word, suffix)
            return word[:lid], true
        }

    case 8:
        // ira   iras
        //
        // Replace with 'ir' if in RV and preceded by 'e'
        if strings.HasSuffix(rv, suffix) {
            if strings.HasSuffix(word, "e" + suffix) {
                lid := strings.LastIndex(word, suffix)
                return word[:lid]+"ir", true
            }
        }
    }

    return word, false
}

// Second step in the portuguese stemming porter algorithm. This
// function removes verb suffixes and returns the resultant word and a
// boolean indicating if the word was modified.
func (ps *PorterStemmer) step2 (word, r1, r2, rv string) (string, bool) {
    // Search for the longest among the known suffixes in RV, if found
    // delete.
    suffix, _ := ps.step2SuffixTree.LongestSuffix(rv)

    if suffix == "" {
        return word, false
    }

    lid := strings.LastIndex(word, suffix)
    return word[:lid], true
}

// Third step in the stemming process. Delete suffix 'i' if in RV and
// preceded by 'c'. Returns the resultant word and a boolean indicating
// if the word was modified.
func (ps *PorterStemmer) step3 (word, r1, r2, rv string) (string, bool) {
    // Delete suffix 'i' if in RV and preceded by 'c'
    if strings.HasSuffix(word, "ci") && strings.HasSuffix(rv, "i") {
        return word[:len(word)-1], true
    }
    return word, false
}

// Forth step. Removes residual suffixes. Returns the resultant word and
// a boolean indicating if the word was modified.
func (ps *PorterStemmer) step4 (word, r1, r2, rv string) (string, bool) {
    // If the word ends with one of the suffixes
    // os   a   i   o   á   í   ó
    // in RV, delete it
    suffix, _ := ps.step4SuffixTree.LongestSuffix(rv)

    if suffix == "" {
        return word, false
    }

    lid := strings.LastIndex(word, suffix)
    return word[:lid], true
}

// Fifth step. Returns the resultant word and a boolean indicating if
// the word was modified.
func (ps *PorterStemmer) step5 (word, r1, r2, rv string) (string, bool) {
    // If the word ends with one of
    // e   é   ê
    // in RV, delete it, and if preceded by 'gu' (or 'ci') with the 'u'
    // (or 'i') in RV, delete the u (or i). 
    // Or if the word ends 'ç' remove the cedilla
    suffix, _ := ps.step5SuffixTree.LongestSuffix(rv)

    if suffix == "" {
        // Check if word ends with 'ç'
        if (strings.HasSuffix(word, "ç")) {
            lid := strings.LastIndex(word, "ç")
            return word[:lid]+"c", true
        } else {
            return word, false
        }
    }

    if strings.HasSuffix(rv, "u"+suffix) &&
    strings.HasSuffix(word, "gu"+suffix) {
        lid := strings.LastIndex(word, "u"+suffix)
        return word[:lid], true
    } else if strings.HasSuffix(rv, "i"+suffix) &&
    strings.HasSuffix(word, "ci"+suffix) {
        lid := strings.LastIndex(word, "i"+suffix)
        return word[:lid], true
    }

    lid := strings.LastIndex(word, suffix)
    return word[:lid], true
}

// Stem executes all steps necessary to obtain a given word's stem. This
// function is used for portuguese stemming only.
func (ps *PorterStemmer) Stem (word string) string {
    stem := ps.expandNasalisedVowels(word)
    modified := false
    r1 := ps.r(stem)
    r2 := ps.r(r1)
    rv := ps.rv(stem)

    // Always do step 1.
    stem, modified = ps.step1(stem, r1, r2, rv)

    // Do step 2 if no ending was removed by step 1.
    if !modified {
        stem, modified = ps.step2(stem, r1, r2, rv)
    }

    // Update R1, R2, RV if modified
    if modified {
        // If the last step to be obeyed — either step 1 or 2 — altered the
        // word, do step 3.
        r1 = ps.r(stem)
        r2 = ps.r(r1)
        rv = ps.rv(stem)

        stem, modified = ps.step3(stem, r1, r2, rv)
    } else {
        // Alternatively, if neither steps 1 nor 2 altered the word, 
        // do step 4.
        stem, modified = ps.step4(stem, r1, r2, rv)
    }

    if modified {
        r1 = ps.r(stem)
        r2 = ps.r(r1)
        rv = ps.rv(stem)
    }

    // Always do step 5.
    stem, modified = ps.step5(stem, r1, r2, rv)
    stem = ps.contractNasalisedVowels(stem)
    return stem
}
