package main

import (
	"fmt"
	"strings"

	"github.com/your/repo/utils/base"
	"github.com/your/repo/utils/es"
	"github.com/your/repo/utils/tikv"
	"github.com/your/repo/utils/zk"
)

type MultiStorageShell struct {
	.backends map[string]func(string) base.BaseCRUDInterface
}

func NewMultiStorageShell() *MultiStorageShell {
	return &MultiStorageShell{
		.backends: make(map[string]func(string) base.BaseCRUDInterface),
	}
}

func (m *MultiStorageShell) RegisterBackend(name string, factory func(string) base.BaseCRUDInterface) {
	m.backends[name] = factory
}

func (m *MultiStorageShell) Start() {
	fmt.Println("=== Multi-Storage CRUD Tool ===")
	fmt.Println("Version: 1.0")
	fmt.Println("Available Commands: es ip:port, tikv ip:port, zk ip:port, help, exit")

	for {
		fmt.Print(">>> ")
		var cmd string
		fmt.Scanln(&cmd)

		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		if cmd == "exit" {
			fmt.Println("Exiting...")
			break
		}

		if cmd == "help" {
			fmt.Println("Available Commands: es ip:port, tikv ip:port, zk ip:port, help, exit")
			continue
		}

		parts := strings.SplitN(cmd, " ", 2)
		if len(parts) < 2 {
			fmt.Println("Invalid command format. Use 'help' for more information.")
			continue
		}

		backendName := parts[0]
		connectionInfo := parts[1]

		factory, exists := m.backends[backendName]
		if !exists {
			fmt.Printf("Unsupported backend: %s\n", backendName)
			continue
		}

		backend := factory(connectionInfo)
		backend.Run()
	}
}

func main() {
	shell := NewMultiStorageShell()
	shell.RegisterBackend("es", func(info string) base.BaseCRUDInterface {
		return es.NewEsCrud(info)
	})
	shell.RegisterBackend("tikv", func(info string) base.BaseCRUDInterface {
		return tikv.NewTikvCrud(info)
	})
	shell.RegisterBackend("zk", func(info string) base.BaseCRUDInterface {
		return zk.NewZkCrud(info)
	})

	shell.Start()
}