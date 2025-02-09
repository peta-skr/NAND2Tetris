package main

import (
	"log"
	"os"
	"testing"
)

// ジャンプ命令を実装する際にランダム文字列を使用しているため、テストが失敗すると思う
// 失敗したファイルだけDiffで確認する

func TestVM(t *testing.T) {
	t.Run("Basic Test", func(t *testing.T) {
		vm("./test/BasicTest/BasicTest.vm", "./test/BasicTest/BasicTest.asm")
		got, err := os.ReadFile("./test/BasicTest/BasicTest.asm")

		if err != nil {
			log.Fatal(err)
		}

		want, err := os.ReadFile("./test/BasicTest/BasicTest_comp.asm")

		if err != nil {
			log.Fatal(err)
		}

		if string(got) != string(want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Pointer Test", func(t *testing.T) {
		vm("./test/PointerTest/PointerTest.vm", "./test/PointerTest/PointerTest.asm")
		got, err := os.ReadFile("./test/PointerTest/PointerTest.asm")

		if err != nil {
			log.Fatal(err)
		}

		want, err := os.ReadFile("./test/PointerTest/PointerTest_comp.asm")

		if err != nil {
			log.Fatal(err)
		}

		if string(got) != string(want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Simple Add Test", func(t *testing.T) {
		vm("./test/SimpleAdd/SimpleAdd.vm", "./test/SimpleAdd/SimpleAdd.asm")
		got, err := os.ReadFile("./test/SimpleAdd/SimpleAdd.asm")

		if err != nil {
			log.Fatal(err)
		}

		want, err := os.ReadFile("./test/SimpleAdd/SimpleAdd_comp.asm")

		if err != nil {
			log.Fatal(err)
		}

		if string(got) != string(want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Stack Test", func(t *testing.T) {
		vm("./test/StackTest/StackTest.vm", "./test/StackTest/StackTest.asm")
		got, err := os.ReadFile("./test/StackTest/StackTest.asm")

		if err != nil {
			log.Fatal(err)
		}

		want, err := os.ReadFile("./test/StackTest/StackTest_comp.asm")

		if err != nil {
			log.Fatal(err)
		}

		if string(got) != string(want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("Static Test", func(t *testing.T) {
		vm("./test/StaticTest/StaticTest.vm", "./test/StaticTest/StaticTest.asm")
		got, err := os.ReadFile("./test/StaticTest/StaticTest.asm")

		if err != nil {
			log.Fatal(err)
		}

		want, err := os.ReadFile("./test/StaticTest/StaticTest_comp.asm")

		if err != nil {
			log.Fatal(err)
		}

		if string(got) != string(want) {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
