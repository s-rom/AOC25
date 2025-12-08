package main

type Circuit struct {
	id       int
	quantity int
	index    int
}

type MaxQueue []*Circuit

func (mq MaxQueue) Len() int {
	return len(mq)
}

func (mq MaxQueue) Less(i, j int) bool {
	return mq[i].quantity > mq[j].quantity
}

func (mq MaxQueue) Swap(i, j int) {
	mq[i], mq[j] = mq[j], mq[i]
	mq[i].index = i
	mq[j].index = j
}

func (mq *MaxQueue) Push(x any) {
	n := len(*mq)
	item := x.(*Circuit)
	item.index = n
	*mq = append(*mq, item)
}

func (mq *MaxQueue) Pop() any {
	old := *mq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*mq = old[0 : n-1]
	return item
}
