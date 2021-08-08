package gitlab

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ############################################################################
// # PendingTask DiffersFrom tests
// ############################################################################
func Test_PendngTask_Compare_Nil(t *testing.T) {
	var one *PendingTask

	assert.False(t, one.DiffersFrom(nil))
}

func Test_PendngTask_Compare_Empty(t *testing.T) {
	one := &PendingTask{}
	another := &PendingTask{}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_Nil_And_Empty(t *testing.T) {
	var one *PendingTask
	another := &PendingTask{}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_Empty_and_Nil(t *testing.T) {
	one := &PendingTask{}

	assert.True(t, one.DiffersFrom(nil))
}

func Test_PendngTask_Compare_Author_Equals(t *testing.T) {
	one := &PendingTask{AuthorUsername: "user"}
	another := &PendingTask{AuthorUsername: "user"}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_Author_Differs(t *testing.T) {
	one := &PendingTask{AuthorUsername: "user1"}
	another := &PendingTask{AuthorUsername: "user2"}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_LastCommentAuthor_Equals(t *testing.T) {
	one := &PendingTask{LastCommentatorUsername: "user"}
	another := &PendingTask{LastCommentatorUsername: "user"}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_LastCommentAuthor_Differs(t *testing.T) {
	one := &PendingTask{LastCommentatorUsername: "user1"}
	another := &PendingTask{LastCommentatorUsername: "user2"}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_ID_Equals(t *testing.T) {
	one := &PendingTask{ID: "123"}
	another := &PendingTask{ID: "123"}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_ID_Differs(t *testing.T) {
	one := &PendingTask{ID: "123"}
	another := &PendingTask{ID: "456"}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_ProjectID_Equals(t *testing.T) {
	one := &PendingTask{ProjectID: "123"}
	another := &PendingTask{ProjectID: "123"}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_ProjectID_Differs(t *testing.T) {
	one := &PendingTask{ProjectID: "123"}
	another := &PendingTask{ProjectID: "456"}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_WebURL_Equals(t *testing.T) {
	one := &PendingTask{WebURL: "123"}
	another := &PendingTask{WebURL: "123"}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_WebURL_Differs(t *testing.T) {
	one := &PendingTask{WebURL: "123"}
	another := &PendingTask{WebURL: "456"}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_IID_Equals(t *testing.T) {
	one := &PendingTask{IID: "123"}
	another := &PendingTask{IID: "123"}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_IID_Differs(t *testing.T) {
	one := &PendingTask{IID: "123"}
	another := &PendingTask{IID: "456"}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_UserNotesCount_Equals(t *testing.T) {
	one := &PendingTask{UserNotesCount: 123}
	another := &PendingTask{UserNotesCount: 123}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_UserNotesCount_Differs(t *testing.T) {
	one := &PendingTask{UserNotesCount: 123}
	another := &PendingTask{UserNotesCount: 456}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_LastUpdated_Equals(t *testing.T) {
	pointInTime := time.Now()
	one := &PendingTask{LastUpdated: pointInTime}
	another := &PendingTask{LastUpdated: pointInTime}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTask_Compare_LastUpdated_Differs(t *testing.T) {
	pointInTime := time.Now()
	one := &PendingTask{LastUpdated: pointInTime}
	another := &PendingTask{LastUpdated: pointInTime.Add(time.Second)}

	assert.True(t, one.DiffersFrom(another))
}

// ############################################################################
// # PendingTaskList DiffersFrom tests
// ############################################################################
func Test_PendngTaskList_Compare_Nil(t *testing.T) {
	var one PendingTaskList

	assert.False(t, one.DiffersFrom(nil))
}

func Test_PendngTaskList_Compare_Empty(t *testing.T) {
	one := PendingTaskList{}
	another := PendingTaskList{}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTaskList_Compare_Nil_And_Empty(t *testing.T) {
	var one PendingTaskList
	another := PendingTaskList{}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTaskList_Compare_Empty_And_Nil(t *testing.T) {
	one := PendingTaskList{}

	assert.False(t, one.DiffersFrom(nil))
}

func Test_PendngTaskList_Compare_Two_And_One(t *testing.T) {
	one := PendingTaskList{
		PendingTask{ID: "123"},
		PendingTask{ID: "123"},
	}
	another := PendingTaskList{
		PendingTask{ID: "123"},
	}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTaskList_Compare_One_And_Two(t *testing.T) {
	one := PendingTaskList{
		PendingTask{ID: "123"},
	}
	another := PendingTaskList{
		PendingTask{ID: "123"},
		PendingTask{ID: "123"},
	}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTaskList_Compare_Two_And_Two_Equal(t *testing.T) {
	one := PendingTaskList{
		PendingTask{ID: "123"},
		PendingTask{ID: "456"},
	}
	another := PendingTaskList{
		PendingTask{ID: "123"},
		PendingTask{ID: "456"},
	}

	assert.False(t, one.DiffersFrom(another))
}

func Test_PendngTaskList_Compare_Two_And_Two_LastPairDiffers(t *testing.T) {
	one := PendingTaskList{
		PendingTask{ID: "123"},
		PendingTask{ID: "456"},
	}
	another := PendingTaskList{
		PendingTask{ID: "123"},
		PendingTask{ID: "4567"},
	}

	assert.True(t, one.DiffersFrom(another))
}

func Test_PendngTaskList_Compare_One_And_One_Equal(t *testing.T) {
	one := PendingTaskList{
		PendingTask{ID: "123"},
	}
	another := PendingTaskList{
		PendingTask{ID: "123"},
	}

	assert.False(t, one.DiffersFrom(another))
}
func Test_PendngTaskList_Compare_One_And_One_Differs(t *testing.T) {
	one := PendingTaskList{
		PendingTask{ID: "123"},
	}
	another := PendingTaskList{
		PendingTask{ID: "456"},
	}

	assert.True(t, one.DiffersFrom(another))
}
