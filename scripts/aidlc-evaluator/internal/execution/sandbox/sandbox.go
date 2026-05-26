package sandbox

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	dockerclient "github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/mach4-braai/aidlc-workflows/aidlc-evaluator/internal/execution"
)

// Config describes Docker sandbox parameters.
type Config struct {
	Image       string
	Dockerfile  string
	WorkDir     string
	MemoryLimit int64
	CPUQuota    int64
	NetworkMode string
}

const defaultImage = "aidlc-sandbox:latest"

// DefaultConfig returns the default sandbox configuration.
func DefaultConfig() Config {
	dockerfilePath := findDockerfile()
	return Config{
		Image:       defaultImage,
		Dockerfile:  dockerfilePath,
		WorkDir:     "/workspace",
		MemoryLimit: 512 * 1024 * 1024, // 512 MB
		CPUQuota:    50000,              // 50% of one CPU
		NetworkMode: "none",
	}
}

func findDockerfile() string {
	candidates := []string{
		"docker/sandbox/Dockerfile",
		"../../../docker/sandbox/Dockerfile",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	return "docker/sandbox/Dockerfile"
}

// Build builds the sandbox Docker image from the configured Dockerfile.
func Build(cfg Config) error {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()
	buildCtx, err := createBuildContext(cfg.Dockerfile)
	if err != nil {
		return fmt.Errorf("build context: %w", err)
	}

	resp, err := cli.ImageBuild(ctx, buildCtx, types.ImageBuildOptions{
		Tags:       []string{cfg.Image},
		Dockerfile: "Dockerfile",
		Remove:     true,
	})
	if err != nil {
		return fmt.Errorf("image build: %w", err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return nil
}

// Run executes cmd inside a fresh sandbox container and returns results.
func Run(cfg Config, cmd string) (execution.CommandResult, error) {
	cli, err := dockerclient.NewClientWithOpts(dockerclient.FromEnv, dockerclient.WithAPIVersionNegotiation())
	if err != nil {
		return execution.CommandResult{}, fmt.Errorf("docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      cfg.Image,
		Cmd:        []string{"/bin/sh", "-c", cmd},
		WorkingDir: cfg.WorkDir,
	}, &container.HostConfig{
		Resources: container.Resources{
			Memory:   cfg.MemoryLimit,
			CPUQuota: cfg.CPUQuota,
		},
		NetworkMode: container.NetworkMode(cfg.NetworkMode),
		AutoRemove:  false,
	}, nil, nil, "")
	if err != nil {
		return execution.CommandResult{}, fmt.Errorf("container create: %w", err)
	}
	containerID := resp.ID
	defer cli.ContainerRemove(ctx, containerID, container.RemoveOptions{Force: true})

	if err := cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return execution.CommandResult{}, fmt.Errorf("container start: %w", err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	var exitCode int
	select {
	case err := <-errCh:
		if err != nil {
			return execution.CommandResult{}, fmt.Errorf("container wait: %w", err)
		}
	case status := <-statusCh:
		exitCode = int(status.StatusCode)
	}

	logs, err := cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return execution.CommandResult{ExitCode: exitCode}, nil
	}
	defer logs.Close()

	var stdout, stderr strings.Builder
	buf := make([]byte, 8*1024)
	for {
		n, err := logs.Read(buf)
		if n > 0 {
			stdout.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}

	return execution.CommandResult{
		ExitCode: exitCode,
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
	}, nil
}

func createBuildContext(dockerfilePath string) (io.Reader, error) {
	data, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	hdr := &tar.Header{Name: "Dockerfile", Size: int64(len(data)), Mode: 0644}
	if err := tw.WriteHeader(hdr); err != nil {
		return nil, err
	}
	if _, err := tw.Write(data); err != nil {
		return nil, err
	}
	tw.Close()
	return &buf, nil
}
