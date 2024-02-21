package searcher

func getPath(doc map[string]any, parts []string) (any, bool) {
	var segment any = doc

	for _, part := range parts {
		m, ok := segment.(map[string]any)

		if !ok {
			return nil, false
		}

		if segment, ok = m[part]; !ok {
			return nil, false
		}
	}

	return segment, true
}
