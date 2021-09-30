package convert

import (
	"encoding/csv"
	"log"
	"os"
	"regexp"
	"strings"
)

const (
	posRightAnswers = 6
	posWrongAnswers = 8
	posAnswers      = 25
	posAttempts     = 26
	postfixAttempts = "_attempt"
	postfixAnswers  = "_answer"
	postfixCorrect  = "_correct"
)

func nameToSnake(str string) string {
	str = strings.Trim(str, " ")
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ToLower(str)
	return str
}

func contains(name string, assessments []string) bool {
	for _, i := range assessments {
		if name == i {
			return true
		}
	}
	return false
}

func parseCodioJoin(val string, postfix string, assessments []string) map[string]string {
	ret := make(map[string]string)
	header := ""
	body := ""
	nameMatch, err := regexp.Compile(`^.*?:$`)
	if err != nil {
		log.Fatal(err)
	}
	split := strings.Split(val, "\n")
	for _, text := range split {
		// text := scanner.Text()
		if nameMatch.MatchString(text) {
			assessmentName := strings.TrimSuffix(text, ":")
			assessmentName = nameToSnake(assessmentName)
			if contains(assessmentName, assessments) {
				if len(header) > 0 {
					ret[header] = body
				}
				header = assessmentName + postfix
				body = ""
				continue
			}
		}
		body = body + text + "\n"
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

// func join(record []string, keysMap []string, values map[string]string) []string {
// 	for _, v := range keysMap {
// 		value, ok := values[v]
// 		if !ok {
// 			value = ""
// 		}
// 		record = append(record, value)
// 	}
// 	return record
// }

func uniq(s []string) []string {
	m := make(map[string]bool)
	for _, item := range s {
		if len(item) == 0 {
			continue
		}
		m[item] = true
	}

	var result []string
	for item := range m {
		result = append(result, item)
	}
	return result
}

func extractAssessments(records [][]string) []string {
	assessments := []string{}
	for _, record := range records {
		assessmentsSplit := strings.Split(record[posRightAnswers]+","+record[posWrongAnswers], ",")
		assessments = append(assessments, assessmentsSplit...)
	}

	for k, v := range assessments {
		assessments[k] = nameToSnake(v)
	}
	assessments = uniq(assessments)
	return assessments
}

func parseCodioCorrectness(record []string) map[string]string {
	rightAnswers := strings.Split(record[posRightAnswers], ", ")
	wrongAnswers := strings.Split(record[posWrongAnswers], ", ")
	res := map[string]string{}
	for _, v := range rightAnswers {
		key := nameToSnake(v)
		res[key+postfixCorrect] = "correct"
	}
	for _, v := range wrongAnswers {
		key := nameToSnake(v)
		res[key+postfixCorrect] = "incorrect"
	}
	return res
}

func merge(ms ...map[string]string) map[string]string {
	res := map[string]string{}
	for _, m := range ms {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
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

	assessments := extractAssessments(records)
	for _, v := range assessments {
		keys = append(keys, v+postfixAnswers)
		keys = append(keys, v+postfixAttempts)
		keys = append(keys, v+postfixCorrect)
	}

	header = append(header, keys...)
	for k, v := range header {
		header[k] = nameToSnake(v)
	}

	csvOut.Write(header)

	for _, record := range records {
		answers := parseCodioJoin(record[posAnswers], postfixAnswers, assessments)
		attempts := parseCodioJoin(record[posAttempts], postfixAttempts, assessments)

		correct := parseCodioCorrectness(record)
		additionalFields := merge(answers, attempts, correct)
		for _, field := range keys {
			value, ok := additionalFields[field]
			if !ok {
				// fmt.Printf("missing \"%s\"\n", field)
				record = append(record, "")
			} else {
				record = append(record, value)
			}
		}
		csvOut.Write(record)
	}
}
