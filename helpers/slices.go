package helpers

func SliceStringContains(s []string, search string) bool {
	set := make(map[string]struct{}, len(s))
	for _, s := range s {
		set[s] = struct{}{}
	}

	_, ok := set[search]
	return ok
}
