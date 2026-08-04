package helpers

import "testing"

func TestGenerateNumericOTP(t *testing.T) {
	code, err := GenerateNumericOTP(6)
	if err != nil {
		t.Fatalf("GenerateNumericOTP returned error: %v", err)
	}
	if len(code) != 6 {
		t.Fatalf("expected 6-digit code, got %q (len %d)", code, len(code))
	}
	for _, c := range code {
		if c < '0' || c > '9' {
			t.Fatalf("expected all-numeric code, got %q", code)
		}
	}

	other, err := GenerateNumericOTP(6)
	if err != nil {
		t.Fatalf("GenerateNumericOTP returned error: %v", err)
	}
	if code == other {
		t.Fatalf("expected two generated codes to differ, both were %q", code)
	}
}
