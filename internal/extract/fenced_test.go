package extract

import "testing"

func TestBlock(t *testing.T) {
	raw := "```sql\nSELECT 1\n```\n"
	got, err := Block(raw, "sql")
	if err != nil {
		t.Fatal(err)
	}
	if got != "SELECT 1" {
		t.Fatalf("got %q", got)
	}
}
