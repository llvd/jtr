package modules

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type JEval struct {
	dirPath string
	opEvalPath string
}

func (je *JEval) RmDir() {
	err := os.RemoveAll(je.dirPath)
	if err != nil {
		log.Fatal(err)
	}
}

func (je *JEval) RmEvalStrings() {
	var bcont string
	if dirExist(je.opEvalPath) {
		open, err := os.Open(je.opEvalPath)
		if err != nil {
			log.Fatal(err)
		}

		scanner := bufio.NewScanner(open)
		for scanner.Scan() {
			if strings.Contains(scanner.Text(), "evlsprt") {
				log.Println(scanner.Text())
			} else {
				bcont += scanner.Text()
			}
		}

		err = ioutil.WriteFile(je.opEvalPath, []byte(bcont), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func dirExist(dirpath string) bool {
	_, err := os.Stat(dirpath)
	if err != nil {
		return false
	}
	return true
}

func delRegEntry(pname string) {
	name := strings.ToLower(pname[1:])
	regPath := fmt.Sprintf("HKEY_CURRENT_USER\\Software\\JavaSoft\\Prefs\\jetbrains\\%s", name)
	cmd := exec.Command("reg", "delete", regPath, "/f")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return
	}
}

func getEvalPaths(pname []string) []JEval {
	var result []JEval
	udir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dirs, err := ioutil.ReadDir(udir)
	if err != nil {
		log.Fatal(err)
	}

	for _, maybeDir := range dirs {
		for _, name := range pname {
			if maybeDir.IsDir() && strings.HasPrefix(maybeDir.Name(), name) {
				result = append(result, JEval{
					dirPath: fmt.Sprintf("%s\\%s\\%s", udir, maybeDir.Name(), "config\\eval"),
					opEvalPath: fmt.Sprintf("%s\\%s\\%s", udir, maybeDir.Name(), "config\\options\\other.xml"),
				})
			}
			maybeDir.Name()
		}
	}
	return result
}

func DetectJbPrograms() {
	pnames := [...]string{".CLion", ".IntelliJ", ".GoLand", ".DataGrip", ".RubyMine", ".PyCharm", ".WebStorm"}
	for i, name := range getEvalPaths(pnames[:]) {
		log.Printf("Resetting %s", pnames[i][1:])
		name.RmEvalStrings()
		name.RmDir()
		delRegEntry(pnames[i])
	}
}