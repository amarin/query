package query

// UniqFieldNames returns unique FieldName set from supplied FieldName slice.
func UniqFieldNames(a []FieldName) []FieldName {
	answer := make([]FieldName, 0, len(a))
	for _, str := range a {
		found := false
		for j := 0; j < len(answer); j++ {
			if answer[j] == str {
				found = true
				break
			}
		}
		if !found {
			answer = append(answer, str)
		}
	}
	return answer
}
