package exec

import (
	"bufio"
	"io"
	"os/exec"
)

// Command 执行
func Command(name string, arg ...string) (*exec.Cmd, []string, error) {
	cmd := exec.Command(name, arg...)
	out := []string{}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return cmd, nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return cmd, nil, err
	}

	err = cmd.Start()
	if err != nil {
		return cmd, nil, err
	}

	for _, rd := range []io.ReadCloser{stdout, stderr} {
		reader := bufio.NewReader(rd)
		for {
			line, err := reader.ReadString('\n')
			if err == io.EOF {
				break
			}
			out = append(out, line)
		}
	}

	err = cmd.Wait()
	return cmd, out, err
}
