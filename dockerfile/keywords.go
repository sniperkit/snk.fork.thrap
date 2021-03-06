/*
Sniperkit-Bot
- Date: 2018-08-11 22:25:29.898780201 +0200 CEST m=+0.118184110
- Status: analyzed
*/

package dockerfile

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	// KeyComment is comment in a dockerfile
	KeyComment = "#"
	// KeyFrom is the dockerfile operation
	KeyFrom = "FROM"
	// KeyExpose is the dockerfile operation
	KeyExpose = "EXPOSE"
	// KeyRun is the dockerfile operation
	KeyRun = "RUN"
	// KeyCopy is the dockerfile operation
	KeyCopy = "COPY"
	// KeyAdd is the dockerfile operation
	KeyAdd = "ADD"
	// KeyLabel is the dockerfile operation
	KeyLabel = "LABEL"
	// KeyEnv is the dockerfile operation
	KeyEnv = "ENV"
	// KeyWorkDir is the dockerfile operation
	KeyWorkDir = "WORKDIR"
	// KeyVolume is the dockerfile operation
	KeyVolume = "VOLUME"
	// KeyUser is the dockerfile operation
	KeyUser = "USER"
	// KeyArg is the dockerfile operation
	KeyArg = "ARG"
	// KeyShell is the dockerfile operation
	KeyShell = "SHELL"
	// KeyHealthCheck is the dockerfile operation
	KeyHealthCheck = "HEALTHCHECK"
	// KeyStopSignal is the dockerfile operation
	KeyStopSignal = "STOPSIGNAL"
	// KeyOnBuild is the dockerfile operation
	KeyOnBuild = "ONBUILD"
	// KeyCmd is the dockerfile operation
	KeyCmd = "CMD"
	// KeyEntrypoint is the dockerfile operation
	KeyEntrypoint = "ENTRYPOINT"
)

var (
	errInvalidInstruction   = errors.New("invalid instruction")
	errUnsupportInstruction = errors.New("unsupported instruction")
)

// ParseInstruction parses a raw instruction into a concrete instruction
func ParseInstruction(r *RawInstruction) (Instruction, error) {
	switch r.Op {

	case KeyFrom:
		return ParseFrom(r.Data)
	case KeyArg:
		return ParseArg(r.Data)
	case KeyWorkDir:
		return ParseWorkDir(r.Data)
	case KeyExpose:
		return ParseExpose(r.Data)
	case KeyCopy:
		return ParseCopy(r.Data)
	case KeyComment:
		return ParseComment(r.Data)

	default:
		return nil, errUnsupportInstruction
	}
}

// Arg is the ARG keyword in a dockerfile
type Arg struct {
	Name string
}

// ParseArg parses a byte slice to a Arg struct
func ParseArg(b []byte) (*Arg, error) {
	p := strings.Split(string(b), " ")
	if len(p) != 1 {
		return nil, errors.New("must have one arg name")
	}

	return &Arg{Name: string(b)}, nil
}

// Key returns the instruction key
func (e *Arg) Key() string {
	return KeyArg
}

func (e *Arg) String() string {
	return e.Name
}

// Env is the ENV keyword in a dockerfile
type Env struct {
	Vars map[string]string
}

// Key returns the instruction key
func (e *Env) Key() string {
	return KeyEnv
}

// Clone returns a clone of Env
func (e *Env) Clone() *Env {
	c := &Env{Vars: make(map[string]string, len(e.Vars))}
	for k, v := range e.Vars {
		c.Vars[k] = v
	}
	return c
}

func (e *Env) String() string {
	var out string
	for k, v := range e.Vars {
		out += k + `="` + v + `" `
	}

	return out[:len(out)-1]
}

// WorkDir is the WORKDIR keyword in a dockerfile
type WorkDir struct {
	Path string
}

func (wd *WorkDir) String() string {
	return wd.Path
}

// Key returns the instruction key
func (wd *WorkDir) Key() string {
	return KeyWorkDir
}

// ParseWorkDir parses a workdir instruction from a dockerfile
func ParseWorkDir(b []byte) (*WorkDir, error) {
	if string(b) == "" {
		return nil, errors.New("WORKDIR not specified")
	}
	return &WorkDir{Path: string(b)}, nil
}

// Comment is a comment in a docker file.  This is a single line in itself
type Comment struct {
	Text string
}

func (cmt *Comment) String() string {
	return cmt.Text
}

