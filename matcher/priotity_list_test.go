package matcher

import (
	"fmt"
	"library-manager/logger"
	"math/rand"
	"testing"
)

type Node struct {
	Alias       string
	mExprLength int
}

func (n Node) Less(sortable ISortable) bool {
	var otherNode = sortable.(Node)
	return n.mExprLength < otherNode.mExprLength
}

func (n Node) Equal(sortable ISortable) bool {
	var otherNode = sortable.(Node)
	return n.mExprLength == otherNode.mExprLength
}

func createRandomNode() Node {
	var priority = rand.Intn(20)
	return Node{
		Alias:       fmt.Sprintf("Node:%d", priority),
		mExprLength: priority,
	}
}

func TestPriorityList_AddZ(t *testing.T) {
	var list PriorityList
	for i := 0; i < 10; i++ {
		var nodeToAdd = createRandomNode()
		logger.Info.Println(nodeToAdd)
		err := list.AddZ(nodeToAdd)
		if err != nil {
			logger.Error.Println(err)
		}
	}
	logger.Info.Printf("list:\n%v", list.mSlice)
}
