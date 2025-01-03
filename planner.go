package main

type Planner struct {
}

func NewPlanner() *Planner {
	return &Planner{}
}

func (p *Planner) distributeFiles(files []string, executors []int) (map[int][]string, []int) {
	filesPerExecutor := make(map[int][]string)
	usedExecutors := []int{}

	for i, file := range files {
		idx := i % len(executors)
		filesPerExecutor[executors[idx]] = append(filesPerExecutor[executors[idx]], file)

		if len(filesPerExecutor[executors[idx]]) == 1 {
			usedExecutors = append(usedExecutors, executors[idx])
		}
	}
	return filesPerExecutor, usedExecutors
}
