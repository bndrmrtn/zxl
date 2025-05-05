package extensions

import (
	"bufio"
	"fmt"
	"os/exec"
	"slices"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Manager struct {
	extensions []Extension
	logger     *zap.Logger
}

type Extension struct {
	Name   string
	Cmd    *exec.Cmd
	Port   string
	Conn   *grpc.ClientConn
	Client ExtensionManagerClient
}

func (m *Manager) AddExtension(name string, cmdPath string) {
	cmd := exec.Command(cmdPath)
	stdout, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		m.logger.Error("Error starting command", zap.String("command", cmdPath), zap.Error(err))
		return
	}

	scanner := bufio.NewScanner(stdout)
	var port string
	go func() {
		for scanner.Scan() {
			port = scanner.Text()
		}
	}()

	time.Sleep(2 * time.Second)

	conn, err := grpc.NewClient("127.0.0.1:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		m.logger.Error("Failed to connect to gRPC server", zap.String("port", port), zap.Error(err))
		return
	}

	ext := Extension{
		Name:   name,
		Cmd:    cmd,
		Port:   port,
		Conn:   conn,
		Client: NewExtensionManagerClient(conn),
	}
	m.extensions = append(m.extensions, ext)

	go func() {
		err := cmd.Wait()
		if err != nil {
			m.logger.Error("Command error", zap.String("name", name), zap.Error(err))
			m.RestartExtension(name, cmdPath)
		}
	}()
}

func (m *Manager) RestartExtension(name string, cmdPath string) {
	m.StopExtension(name)
	m.AddExtension(name, cmdPath)
}

func (m *Manager) StopExtension(name string) {
	for i, ext := range m.extensions {
		if ext.Name == name {
			m.logger.Info("Stopping extension", zap.String("name", name))
			ext.Cmd.Process.Kill()
			ext.Conn.Close()
			m.extensions = slices.Delete(m.extensions, i, i+1)
			return
		}
	}
	m.logger.Warn("Extension not found", zap.String("name", name))
}

func (m *Manager) StopAll() {
	for _, ext := range m.extensions {
		m.logger.Info("Stopping extension", zap.String("name", ext.Name))
		ext.Cmd.Process.Kill()
		ext.Conn.Close()
	}
	m.extensions = []Extension{}
}

func (m *Manager) GetExtensionClient(name string) (ExtensionManagerClient, error) {
	for _, ext := range m.extensions {
		if ext.Name == name {
			return ext.Client, nil
		}
	}
	return nil, fmt.Errorf("extension not found: %s", name)
}
