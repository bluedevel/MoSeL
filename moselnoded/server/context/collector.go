package context

import (
	"github.com/WE-Development/mosel/api"
	"os"
	"os/exec"
	"bytes"
	"strings"
	"log"
	"io/ioutil"
)

type collector struct {
	scriptFolder string
	scripts      []string
}

func NewCollector() *collector {
	return &collector{
		scriptFolder: "/tmp/mosel",
		scripts: make([]string, 0),
	}
}

func (collector *collector) AddScript(name string, src []byte) error {
	filePath := collector.scriptFolder + "/" + name

	if _, err := mkdirIfNotExist(collector.scriptFolder, 0764); err != nil {
		return err
	}

	err := ioutil.WriteFile(filePath, src, 0664)

	if err != nil {
		log.Println(err)
		return err
	}

	collector.scripts = append(collector.scripts, name)
	log.Printf("Added script %s", name)
	return nil
}

func (collector *collector) FillNodeInfo(info *api.NodeInfo) {
	for _, script := range collector.scripts {
		executeScript(collector.scriptFolder + "/" + script, info)
	}
}

func executeScript(script string, info *api.NodeInfo) {
	cmd := exec.Command("/bin/bash", script)
	out := &bytes.Buffer{}
	cmd.Stdout = out

	res := make(map[string]string)
	for _, line := range
		strings.Split(out.String(), "\n") {
		parts := strings.SplitN(line, ":", 2)
		graph := parts[0]
		value := parts[1]
		res[graph] = value
	}
	(*info)[script] = res
}

func mkdirIfNotExist(path string, perm os.FileMode) (bool, error) {
	if ok, _ := exists(path); !ok {
		err := os.Mkdir(path, perm)
		return err != nil, err
	}
	return false, nil
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	log.Println(err)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}