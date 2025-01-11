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

		os.WriteFile("./output/Add_output.txt", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler MaxL.asm", func(t *testing.T) {
		got := Assemble("./test/MaxL.asm")
		want := fileRead("./test/MaxL_binary.txt")

		os.WriteFile("./output/MaxL_output.txt", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler PongL.asm", func(t *testing.T) {
		got := Assemble("./test/PongL.asm")
		want := fileRead("./test/PongL_binary.txt")

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler RectL.asm", func(t *testing.T) {
		got := Assemble("./test/RectL.asm")
		want := fileRead("./test/RectL_binary.txt")

		os.WriteFile("./output/RectL_output.txt", []byte(got), os.ModeAppend)

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