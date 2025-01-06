package parser

import "testing"


func TestParser(t *testing.T) {

    t.Run("test initialize function", func(t *testing.T) {

        _, err := Initialize("../test/Add.asm")

        if err != nil {
            t.Errorf("run errror")
        }
    })
}

func TestHasMoreCommands(t *testing.T) {
    t.Run("open multiple line file", func(t *testing.T) {
        inputData, err := Initialize("../test/Add.asm")

        if err != nil {
            t.Errorf("cannot open file")
        }

        got := inputData.HasMoreCommands()
        want := true

        if got != want {
            t.Errorf("expected %t but got %t", want, got)
        }
    })

    t.Run("open empty file", func(t *testing.T) {
        inputData, err := Initialize("../test/None.asm")

        if err != nil {
            t.Errorf("cannot open file")
        }

        got := inputData.HasMoreCommands()
        want := false

        if got != want {
            t.Errorf("expected %t but got %t", want, got)
        }
    })

    t.Run("open one line file", func(t *testing.T) {
        inputData, err := Initialize("../test/OneLine.asm")

        if err != nil {
            t.Errorf("cannot open file")
        }

        got := inputData.HasMoreCommands()
        want := true

        if got != want {
            t.Errorf("expected %t but got %t", want, got)
        }
    })
}

func TestAdvance(t *testing.T) {
    t.Run("")    
}