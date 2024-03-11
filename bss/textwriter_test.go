package bss

func TestTextWriter(t *testing.T) {

	params := TextFileParameters{}
	w, err := NewTextFileWriter("1.txt", params)
	if err != nil {
		t.Fatal(err)
	}

	var users []UserRecord

	for _, u := range users {
		if err := w.Append(u); err != nil {
			t.Fatal(err)
		}
	}

	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
}
