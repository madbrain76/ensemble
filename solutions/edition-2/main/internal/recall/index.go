// Package recall owns keyword retrieval, relevance judgment and attachment text.
// Its only application dependency is the common vocabulary and parent interfaces.
package recall

import (
	"math"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// An index belongs to one source so a large document directory cannot dilute
// rare terms in personal memories. Chunks retain original bytes for attribution.
type chunk struct {
	file, header, content string
	tf                    map[string]int
	length                int
}
type index struct {
	chunks  []chunk
	df      map[string]int
	average float64
}
type hit struct {
	chunk
	source string
	score  float64
}

func tokens(s string) []string {
	var out []string
	for _, word := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsNumber(r) }) {
		if utf8.RuneCountInString(word) >= 2 {
			out = append(out, word)
		}
	}
	return out
}

// Stop words remain in documents: dropping them there changes length statistics.
func queryTokens(s string) []string {
	const stops = " a an and are as at be been but by can could did do does for from had has have how if in into is it its may me my not of on or our should so that the their them then there these they this to was we were what when where which who will with would you your "
	var out []string
	seen := map[string]bool{}
	for _, word := range tokens(s) {
		if !strings.Contains(stops, " "+word+" ") && !seen[word] {
			out = append(out, word)
			seen[word] = true
		}
	}
	return out
}
func build(chunks []chunk) index {
	x := index{chunks: chunks, df: map[string]int{}}
	for i := range x.chunks {
		c := &x.chunks[i]
		c.tf = map[string]int{}
		for _, word := range tokens(c.content) {
			c.tf[word]++
			c.length++
		}
		for word := range c.tf {
			x.df[word]++
		}
		x.average += float64(c.length)
	}
	if len(chunks) > 0 {
		x.average /= float64(len(chunks))
	}
	return x
}
func (x index) search(source string, query []string, floor float64) []hit {
	var out []hit
	if x.average == 0 {
		return out
	}
	for _, c := range x.chunks {
		score := 0.0
		for _, word := range query {
			f := float64(c.tf[word])
			// Smoothed IDF stays positive even in a one-document archive.
			idf := math.Log(1 + (float64(len(x.chunks)-x.df[word])+0.5)/(float64(x.df[word])+0.5))
			score += idf * f * 2.2 / (f + 1.2*(0.25+0.75*float64(c.length)/x.average))
		}
		if score > 0 && score >= floor {
			out = append(out, hit{c, source, score})
		}
	}
	return out
}
func rank(hits []hit) {
	sort.SliceStable(hits, func(i, j int) bool { return hits[i].score > hits[j].score })
}

// Split only at the requested markdown depth. Preambles remain searchable, and
// a child heading carries its parent breadcrumb when a large section is divided.
func sections(text, marker, parent string) []chunk {
	var out []chunk
	header := parent
	var body strings.Builder
	flush := func() {
		if content := strings.TrimSpace(body.String()); content != "" {
			out = append(out, chunk{header: header, content: content})
		}
		body.Reset()
	}
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, marker) {
			flush()
			header = strings.TrimSpace(strings.TrimPrefix(line, marker))
			if parent != "" {
				header = parent + " > " + header
			}
		}
		body.WriteString(line)
		body.WriteByte('\n')
	}
	flush()
	return out
}
func chunks(file, text string) []chunk {
	var out []chunk
	for _, section := range sections(text, "## ", "") {
		parts := []chunk{section}
		if len(section.content) > 4096 {
			parts = sections(section.content, "### ", section.header)
		}
		for _, part := range parts {
			if len(part.content) <= 4096 {
				part.file = file
				out = append(out, part)
				continue
			}
			// Paragraphs are the last semantic boundary. An indivisible paragraph
			// may be larger; the attachment's hard byte cap still applies.
			body := ""
			for _, paragraph := range strings.Split(part.content, "\n\n") {
				if body != "" && len(body)+len(paragraph)+2 > 4096 {
					out = append(out, chunk{file: file, header: part.header, content: body})
					body = ""
				}
				if body != "" {
					body += "\n\n"
				}
				body += paragraph
			}
			if body != "" {
				out = append(out, chunk{file: file, header: part.header, content: body})
			}
		}
	}
	return out
}
