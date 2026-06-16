package clipboard

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

var ErrUnavailable = errors.New("clipboard provider unavailable")

type Copier interface {
	Copy(ctx context.Context, text string) error
}

type Provider struct {
	Name string
	Args []string
}

type CommandCopier struct {
	provider Provider
}

func NewCommandCopier(provider Provider) *CommandCopier {
	return &CommandCopier{provider: provider}
}

func (c *CommandCopier) Copy(ctx context.Context, text string) error {
	if _, err := exec.LookPath(c.provider.Name); err != nil {
		return fmt.Errorf("%w: %s", ErrUnavailable, c.provider.Name)
	}
	cmd := exec.CommandContext(ctx, c.provider.Name, c.provider.Args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	if _, err := stdin.Write([]byte(text)); err != nil {
		_ = stdin.Close()
		_ = cmd.Wait()
		return err
	}
	if err := stdin.Close(); err != nil {
		_ = cmd.Wait()
		return err
	}
	return cmd.Wait()
}

type MultiCopier struct {
	copiers []Copier
}

func NewMulti(copiers ...Copier) *MultiCopier {
	return &MultiCopier{copiers: copiers}
}

func NewLinuxCopier() *MultiCopier {
	return NewMulti(
		NewCommandCopier(Provider{Name: "wl-copy"}),
		NewCommandCopier(Provider{Name: "xclip", Args: []string{"-selection", "clipboard"}}),
		NewCommandCopier(Provider{Name: "xsel", Args: []string{"--clipboard", "--input"}}),
	)
}

func (m *MultiCopier) Copy(ctx context.Context, text string) error {
	var last error
	for _, copier := range m.copiers {
		if err := copier.Copy(ctx, text); err == nil {
			return nil
		} else {
			last = err
		}
	}
	if last != nil {
		return fmt.Errorf("%w: %v", ErrUnavailable, last)
	}
	return ErrUnavailable
}
