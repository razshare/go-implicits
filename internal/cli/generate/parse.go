package generate

import (
	"errors"
	"fmt"
	"strings"

	"github.com/razshare/go-implicits/globals"
)

func Parse(source string, build Build) (text string, err error) {
	var find strings.Builder
	var repl strings.Builder

	next := func(line string) (Mod, error) {
		find.Reset()
		repl.Reset()
		state := Start
		for _, char := range line {
			switch state {
			case Start:
				switch char {
				case '"':
					state = ReadingOriginalString
				}
			case ReadingOriginalString:
				switch char {
				case '\\':
					state = EscapingOriginal
				case '"':
					state = DoneReadingOriginalString
				default:
					find.WriteRune(char)
				}
			case DoneReadingOriginalString:
				switch char {
				case ' ':
					state = ExpectingReplacementString
				}
			case ExpectingReplacementString:
				switch char {
				case '"':
					state = ReadingReplacementString
				default:
				}
			case ReadingReplacementString:
				switch char {
				case '\\':
					state = EscapingReplacement
				case '"':
					state = DoneReadingReplacementString
					return Mod{
						Pattern:     find.String(),
						Replacement: repl.String(),
					}, nil
				default:
					repl.WriteRune(char)
				}
			case EscapingOriginal:
				repl.WriteRune(char)
				state = ReadingOriginalString
				continue
			case EscapingReplacement:
				repl.WriteRune(char)
				state = ReadingReplacementString
				continue
			default:
				switch char {
				case ' ':
					// Noop.
					continue
				}
				state = Invalid
				return Mod{}, fmt.Errorf("invalid state, expecting \",\", received \"%c\" instead", char)
			}
		}
		return Mod{}, errors.New("invalid mod")
	}

	var mod Mod
	var sb strings.Builder

	modsGlobal := make([]Mod, 0)
	modsGlobalLen := 0
	modsLocal := make([]Mod, 0)
	modsLocalLen := 0

	for _, line := range strings.Split(source, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, globals.CodegenGlobalModHint) {
			offset := len(globals.CodegenGlobalModHint)
			if mod, err = next(line[offset:]); err != nil {
				return
			}
			modsGlobal = append(modsGlobal, mod)
			modsGlobalLen++
			continue
		} else if strings.HasPrefix(trimmed, globals.CodegenLineModHint) {
			o := len(globals.CodegenLineModHint)
			if mod, err = next(line[o:]); err != nil {
				return
			}
			modsLocal = append(modsLocal, mod)
			modsLocalLen++
			continue
		}

		if modsGlobalLen > 0 {
			if err = build(Block{Mods: modsGlobal, Line: &line}); err != nil {
				return
			}
		}

		if modsLocalLen > 0 {
			if err = build(Block{Mods: append(modsGlobal, modsLocal...), Line: &line}); err != nil {
				return
			}
			modsLocal = make([]Mod, 0)
			modsLocalLen = 0
		}

		sb.WriteString(line + "\n")
	}

	text = sb.String()

	return
}
