package main

import "github.com/arisinghackers/goxploit/pkg/msfrpc"

func main() {
	msf_rpc_library_generator := msfrpc.MsfLibraryGenerator{}

	err := msf_rpc_library_generator.GenerateLibrary()
	if err != nil {
		panic(err)
	}

}
