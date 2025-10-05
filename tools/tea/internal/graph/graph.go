// Package graph provides simple directed graph utilities such as cycle detection.
package graph

// FindCycles returns a list of cycles discovered via DFS in a directed graph.
func FindCycles(adj map[string][]string) [][]string {
	var cycles [][]string
	visited := make(map[string]bool)
	stack := make(map[string]bool)
	path := make([]string, 0)

	var dfs func(string)
	dfs = func(u string) {
		visited[u] = true
		stack[u] = true
		path = append(path, u)

		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v)
			} else if stack[v] {
				// Found a cycle; extract subpath starting from v
				var cyc []string
				for i := len(path) - 1; i >= 0; i-- {
					cyc = append([]string{path[i]}, cyc...)
					if path[i] == v {
						break
					}
				}
				if len(cyc) > 0 {
					cycles = append(cycles, cyc)
				}
			}
		}

		stack[u] = false
		path = path[:len(path)-1]
	}

	for u := range adj {
		if !visited[u] {
			dfs(u)
		}
	}
	return cycles
}
