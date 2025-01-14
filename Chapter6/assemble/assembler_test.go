package assemble

import (
	"fmt"
	"os"
	"testing"
)

func TestAssembler(t *testing.T) {
	t.Run("assembler Add.asm", func(t *testing.T) {
		got := Assemble("../test/Add.asm")
		want := fileRead("../test/compare/Add.hack")

		// os.WriteFile("./test/output/Add.hack", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler MaxL.asm", func(t *testing.T) {
		got := Assemble("../test/MaxL.asm")
		want := fileRead("../test/compare/MaxL.hack")

		// os.WriteFile("./test/output/MaxL_output.txt", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler PongL.asm", func(t *testing.T) {
		got := Assemble("../test/PongL.asm")
		want := fileRead("../test/compare/PongL.hack")

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler RectL.asm", func(t *testing.T) {
		got := Assemble("../test/RectL.asm")
		want := fileRead("../test/compare/RectL.hack")

		// os.WriteFile("./output/RectL.hack", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler Max.asm", func(t *testing.T) {
		got := Assemble("../test/Max.asm")
		want := fileRead("../test/compare/Max.hack")

		os.WriteFile("../test/output/Max.hack", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("assembler Rect.asm", func(t *testing.T) {
		got := Assemble("../test/Rect.asm")
		want := fileRead("../test/compare/Rect.hack")

		os.WriteFile("../test/output/Rect.hack", []byte(got), os.ModeAppend)

		if want != got {
			t.Errorf("got %q, want %q", got, want)
		}
	})
	
	t.Run("assembler Pong.asm", func(t *testing.T) {
		got := Assemble("../test/Pong.asm")
		want := fileRead("../test/compare/Pong.hack")

		os.WriteFile("../test/output/Pong.hack", []byte(got), os.ModeAppend)

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