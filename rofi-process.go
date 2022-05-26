package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

var last *RofiBlocksEvent

type RofiProcess struct {
	args    []string // just a list of strings for now
	ctx     context.Context
	cancel  context.CancelFunc
	command *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
}

func StartRofiProcess(args ...string) (*RofiProcess, error) {
	ctx, cancel := context.WithCancel(context.Background())
	cmdArgs := append([]string{"-show", "blocks", "-modes", "blocks,keys"}, args...)
	cmd := exec.CommandContext(ctx, "rofi", cmdArgs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return nil, err
	}
	err = cmd.Start()
	if err != nil {
		cancel()
		return nil, err
	}
	rofi := &RofiProcess{
		args:    args,
		ctx:     ctx,
		cancel:  cancel,
		command: cmd,
		stdin:   stdin,
		stdout:  stdout,
	}
	logger.Printf("Started Rofi (pid %d)\n", rofi.command.Process.Pid)
	go listenProcessExit(rofi)
	return rofi, nil
}

func (p *RofiProcess) Stop() (*os.ProcessState, error) {
	// runs p.command.Process.Kill(),
	// which does not wait for the process to exit
	p.cancel()
	// block until cancel is complete and close stdin/stdout
	err := p.command.Wait()
	return p.command.ProcessState, err
}

func (p *RofiProcess) Send(message string) error {
	messageBytes := []byte(message)
	bytesWritten, err := p.stdin.Write(messageBytes)
	if bytesWritten != len(messageBytes) {
		return fmt.Errorf("Only wrote %d bytes of %d-byte message; %v", bytesWritten, len(messageBytes), err)
	}
	return err
}

func listenProcessExit(rofi *RofiProcess) {
	state, err := rofi.command.Process.Wait()
	logger.Printf("Rofi process %d exited; exit code: %d; exit err: %+v\n", rofi.command.Process.Pid, state.ExitCode(), err)
	maybeHandleRofiExit(state, err)
}
func sendOutput(rofi *RofiProcess, outputReceiver chan RofiBlocksOutput, receiverReady chan int) {
	for {
		select {
		case <-rofi.ctx.Done():
			return
		default:
			logger.Println("Receiver ready")
			receiverReady <- rofi.command.Process.Pid
			output := <-outputReceiver
			logger.Printf("Got output to send to Rofi; %s", output.ToString())
			rofi.Send(output.ToString())
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func readEvents(rofi *RofiProcess, eventsChan chan RofiBlocksEvent) {
	scanner := bufio.NewScanner(rofi.stdout)
	for scanner.Scan() {
		lineIn := scanner.Bytes()
		// why do we get zero length buffers?
		if len(lineIn) > 0 {
			event := RofiBlocksEvent{Prev: last}
			err := json.Unmarshal(lineIn, &event)
			if err != nil {
				logger.Printf("Failed to parse Rofi output:\n%s\n", string(lineIn))
			} else {
				// this will block if nothing reads
				eventsChan <- event
				last = &event
			}
		}
	}
	// exit when the scanner is done, which should happen when the process is killed
}

func maybeHandleRofiExit(state *os.ProcessState, err error) {
	switch state.ExitCode() {
	case 0:
		if err == nil {
			os.Exit(0)
		}
	case 65:
		logger.Println("Rofi displayed an error and user exited")
		os.Exit(0)
	}
}
