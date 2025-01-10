package main

import (
	"fmt"
	"os"
	"testing"
)

func TestAssembler(t *testing.T) {
	t.Run("assembler Add.asm", func(t *testing.T) {
		got := Assemble("./test/Add.asm")
		want := fileRead("./test/Add_binary.txt")

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}

		
	})
}

func fileRead(path string) string {
	content, err := os.ReadFile(path);
	
	if err != nil {
		fmt.Println("Error reading file: ", err)
		return ""
	}
	return string(content)
}