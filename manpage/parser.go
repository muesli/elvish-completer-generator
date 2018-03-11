package manpage

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	sectionRx = regexp.MustCompile(`(?i)\A\.(PP|IT|SH|TP)\s?(.*)\z`)
	trimRx    = regexp.MustCompile(`(?i)\\f\w|\.RS\s\d+|\.(\w{2})\s?|\\|&`)
	flagTagRx = regexp.MustCompile(`(?i)\.?FL\s`)
)

var (
	// Debug enables/disables debug mode
	Debug = false
)

// Parser is a parser for man pages
type Parser struct {
	path string
	tags []*Tag
}

// New creates a new Parser
func New(path string) *Parser {
	parser := &Parser{
		path: path,
		tags: []*Tag{},
	}

	return parser
}

// GetOptions gets the options from the man page
func (p *Parser) GetOptions() map[string]string {
	opts := map[string]string{}

	for _, tag := range p.tags {
		for opt, desc := range tag.ToOptions() {
			opts[opt] = desc
		}
	}

	return opts
}

// Parse parses the man page
func (p *Parser) Parse() error {
	f, err := os.Open(p.path)
	if err != nil {
		return err
	}

	defer func() { _ = f.Close() }()

	var scanner *bufio.Scanner
	if strings.HasSuffix(p.path, ".gz") {
		rdr, err := gzip.NewReader(f)
		if err != nil {
			return err
		}
		defer func() { _ = rdr.Close() }()
		scanner = bufio.NewScanner(rdr)
	} else {
		scanner = bufio.NewScanner(f)
	}

	return p.scanAndProcess(scanner)
}

func (p *Parser) scanAndProcess(scanner *bufio.Scanner) error {
	var content []string
	var section string

	p.tags = []*Tag{}
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, `.\`) || line == "." {
			// skip comment
			continue
		}

		matches := sectionRx.FindStringSubmatch(line)
		if len(matches) <= 1 {
			content = append(content, cleanContent(line))
			continue
		}

		if section != "" {
			if Debug {
				fmt.Printf("Appending: %s %#v\n", section, content)
			}
			p.tags = append(p.tags, &Tag{Name: strings.ToUpper(section), Content: content})
		}

		if len(matches) > 1 {
			section = matches[1]
			content = []string{}
			if matches[2] != "" {
				content = append(content, cleanContent(matches[2]))
			}
		}

	}

	return nil
}

func cleanContent(content string) string {
	content = flagTagRx.ReplaceAllString(content, "-") // unescape flags
	content = trimRx.ReplaceAllString(content, "")     // remove extra tags

	return strings.TrimRight(content, " ") // some lines have extra spaces at the end
}
