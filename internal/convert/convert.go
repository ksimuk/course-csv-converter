package convert

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
)

type parsed struct {
	answer   string
	attempts string
}

func parseCodioJoin(val string, postfix string) map[string]string {
	ret := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(val))
	header := ""
	body := ""
	nameMatch, err := regexp.Compile(`^.*?:$`)
	if err != nil {
		log.Fatal(err)
	}
	for scanner.Scan() {
		text := scanner.Text()
		if nameMatch.MatchString(text) {
			if len(header) > 0 {
				ret[header] = body
			}
			header = strings.TrimSuffix(text, ":") + postfix
			header = strings.ReplaceAll(header, "_", " ")
			body = ""
		} else {
			body = body + text + "\n"
		}
	}
	if len(header) > 0 {
		ret[header] = body
	}
	return ret
}

func getKeys(val map[string]string) []string {
	keys := make([]string, 0, len(val))
	for k := range val {
		keys = append(keys, k)
	}
	return keys
}

func join(record []string, keysMap []string, values map[string]string) []string {
	for _, v := range keysMap {
		value, ok := values[v]
		if !ok {
			value = ""
		}
		record = append(record, value)
	}
	return record
}

func uniq(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		m[item] = true
	}

	var result []string
	for item := range m {
		result = append(result, item)
	}
	return result
}

func Convert(srcFile string, dstFile string) {
	f, err := os.Open(srcFile)
	if err != nil {
		log.Fatal(err)
	}

	in := csv.NewReader(f)
	fOut, err := os.Create(dstFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fOut.Close()

	csvOut := csv.NewWriter(fOut)
	defer csvOut.Flush()

	header, err := in.Read() // read header
	if err != nil {
		log.Fatal(err)
	}

	records, err := in.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var keys []string

	for _, record := range records {
		answers := parseCodioJoin(record[25], "_anwser")
		attempts := parseCodioJoin(record[26], "_attempts")
		keys = append(keys, getKeys(answers)...)
		keys = append(keys, getKeys(attempts)...)
	}
	keys = uniq(keys)
	header = append(header, keys...)
	csvOut.Write(header)

	for _, record := range records {
		values := parseCodioJoin(record[25], "_anwser")
		for k, v := range parseCodioJoin(record[26], "_attempts") {
			values[k] = v
		}
		res := join(record, keys, values)
		csvOut.Write(res)
	}
}
