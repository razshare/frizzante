package assemblies

import (
	"bufio"
	"os"
	"strings"
)

func Parse(source *os.File, references map[string]map[string]*FunctionInfo, online func(line string)) {
	scanner := bufio.NewScanner(source)
	var fileName string
	var functionName string
	var start uint64
	var index uint64
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		online(line)
		if strings.HasPrefix(line, "TEXT") {
			if file, fileExists := references[fileName]; fileExists {
				if function, functionExists := file[functionName]; functionExists {
					function.BinarySize = index - start
				}
			}
			start = index + 1
			parts := strings.SplitN(line, " ", 3)
			count := len(parts)
			if count > 1 {
				functionName = parts[1]
			}
			if count > 2 {
				fileName = parts[2]
			}
			if file, fileExists := references[fileName]; !fileExists {
				file = map[string]*FunctionInfo{}
				references[fileName] = file
			}
		} else {
			file, fileExists := references[fileName]
			if fileName == "" || !fileExists {
				continue
			}
			function, functionExists := file[functionName]
			if !functionExists {
				function = &FunctionInfo{}
				file[functionName] = function
			}
			function.AssemblyContent += line + "\n"
		}
		index++
	}
	if file, fileExists := references[fileName]; fileExists {
		if function, functionExists := file[functionName]; functionExists {
			function.BinarySize = index - start
		}
	}
}