// Key returns the instruction key
func (cmt *Comment) Key() string {
	return KeyComment
}

// ParseComment parse a comment in a dockerfile
func ParseComment(b []byte) (*Comment, error) {
	s := strings.TrimSpace(string(b))
	if s[0] != '#' {
		return nil, errors.New("invalid comment")
	}
	return &Comment{
		Text: strings.TrimSpace(s[1:]),
	}, nil
}

// Run is the RUN keyword in a dockerfile
type Run struct {
	Command string
}

func (run *Run) String() string {
	return run.Command
}

// Key returns the instruction key
func (run *Run) Key() string {
	return KeyRun
}

// Cmd is a parse docker CMD instruction
type Cmd struct {
	Command string
	Args    []string
}

func (cmd *Cmd) String() string {
	return cmd.Command + " " + strings.Join(cmd.Args, " ")
}

// Key returns the instruction key
func (cmd *Cmd) Key() string {
	return KeyCmd
}

// EntryPoint is the ENTRYPOINT keyword in a dockerfile
type EntryPoint struct {
	Command string
}

func (ep *EntryPoint) String() string {
	return ep.Command
}

// Key returns the instruction key
func (ep *EntryPoint) Key() string {
	return KeyEntrypoint
}

// Volume is the VOLUME keyword in a dockerfile
type Volume struct {
	Paths []string
}

// Key returns the instruction key
func (vol *Volume) Key() string {
	return KeyVolume
}

func (vol *Volume) String() string {
	if len(vol.Paths) == 1 {
		return vol.Paths[0]
	}
	return `["` + strings.Join(vol.Paths, `", "`) + `"]`
}

// Expose line
type Expose struct {
	Port  uint16
	Proto string
}

// Key returns the instruction key
func (expose *Expose) Key() string {
	return KeyExpose
}

func (expose *Expose) String() string {
	if expose.Port == 0 {
		return ""
	}

	out := strconv.FormatUint(uint64(expose.Port), 10)
	if len(expose.Proto) > 0 {
		out += "/" + expose.Proto
	}
	return out
}

// ParseExpose parses an expose instruction
func ParseExpose(b []byte) (*Expose, error) {
	parts := strings.Split(string(b), "/")
	l := len(parts)

	var (
		e   Expose
		p   uint64
		err error
	)

	switch l {
	case 1:
		p, err = strconv.ParseUint(parts[0], 10, 16)
		if err == nil {
			e.Port = uint16(p)
		}

	case 2:
		e.Proto = parts[1]
		p, err = strconv.ParseUint(parts[0], 10, 16)
		if err == nil {
			e.Port = uint16(p)
		}

	default:
		return nil, errors.Wrap(errInvalidInstruction, KeyExpose)

	}

	return &e, nil
}

// Copy line
type Copy struct {
	Source      string
	Destination string
	// Addtional copy options format: --option=value
	Options []string
}

// Key returns the instruction key
func (copy *Copy) Key() string {
	return KeyCopy
}

func (copy *Copy) String() string {
	return strings.Join(copy.Options, " ") + " " +
		copy.Source + " " +
		copy.Destination
}

// ParseCopy parses a COPY instruction
func ParseCopy(b []byte) (*Copy, error) {
	parts := strings.Split(string(b), " ")
	l := len(parts)

	if l < 2 {
		return nil, errors.Wrap(errInvalidInstruction, KeyCopy)
	}

	c := &Copy{
		Source:      parts[l-2],
		Destination: parts[l-1],
	}

	if len(parts) > 2 {
		c.Options = parts[:l-2]
	}

	return c, nil
}

// From line
type From struct {
	Image string
	As    string
}

// Key returns the instruction key
func (from *From) Key() string {
	return KeyFrom
}

func (from *From) String() string {
	if len(from.As) == 0 {
		return from.Image
	}

	return from.Image + " as " + from.As
}

// ParseFrom parses returns a from object
func ParseFrom(b []byte) (*From, error) {
	var (
		s     = string(b)
		parts = strings.Split(s, " ")
		l     = len(parts)
		from  *From
	)

	switch l {
	case 1:
		from = &From{Image: parts[0]}

	case 3:
		from = &From{
			Image: parts[0],
			As:    parts[2],
		}

	default:
		return nil, errors.Wrap(errInvalidInstruction, KeyFrom)
	}

	return from, nil
}
