package fastcopy

import (
	"testing"

	"github.com/abhinav/tmux-fastcopy/internal/ui"
	tcell "github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestGenerateHints(t *testing.T) {
	t.Parallel()

	alphabet := []rune("abc")

	tests := []struct {
		desc    string
		text    string
		matches []Range
		want    []hint
	}{
		{
			desc: "no matches",
			text: "foo",
			want: []hint{},
		},
		{
			desc: "single match",
			text: "foo bar",
			matches: []Range{
				{1, 3}, // f(oo)
			},
			want: []hint{
				{
					Label: "a",
					Text:  "oo",
					Matches: []Range{
						{1, 3}, // f(oo)
					},
				},
			},
		},
		{
			desc: "duplicated match",
			text: "foo bar baz qux",
			matches: []Range{
				{4, 6},  // (ba)r
				{8, 10}, // (ba)z
			},
			want: []hint{
				{
					Label: "a",
					Text:  "ba",
					Matches: []Range{
						{4, 6},  // (ba)r
						{8, 10}, // (ba)z
					},
				},
			},
		},
		{
			desc: "multiple matches",
			text: "foo bar baz qux",
			matches: []Range{
				{0, 3},   // (foo)
				{4, 6},   // (ba)r
				{8, 10},  // (ba)z
				{13, 15}, // q(ux)
			},
			want: []hint{
				{
					Label: "c",
					Text:  "ba",
					Matches: []Range{
						{4, 6},  // (ba)r
						{8, 10}, // (ba)z
					},
				},
				{
					Label: "a",
					Text:  "foo",
					Matches: []Range{
						{0, 3}, // (foo)
					},
				},
				{
					Label: "b",
					Text:  "ux",
					Matches: []Range{
						{13, 15}, // q(ux)
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			got := generateHints(alphabet, tt.text, tt.matches)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestHintAnnotations(t *testing.T) {
	t.Parallel()

	style := Style{
		Normal:         tcell.StyleDefault,
		Match:          tcell.StyleDefault.Foreground(tcell.ColorGreen),
		SkippedMatch:   tcell.StyleDefault.Foreground(tcell.ColorGray),
		HintLabel:      tcell.StyleDefault.Foreground(tcell.ColorRed),
		HintLabelInput: tcell.StyleDefault.Foreground(tcell.ColorYellow),
	}

	tests := []struct {
		desc  string
		give  hint
		input string
		want  []ui.TextAnnotation
	}{
		{
			desc: "multiple matches",
			give: hint{
				Label: "a",
				Text:  "foo",
				Matches: []Range{
					{0, 3},
					{7, 10},
				},
			},
			// [a]oo
			want: []ui.TextAnnotation{
				ui.OverlayTextAnnotation{
					Offset:  0,
					Overlay: "a",
					Style:   style.HintLabel,
				},
				ui.StyleTextAnnotation{
					Offset: 1,
					Length: 2,
					Style:  style.Match,
				},
				ui.OverlayTextAnnotation{
					Offset:  7,
					Overlay: "a",
					Style:   style.HintLabel,
				},
				ui.StyleTextAnnotation{
					Offset: 8,
					Length: 2,
					Style:  style.Match,
				},
			},
		},
		{
			desc: "full input match",
			give: hint{
				Label: "a",
				Text:  "foo",
				Matches: []Range{
					{0, 3},
				},
			},
			input: "a",
			want: []ui.TextAnnotation{
				ui.OverlayTextAnnotation{
					Offset:  0,
					Overlay: "a",
					Style:   style.HintLabelInput,
				},
				ui.StyleTextAnnotation{
					Offset: 1,
					Length: 2,
					Style:  style.Match,
				},
			},
		},
		{
			desc: "multi character label",
			give: hint{
				Label: "ab",
				Text:  "foobar",
				Matches: []Range{
					{1, 7},
				},
			},
			want: []ui.TextAnnotation{
				ui.OverlayTextAnnotation{
					Offset:  1,
					Overlay: "ab",
					Style:   style.HintLabel,
				},
				ui.StyleTextAnnotation{
					Offset: 3,
					Length: 4,
					Style:  style.Match,
				},
			},
		},
		{
			desc: "multi character label/input match",
			give: hint{
				Label: "ab",
				Text:  "foobar",
				Matches: []Range{
					{1, 7},
				},
			},
			input: "a",
			want: []ui.TextAnnotation{
				ui.OverlayTextAnnotation{
					Offset:  1,
					Overlay: "a",
					Style:   style.HintLabelInput,
				},
				ui.OverlayTextAnnotation{
					Offset:  2,
					Overlay: "b",
					Style:   style.HintLabel,
				},
				ui.StyleTextAnnotation{
					Offset: 3,
					Length: 4,
					Style:  style.Match,
				},
			},
		},
		{
			desc: "multi character label/input mismatch",
			give: hint{
				Label: "ab",
				Text:  "foobar",
				Matches: []Range{
					{1, 7},
				},
			},
			input: "x",
			want: []ui.TextAnnotation{
				ui.StyleTextAnnotation{
					Offset: 1,
					Length: 6,
					Style:  style.SkippedMatch,
				},
			},
		},
		{
			desc: "long label",
			give: hint{
				Label: "abcd",
				Text:  "foo",
				Matches: []Range{
					{0, 3},
				},
			},
			want: []ui.TextAnnotation{
				ui.OverlayTextAnnotation{
					Offset:  0,
					Overlay: "abcd",
					Style:   style.HintLabel,
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()

			got := tt.give.Annotations(tt.input, style)
			assert.Equal(t, tt.want, got)
		})
	}
}
