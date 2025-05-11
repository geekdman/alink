package main

import (
	"alink/utils/base"
	"alink/utils/tikv"
	"alink/utils/zk"
	"fmt"
	"github.com/chzyer/readline"
	"strings"
)

type MultiStorageShell struct {
	backends map[string]func(string) base.BaseCRUDInterface
}

func NewMultiStorageShell() *MultiStorageShell {
	return &MultiStorageShell{
		backends: make(map[string]func(string) base.BaseCRUDInterface),
	}
}

func (m *MultiStorageShell) RegisterBackend(name string, factory func(string) base.BaseCRUDInterface) {
	m.backends[name] = factory
}

func (m *MultiStorageShell) Start() {
	// 创建 readline 实例
	rl, err := readline.NewEx(&readline.Config{
		Prompt:      ">>> ",
		HistoryFile: "/tmp/readline.tmp",
		//AutoComplete:    readline.NewAutoComplete(),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	fmt.Println("=== Multi-Storage CRUD Tool ===")
	fmt.Println("Version: 1.0")
	fmt.Println("Available Commands: es ip:port, tikv ip:port, zk ip:port, help, exit")

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}
		fmt.Println("===========")
		cmd := strings.TrimSpace(line)
		if cmd == "" {
			continue
		}

		switch cmd {
		case "exit":
			fmt.Println("Exiting...")
			return
		case "help":
			fmt.Println("Available Commands: es ip:port, tikv ip:port, zk ip:port, help, exit")
		default:
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
}

func main() {
	shell := NewMultiStorageShell()
	shell.RegisterBackend("zk", func(info string) base.BaseCRUDInterface {
		return zk.NewZkCrud(info)
	})
	//shell.RegisterBackend("es", func(info string) base.BaseCRUDInterface {
	//	return es.NewEsCrud(info)
	//})
	shell.RegisterBackend("tikv", func(info string) base.BaseCRUDInterface {
		return tikv.NewTikvCrud(info)
	})
	shell.Start()
}
