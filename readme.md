# memx

memx is a tool for generating droppers that use `memfd_create()` to execute ELF binaries directly from memory, without dropping them on disk. This technique allows for stealthy execution of binaries, as they are never written to disk and leave no traces behind.

## Usage

```
memx -binary <binary_file> -language <language> -process <process_name>
```


### Options

- `-binary <binary_file>`: Path to the binary file to execute.
- `-language <language>`: Language for script generation. Supported languages are python and perl (default: python).
- `-process <process_name>`: Name of the process to execute the binary as.

## Monitoring for /memfd Symlinks

memx also provides an optional monitoring mode to detect open `memfd` instances in processes. This can be useful for detecting the use of `memfd_create()` in real-time.

```
memx -m
```

## Generating Python and Perl Scripts

memx can generate scripts that utilize `memfd_create()` to execute ELF binaries from memory, without ever storing them on disk. It supports Python and Perl scripts. There is also scope for adding support for additional programming languages in the future.

## Dependencies

memx requires Go (Golang) to be installed on the system.

## Installation

To install memx, follow these steps:

1. Install Go (Golang) on your system.
2. Clone the memx repository: `git clone https://github.com/zerodaywolf/memx.git`
3. Navigate to the memx directory: `cd memx`
4. Build the memx tool: `go build`
5. Run memx: `./memx`

## License

This tool is released under the MIT License. See [LICENSE](LICENSE) for more information.

## Disclaimer

This tool is intended for educational and security testing purposes only. The author is not responsible for any illegal or unauthorized use of this tool.

## Contributing

Contributions are welcome! If you have any ideas, suggestions, or improvements, feel free to open an issue or submit a pull request.
