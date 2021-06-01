package main

import (
	"bufio"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/satyrius/gonx"
	"gopkg.in/mcuadros/go-syslog.v2/format"
	"strings"
	"time"
)

type NginxLogFormatter struct {
	parser *Parser
}

func NewNginxLogFormatter(format string) *NginxLogFormatter {
	return &NginxLogFormatter{
		parser: NewParser(format),
	}
}

type LogEntry struct {
	RemoteHost  string
	RemotePort  string
	HostName    string
	URI         string
	RequestTime time.Time
	Method      string
	RespCode    int
	RespLength  int
	UserAgent   string
	Referer     string
}

func (f *NginxLogFormatter) GetParser(line []byte) format.LogParser {
	f.parser.line = string(line)
	f.parser.entry = LogEntry{}
	f.parser.logParts = format.LogParts{}
	return f.parser
}

func (f *NginxLogFormatter) GetSplitFunc() bufio.SplitFunc {
	return nil
}

type Parser struct {
	line       string
	location   *time.Location
	entry      LogEntry
	format     string
	logParts   format.LogParts
	gonxParser *gonx.Parser
}

func (parser Parser) Parse() error {
	start := strings.Index(parser.line, "nginx: ")
	if start < 0 {
		return fmt.Errorf("text log parsing err: not nginx tag")
	}
	entry, err := parser.gonxParser.ParseString(parser.line[start+7:])
	if err != nil {
		return fmt.Errorf("text log parsing err: %w", err)
	}
	for key, value := range entry.Fields() {
		parser.logParts[key] = value
	}
	if err := mapstructure.Decode(entry, &parser.entry); err != nil {
		return err
	}
	return nil
}

func (parser Parser) Dump() format.LogParts {
	return parser.logParts
}

func (parser Parser) Location(location *time.Location) {
	parser.location = location
	return
}

func NewParser(format string) *Parser {
	return &Parser{
		location:   time.UTC,
		format:     format,
		gonxParser: gonx.NewParser(format),
	}
}
