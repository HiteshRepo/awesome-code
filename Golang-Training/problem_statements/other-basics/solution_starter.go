package main

func main() {
	solution := &SolutionImpl{}
	Problem(solution)
}

type SolutionImpl struct{}

func (s *SolutionImpl) GetDivideFunction() func(...int) int {
	return func(nums ...int) int {
		return 0
	}
}