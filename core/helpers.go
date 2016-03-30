package core

func uniqStrs(in []string) (out []string) {
	tab := make(map[string]struct{}) // empty struct because it's the smallest possible value
	for _, str := range in {
		if _, ok := tab[str]; !ok {
			tab[str] = struct{}{}
			out = append(out, str)
		}
	}
	return out
}
