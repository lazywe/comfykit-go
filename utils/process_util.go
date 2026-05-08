package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type Process struct {
	cmd     *exec.Cmd
	stdout  *bufio.Scanner
	stderr  *bufio.Scanner
	done    chan struct{}
	mu      sync.Mutex
	running bool
}

func NewProcess(command string, args ...string) *Process {
	cmd := exec.Command(command, args...)
	return &Process{
		cmd:     cmd,
		done:    make(chan struct{}),
		running: false,
	}
}

func (p *Process) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return fmt.Errorf("process already running")
	}

	stdout, err := p.cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := p.cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := p.cmd.Start(); err != nil {
		return err
	}

	p.stdout = bufio.NewScanner(stdout)
	p.stderr = bufio.NewScanner(stderr)
	p.running = true

	go func() {
		for p.stdout.Scan() {
			fmt.Println(p.stdout.Text())
		}
	}()

	go func() {
		for p.stderr.Scan() {
			fmt.Println(p.stderr.Text())
		}
	}()

	go func() {
		p.cmd.Wait()
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
		close(p.done)
	}()

	return nil
}

func (p *Process) Wait() error {
	<-p.done
	if !p.cmd.ProcessState.Success() {
		return fmt.Errorf("process exited with code %d", p.cmd.ProcessState.ExitCode())
	}
	return nil
}

func (p *Process) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return fmt.Errorf("process not running")
	}

	if err := p.cmd.Process.Kill(); err != nil {
		return err
	}

	<-p.done
	return nil
}

func (p *Process) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}

func RunCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, string(output))
	}
	return string(output), nil
}

func RunCommandWithEnv(env map[string]string, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)

	currentEnv := os.Environ()
	for k, v := range env {
		currentEnv = append(currentEnv, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = currentEnv

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, string(output))
	}
	return string(output), nil
}

func FindProcessByName(name string) ([]int, error) {
	var pids []int

	switch runtime.GOOS {
	case "windows":
		output, err := RunCommand("tasklist", "/FI", fmt.Sprintf("IMAGENAME eq %s.exe", name))
		if err != nil {
			return nil, err
		}
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, name+".exe") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					var pid int
					if _, err := fmt.Sscanf(parts[1], "%d", &pid); err == nil {
						pids = append(pids, pid)
					}
				}
			}
		}
	default:
		output, err := RunCommand("pgrep", "-f", name)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			var pid int
			if _, err := fmt.Sscanf(line, "%d", &pid); err == nil {
				pids = append(pids, pid)
			}
		}
	}

	return pids, nil
}

func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}
