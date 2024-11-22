package main

type Planner struct {
}

func NewPlanner() *Planner {
	return &Planner{}
}

func (p *Planner) distributeFiles(files []string, executors []string) (map[int][]string, []int) {
	filesPerExecutor := make(map[int][]string)
	usedExecutors := []int{}

	for i, file := range files {
		executorIdx := i % len(executors)
		filesPerExecutor[executorIdx] = append(filesPerExecutor[executorIdx], file)

		if len(filesPerExecutor[executorIdx]) == 1 {
			usedExecutors = append(usedExecutors, executorIdx)
		}
	}
	return filesPerExecutor, usedExecutors
}
