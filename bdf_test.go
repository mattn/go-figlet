package figlet

import "testing"

func TestSynthesizeHalfwidthChar(t *testing.T) {
	src := bdfChar{
		encoding: 'オ',
		width:    8,
		bitmap:   "18003C007E00DB00FF00DB007E003C00",
	}

	got, ok := synthesizeHalfwidthChar(src, 'ｵ', 8)
	if !ok {
		t.Fatalf("synthesizeHalfwidthChar returned !ok")
	}
	if got.encoding != 'ｵ' {
		t.Fatalf("encoding = %q, want %q", got.encoding, 'ｵ')
	}
	if got.width != 4 {
		t.Fatalf("width = %d, want 4", got.width)
	}

	rows, ok := decodeBDFBitmap(got.bitmap, got.width, 8)
	if !ok {
		t.Fatalf("decodeBDFBitmap returned !ok")
	}
	left, right, ok := bitmapBounds(rows)
	if !ok {
		t.Fatalf("bitmapBounds returned !ok")
	}
	if left != 0 || right != 3 {
		t.Fatalf("bounds = (%d,%d), want (0,3)", left, right)
	}
}

func TestAppendSyntheticHalfwidthKatakana(t *testing.T) {
	chars := []bdfChar{
		{
			encoding: 'オ',
			width:    8,
			bitmap:   "18003C007E00DB00FF00DB007E003C00",
		},
	}

	got := appendSyntheticHalfwidthKatakana(chars, 8)
	found := false
	for _, ch := range got {
		if ch.encoding == 'ｵ' {
			found = true
			if ch.width != 4 {
				t.Fatalf("synthetic width = %d, want 4", ch.width)
			}
		}
	}
	if !found {
		t.Fatalf("synthetic halfwidth katakana not added")
	}
}
