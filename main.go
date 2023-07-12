package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"
	"io/ioutil"
	"os"
	"strings"
	"encoding/hex"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func monitorMemfdSymlink() {
	fmt.Println("Monitoring for open memfds in processes...")

	for {
		procs, err := ioutil.ReadDir("/proc")
		if err != nil {
			fmt.Printf("Failed to read /proc directory: %v\n", err)
			continue
		}

		for _, proc := range procs {
			if !proc.IsDir() {
				continue
			}

			pid := proc.Name()
			exePath := filepath.Join("/proc", pid, "exe")

			exe, err := os.Readlink(exePath)
			if err != nil {
				continue
			}

			if strings.HasPrefix(exe, "/memfd") {
				fmt.Printf("Found open memfd in PID: %s [%s]\n", pid, time.Now().Format("2006-01-02 15:04:05"))
			}
		}

		time.Sleep(5 * time.Second)
	}
}


func generatePythonScript(binaryData []byte, processName string) bool {
	script := "#!/usr/bin/env python3\n"
	script += "import os\n"
	script += "import sys\n\n"
	script += "fd = os.memfd_create('', os.MFD_CLOEXEC)\n"
	script += "if fd == -1:\n"
	script += "    sys.exit('[-] failed to create memfd')\n"
	script += "else:\n"
	script += "    print('[+] created memfd\\n\\t[+] fd:', fd, '\\n\\t[+] PID:', os.getpid())\n\n"
	script += "f = os.fdopen(fd, 'wb')\n\n"

	chunkSize := 32
	numChunks := len(binaryData) / chunkSize
	if len(binaryData)%chunkSize != 0 {
		numChunks++
	}

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(binaryData) {
			end = len(binaryData)
		}
		chunk := binaryData[start:end]
		hexChunk := hex.EncodeToString(chunk)
		script += fmt.Sprintf("f.write(bytes.fromhex('%s'))\n", hexChunk)
	}

	script += "\nprint('[+] wrote ELF binary to memory')\n\n"
	script += "print('[*] executing binary..')\n"
	script += "os.execv(f\"/proc/{os.getpid()}/fd/{fd}\", ['" + processName + "'])\n"

	f, err := os.Create("memfd.py")
	check(err)

	defer f.Close()

	_, err = f.WriteString(script)
	check(err)

	fmt.Println("wrote to memfd.py")

	return true
}

func generatePerlScript(binaryData []byte, processName string) bool {
	script := "#!/usr/bin/env perl\n"
	script += "use warnings;\n"
	script += "use strict;\n\n"
	script += "my $name = \"\";\n"
	script += "my $fd = syscall(319, $name, 1);\n"
	script += "if (-1 == $fd) {\n"
	script += "    die \"[-] failed to create memfd\";\n"
	script += "} else {\n"
	script += " print \"[+] created memfd\n\t[+] fd: $fd\n\t[+] PID: $$\n\n\"};\n";
	script += "open(my $FH, '>&='.$fd) or die \"open: $!\";\n"
	script += "select((select($FH), $|=1)[0]);\n\n"

	chunkSize := 32
	numChunks := len(binaryData) / chunkSize
	if len(binaryData)%chunkSize != 0 {
		numChunks++
	}

	for i := 0; i < numChunks; i++ {
		start := i * chunkSize
		end := start + chunkSize
		if end > len(binaryData) {
			end = len(binaryData)
		}
		chunk := binaryData[start:end]
		hexChunk := hex.EncodeToString(chunk)
		script += fmt.Sprintf("print $FH pack 'H*', '%s';\n", hexChunk)
	}

	script += "print \"[+] wrote ELF binary to memory\\n\";\n\n"
	script += "print \"[*] executing binary..\\n\";\n"
	script += "exec {\"/proc/$$/fd/$fd\"} \"%s\"\n"

	err := ioutil.WriteFile("memfd.pl", []byte(fmt.Sprintf(script, processName)), 0644)
	check(err)

	fmt.Println("wrote to memfd.pl")
	return true
}


func main() {
	binaryPath := flag.String("binary", "", "Path to the binary file")
	language := flag.String("language", "python", "Language for script generation (python, perl)")
	processName := flag.String("process", "", "Name of the process")
	monitor := flag.Bool("m", false, "Monitor for /memfd symlinks")
	flag.Parse()

	if *monitor {
		monitorMemfdSymlink()
		return
	}


	if *binaryPath == "" || *processName == "" {
		fmt.Println("Usage: memx -binary <binary_file> -language <language> -process <process_name>")
		os.Exit(1)
	}

	binaryData, err := ioutil.ReadFile(*binaryPath)
	if err != nil {
		fmt.Printf("Failed to read binary file: %v\n", err)
		os.Exit(1)
	}

	var script string
	switch strings.ToLower(*language) {
	case "python":
		_ = generatePythonScript(binaryData, *processName)
	case "perl":
		_ = generatePerlScript(binaryData, *processName)
	default:
		fmt.Println("Invalid language. Supported languages are python and perl.")
		os.Exit(1)
	}

	fmt.Println(script)
}

