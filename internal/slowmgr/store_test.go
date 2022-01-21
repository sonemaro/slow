package slowmgr

import (
	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Paginate(t *testing.T) {
	fakeData, err := generateFake(5)
	if err != nil {
		assert.Fail(t, "cannot generate fake data", err)
	}
	p := Paginate(fakeData, 2, 3)
	assert.Len(t, p, 2, "expected: %d | got: %d", 2, len(p))
	assert.EqualValues(t, p[0], fakeData[3], "objects are not equal. Obj1:%v | Obj2:%v", p[0], fakeData[3])

}

func generateFake(n int) ([]*Slow, error) {
	var fakeData []*Slow
	for i := 0; i < n; i++ {
		s := Slow{}
		err := faker.FakeData(&s)
		if err != nil {
			return nil, err
		}
		fakeData = append(fakeData, &s)
	}
	return fakeData, nil
}
