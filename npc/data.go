package npc

import (
	"bufio"
	"embed"
	"encoding/gob"
	"io"
	"os"
	"regexp"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func IsRetailResponse(name, s string) bool {
	return IsNPC(name) && retailNPCs[name].Match(s)
}

func IsNPC(name string) bool {
	_, ok := retailNPCs[name]
	return ok
}

type ResponseMap map[string]Responses

func (m ResponseMap) Match(name, s string) bool {
	return m[name].Match(s)
}

type Responses struct {
	Name     string
	FileName string
	Static   map[string]bool
	Dynamic  []DynamicResponse
	Includes []string
}

func (r Responses) Match(s string) bool {
	if r.Static != nil && r.Static[s] {
		return true
	}

	for _, d := range r.Dynamic {
		if d.Match(s) {
			return true
		}
	}

	for _, include := range r.Includes {
		if retailIncludes[include].Match(s) {
			return true
		}
	}

	return false
}

type DynamicResponse struct {
	Prefix   string
	RE       string
	cachedRE *regexp.Regexp
}

func (r DynamicResponse) Match(s string) bool {
	if !strings.HasPrefix(s, r.Prefix) {
		return false
	}

	if r.cachedRE == nil {
		r.cachedRE = regexp.MustCompile(r.RE)
	}

	return r.cachedRE.MatchString(s)
}

//go:embed retaildata/*
var retailData embed.FS

var retailNPCs map[string]Responses
var retailIncludes map[string]Responses
var nameRE = regexp.MustCompile(`^\s*Name = "([^"]+)"`)
var includeRE = regexp.MustCompile(`^\s*@"([^"]+)"`)

func init() {
	retailNPCs = map[string]Responses{}
	retailIncludes = map[string]Responses{}

	if f, err := retailData.Open("retaildata/__npcs.gob"); err == nil {
		dec := gob.NewDecoder(f)
		dec.Decode(&retailNPCs)
	}

	if f, err := retailData.Open("retaildata/__includes.gob"); err == nil {
		dec := gob.NewDecoder(f)
		dec.Decode(&retailIncludes)
	}

	if len(retailNPCs) > 0 && len(retailIncludes) > 0 {
		return
	}

	files, _ := retailData.ReadDir("retaildata")
	for _, file := range files {
		isInclude := strings.HasSuffix(file.Name(), ".ndb")
		isNPC := strings.HasSuffix(file.Name(), ".npc")
		if !isInclude && !isNPC {
			continue
		}

		f, _ := retailData.Open("retaildata/" + file.Name())
		uReader := transform.NewReader(f, charmap.ISO8859_1.NewDecoder())
		r := parseFile(uReader, file.Name(), isInclude)

		if isInclude {
			retailIncludes[r.FileName] = r
		} else {
			retailNPCs[r.Name] = r
		}

		f.Close()
	}

	if f, err := os.Create("npc/retaildata/__npcs.gob"); err == nil {
		enc := gob.NewEncoder(f)
		enc.Encode(retailNPCs)
	}

	if f, err := os.Create("npc/retaildata/__includes.gob"); err == nil {
		enc := gob.NewEncoder(f)
		enc.Encode(retailIncludes)
	}
}

func parseFile(r io.Reader, fName string, isInclude bool) Responses {
	resp := Responses{
		Static:   map[string]bool{},
		FileName: fName,
	}

	scanner := bufio.NewScanner(r)
	needsName := !isInclude

	for scanner.Scan() {
		line := scanner.Text()
		if needsName && nameRE.MatchString(line) {
			needsName = false
			resp.Name = nameRE.FindStringSubmatch(line)[1]
			continue
		}

		if strings.Contains(line, "->") {
			line = strings.Split(line, "->")[1]
			responses := strings.Split(line, `"`)
			if len(responses) == 1 {
				continue
			}

			for i := 1; i < len(responses); i += 2 {
				r := responses[i]
				if strings.Contains(r, "%") {
					prefix := strings.Split(r, "%")[0]
					r = regexp.QuoteMeta(r)
					r = strings.ReplaceAll(r, "%P", `\d+`)
					r = strings.ReplaceAll(r, "%A", `\d+`)
					r = strings.ReplaceAll(r, "%T", `\d{2}:\d{2}`)
					r = strings.ReplaceAll(r, "%N", `[A-Za-z' -]+`)

					resp.Dynamic = append(resp.Dynamic, DynamicResponse{
						Prefix: prefix,
						RE:     r,
					})
				} else {
					resp.Static[r] = true
				}
			}
		}

		if strings.HasPrefix(line, "@") {
			resp.Includes = append(resp.Includes, includeRE.FindStringSubmatch(line)[1])
		}
	}
	return resp
}
