package chunk

func SplitText(text string, size int) []string {

	var chunks []string

	for i := 0; i < len(text); i += size {

		end := i + size

		if end > len(text) {
			end = len(text)
		}

		chunks = append(chunks, text[i:end])
	}

	return chunks
}
