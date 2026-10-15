package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/rix4uni/paramfilter/banner"
	"github.com/spf13/pflag"
)

var reCmd4 = regexp.MustCompile(`(?i).(jpg|jpeg|gif|css|tif|tiff|png|ttf|woff|woff2|icon|pdf|svg|txt|js)`)

func matchCmd1(line string) bool {
	dotIdx := strings.LastIndexByte(line, '.')
	if dotIdx == -1 || dotIdx == len(line)-1 {
		return false
	}
	switch strings.ToLower(line[dotIdx+1:]) {
	case "php", "asp", "aspx", "jsp", "jspx", "xml", "html":
		return true
	default:
		return false
	}
}

func matchCmd2(line string) bool {
	var offset int
	if strings.HasPrefix(line, "http://") {
		offset = 7
	} else if strings.HasPrefix(line, "https://") {
		offset = 8
	} else {
		return true
	}

	slashIdx := strings.IndexByte(line[offset:], '/')
	if slashIdx <= 0 {
		return true
	}

	remainder := line[offset+slashIdx+1:]
	return strings.IndexByte(remainder, '.') == -1
}

func cutSedF13(line string) string {
	slashCount := 0
	cutEnd := len(line)
	for i := 0; i < len(line); i++ {
		if line[i] == '/' {
			slashCount++
			if slashCount == 3 {
				cutEnd = i
				break
			}
		}
	}
	part := line[:cutEnd]
	for len(part) > 0 && part[len(part)-1] == '/' {
		part = part[:len(part)-1]
	}
	return part + "/"
}

func matchCmd4(line string) bool {
	return !reCmd4.MatchString(line)
}

func loadExisting(filename string, set map[string]struct{}) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, 256*1024)
	for {
		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")
			if line != "" {
				set[line] = struct{}{}
			}
		}
		if err != nil {
			break
		}
	}
}

func main() {
	var (
		endpointsFile string
		paramsFile    string
		silent        bool
		version       bool
	)

	pflag.StringVarP(&endpointsFile, "endpoints", "e", "endpoints.txt", "Output file for filtered endpoints")
	pflag.StringVarP(&paramsFile, "params", "p", "params.txt", "Output file for filtered parameters")
	pflag.BoolVar(&silent, "silent", false, "Silent mode")
	pflag.BoolVar(&version, "version", false, "Print the version of the tool and exit")
	pflag.Parse()

	if version {
		banner.PrintBanner()
		banner.PrintVersion()
		os.Exit(0)
	}

	if !silent {
		banner.PrintBanner()
	}

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return
	}

	seenEndpoints := make(map[string]struct{})
	seenParams := make(map[string]struct{})

	loadExisting(endpointsFile, seenEndpoints)
	loadExisting(paramsFile, seenParams)
	openFlag := os.O_CREATE | os.O_WRONLY | os.O_APPEND

	epFile, err := os.OpenFile(endpointsFile, openFlag, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating/opening endpoints file %s: %v\n", endpointsFile, err)
		os.Exit(1)
	}
	defer epFile.Close()
	epWriter := bufio.NewWriterSize(epFile, 256*1024)
	defer epWriter.Flush()

	paramFile, err := os.OpenFile(paramsFile, openFlag, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating/opening params file %s: %v\n", paramsFile, err)
		os.Exit(1)
	}
	defer paramFile.Close()
	paramWriter := bufio.NewWriterSize(paramFile, 256*1024)
	defer paramWriter.Flush()

	startTime := time.Now()
	totalLines := 0
	newEndpointsCount := 0
	newParamsCount := 0

	scanner := bufio.NewReaderSize(os.Stdin, 256*1024)

	for {
		line, err := scanner.ReadString('\n')
		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")
			line = strings.TrimSpace(line)
			if line != "" {
				totalLines++
				hasEqual := strings.IndexByte(line, '=') != -1

				if !hasEqual {
					if matchCmd1(line) {
						if _, exists := seenEndpoints[line]; !exists {
							seenEndpoints[line] = struct{}{}
							epWriter.WriteString(line)
							epWriter.WriteByte('\n')
							newEndpointsCount++
						}
					}
					if matchCmd2(line) {
						if _, exists := seenEndpoints[line]; !exists {
							seenEndpoints[line] = struct{}{}
							epWriter.WriteString(line)
							epWriter.WriteByte('\n')
							newEndpointsCount++
						}
					}
				} else {
					if matchCmd4(line) {
						if _, exists := seenParams[line]; !exists {
							seenParams[line] = struct{}{}
							paramWriter.WriteString(line)
							paramWriter.WriteByte('\n')
							newParamsCount++
						}
					}
				}

				origin := cutSedF13(line)
				if _, exists := seenEndpoints[origin]; !exists {
					seenEndpoints[origin] = struct{}{}
					epWriter.WriteString(origin)
					epWriter.WriteByte('\n')
					newEndpointsCount++
				}
			}
		}
		if err != nil {
			break
		}
	}

	epWriter.Flush()
	paramWriter.Flush()

	elapsed := time.Since(startTime)
	speed := float64(0)
	if elapsed.Seconds() > 0 {
		speed = float64(totalLines) / elapsed.Seconds()
	}
	fmt.Fprintf(os.Stderr, "[✓] Total URLs Processed : %d\n", totalLines)
	fmt.Fprintf(os.Stderr, "[✓] New Endpoints Saved  : %d -> %s\n", newEndpointsCount, endpointsFile)
	fmt.Fprintf(os.Stderr, "[✓] New Params Saved     : %d -> %s\n", newParamsCount, paramsFile)
	fmt.Fprintf(os.Stderr, "[✓] Time Elapsed         : %v (%.0f URLs/sec)\n", elapsed.Round(time.Millisecond), speed)
}
