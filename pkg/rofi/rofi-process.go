package rofi

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

type RofiProcess struct {
	Args      []string // just a list of strings for now
	Output    *RofiBlocksOutput
	LastEvent Event
	ctx       context.Context
	cancel    context.CancelFunc
	command   *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
}

func MakeRofiProcess(args ...string) (*RofiProcess, error) {
	ctx, cancel := context.WithCancel(context.Background())
	r := &RofiProcess{
		Args:    append([]string{"-show", "blocks", "-modes", "blocks"}, args...),
		Output:  &RofiBlocksOutput{},
		ctx:     ctx,
		cancel:  cancel,
	}
	cmd := exec.CommandContext(ctx, "rofi", r.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	r.command = cmd
	r.stdin = stdin
	r.stdout = stdout
	return r, nil
}

func (r *RofiProcess) Start() error {
	if err := r.command.Start(); err != nil {
		return err
	}
	go r.sendOutput()
	return nil
}

func (r *RofiProcess) Stop() (*os.ProcessState, error) {
	// runs p.command.Process.Kill(),
	// which does not wait for the process to exit
	r.cancel()
	// block until cancel is complete and close stdin/stdout
	err := r.command.Wait()
	state := r.command.ProcessState
	return state, err
}

// Handle the case of the process dying before we call Stop()
func (r *RofiProcess) ListenProcessExit() (*os.ProcessState, error) {
	state, err := r.command.Process.Wait()
	r.cancel() // is this necessary?
	log.Printf("Rofi process %d exited; exit code: %d; exit err: %+v\n%v", r.command.Process.Pid, state.ExitCode(), err, state)
	return state, err
}

func (r *RofiProcess) ReadEvents(eventsChan chan<- Event) {
	scanner := bufio.NewScanner(r.stdout)
	for scanner.Scan() {
		lineIn := scanner.Bytes()
		// why do we get zero length buffers?
		if len(lineIn) < 1 {
			continue
		}
		event, err := r.readEvent(lineIn)
		if err != nil {
			log.Printf("Failed to parse Rofi event: %v: %s", err, string(lineIn))
			continue
		}
		// this will block if nothing reads; is that fine?
		eventsChan <- event
		r.LastEvent = event
	}
	// exit when the scanner is done, which should happen when the process is killed
}

func (r *RofiProcess) SendUpdates(updates <-chan []OutputUpdate) {
	for {
		select {
		case <-r.ctx.Done():
			return
		case changes := <-updates:
			r.Output.ApplyAll(changes)
		}
	}
}

func (r *RofiProcess) sendOutput() {
	send := func() {
		if r.Output.Changes > 0 {
			msg, err := r.Output.MarshalJson()
			if err != nil {
				log.Printf("ERROR: marshal output to json: %v", err)
				return
			}
			msg = append(msg, '\n')
			if _, err := r.stdin.Write(msg); err != nil {
				log.Printf("Write failed: %v", err)
				return
			}
			r.Output.Changes = 0
		}
	}

	send()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-time.After(50 * time.Millisecond):
			send()
		}
	}
}

func (r *RofiProcess) readEvent(data []byte) (Event, error) {
	var raw rawEvent
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	switch raw.Name {
	case "input change":
		return &InputChangeEvent{Value: raw.Value, prev: r.LastEvent}, nil
	case "custom key":
		id, err := strconv.Atoi(raw.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to read custom key event key id: %w", err)
		}
		return &CustomKeyEvent{KeyId: id, prev: r.LastEvent}, nil
	case "execute custom input":
		return &ExecCustomEntryEvent{Value: raw.Value, prev: r.LastEvent}, nil
	}
	idInt, err := strconv.Atoi(raw.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to read line id: %w", err)
	}
	id := LineId(idInt)
	switch raw.Name {
	case "active entry":
		return &ActiveEntryEvent{LineId: id, prev: r.LastEvent}, nil
	case "select entry":
		return &SelectEntryEvent{LineId: id, prev: r.LastEvent}, nil
	case "delete entry":
		return &DeleteEntryEvent{LineId: id, prev: r.LastEvent}, nil
	default:
		return nil, fmt.Errorf("received unknown event '%s'", raw.Name)
	}
}
