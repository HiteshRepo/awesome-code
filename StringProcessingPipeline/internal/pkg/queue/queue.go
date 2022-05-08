package queue

type Queue struct {
	data []string
}

func NewQueue() *Queue {
	return &Queue{data: make([]string, 0)}
}

func (q *Queue) Push(datum string) {
	q.data = append(q.data, datum)
}

func (q *Queue) Pop() string {
	str := q.data[0]
	q.data = q.data[1:]
	return str
}

func (q *Queue) IsEmpty() bool {
	return len(q.data) == 0
}
