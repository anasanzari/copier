package copier_test

import (
	"errors"

	"reflect"
	"testing"
	"time"

	"github.com/anasanzari/copier"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Name     string
	Birthday *time.Time
	Nickname string
	Role     string
	Age      int32
	FakeAge  *int32
	Notes    []string
	flags    []byte
}

func (user User) DoubleAge() int32 {
	return 2 * user.Age
}

type Employee struct {
	Name      string
	Birthday  *time.Time
	Nickname  *string
	Age       int64
	FakeAge   int
	EmployeID int64
	DoubleAge int32
	SuperRule string
	Notes     []string
	flags     []byte
}

func (employee *Employee) Role(role string) {
	employee.SuperRule = "Super " + role
}

func checkEmployee(employee Employee, user User, t *testing.T, testCase string) {
	if employee.Name != user.Name {
		t.Errorf("%v: Name haven't been copied correctly.", testCase)
	}
	if employee.Nickname == nil || *employee.Nickname != user.Nickname {
		t.Errorf("%v: NickName haven't been copied correctly.", testCase)
	}
	if employee.Birthday == nil && user.Birthday != nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Birthday != nil && user.Birthday == nil {
		t.Errorf("%v: Birthday haven't been copied correctly.", testCase)
	}
	if employee.Age != int64(user.Age) {
		t.Errorf("%v: Age haven't been copied correctly.", testCase)
	}
	if user.FakeAge != nil && employee.FakeAge != int(*user.FakeAge) {
		t.Errorf("%v: FakeAge haven't been copied correctly.", testCase)
	}
	if employee.DoubleAge != user.DoubleAge() {
		t.Errorf("%v: Copy from method doesn't work", testCase)
	}
	if employee.SuperRule != "Super "+user.Role {
		t.Errorf("%v: Copy to method doesn't work", testCase)
	}
	if !reflect.DeepEqual(employee.Notes, user.Notes) {
		t.Errorf("%v: Copy from slice doen't work", testCase)
	}
}

func TestCopyStruct(t *testing.T) {
	var fakeAge int32 = 12
	user := User{Name: "Jinzhu", Nickname: "jinzhu", Age: 18, FakeAge: &fakeAge, Role: "Admin", Notes: []string{"hello world", "welcome"}, flags: []byte{'x'}}
	employee := Employee{}

	if err := copier.Copy(employee, &user); err == nil {
		t.Errorf("Copy to unaddressable value should get error")
	}

	copier.Copy(&employee, &user)
	checkEmployee(employee, user, t, "Copy From Ptr To Ptr")

	employee2 := Employee{}
	copier.Copy(&employee2, user)
	checkEmployee(employee2, user, t, "Copy From Struct To Ptr")

	employee3 := Employee{}
	ptrToUser := &user
	copier.Copy(&employee3, &ptrToUser)
	checkEmployee(employee3, user, t, "Copy From Double Ptr To Ptr")

	employee4 := &Employee{}
	copier.Copy(&employee4, user)
	checkEmployee(*employee4, user, t, "Copy From Ptr To Double Ptr")
}

func TestCopyFromStructToSlice(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}
	employees := []Employee{}

	if err := copier.Copy(employees, &user); err != nil && len(employees) != 0 {
		t.Errorf("Copy to unaddressable value should get error")
	}

	if copier.Copy(&employees, &user); len(employees) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(employees[0], user, t, "Copy From Struct To Slice Ptr")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, user); len(*employees2) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee((*employees2)[0], user, t, "Copy From Struct To Double Slice Ptr")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, user); len(employees3) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*(employees3[0]), user, t, "Copy From Struct To Ptr Slice Ptr")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, user); len(*employees4) != 1 {
		t.Errorf("Should only have one elem when copy struct to slice")
	} else {
		checkEmployee(*((*employees4)[0]), user, t, "Copy From Struct To Double Ptr Slice Ptr")
	}
}

func TestCopyFromSliceToSlice(t *testing.T) {
	users := []User{User{Name: "Jinzhu", Age: 18, Role: "Admin", Notes: []string{"hello world"}}, User{Name: "Jinzhu2", Age: 22, Role: "Dev", Notes: []string{"hello world", "hello"}}}
	employees := []Employee{}

	if copier.Copy(&employees, users); len(employees) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(employees[0], users[0], t, "Copy From Slice To Slice Ptr @ 1")
		checkEmployee(employees[1], users[1], t, "Copy From Slice To Slice Ptr @ 2")
	}

	employees2 := &[]Employee{}
	if copier.Copy(&employees2, &users); len(*employees2) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee((*employees2)[0], users[0], t, "Copy From Slice Ptr To Double Slice Ptr @ 1")
		checkEmployee((*employees2)[1], users[1], t, "Copy From Slice Ptr To Double Slice Ptr @ 2")
	}

	employees3 := []*Employee{}
	if copier.Copy(&employees3, users); len(employees3) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*(employees3[0]), users[0], t, "Copy From Slice To Ptr Slice Ptr @ 1")
		checkEmployee(*(employees3[1]), users[1], t, "Copy From Slice To Ptr Slice Ptr @ 2")
	}

	employees4 := &[]*Employee{}
	if copier.Copy(&employees4, users); len(*employees4) != 2 {
		t.Errorf("Should have two elems when copy slice to slice")
	} else {
		checkEmployee(*((*employees4)[0]), users[0], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 1")
		checkEmployee(*((*employees4)[1]), users[1], t, "Copy From Slice Ptr To Double Ptr Slice Ptr @ 2")
	}
}

func TestEmbedded(t *testing.T) {
	type Base struct {
		BaseField1 int
		BaseField2 int
	}

	type Embed struct {
		EmbedField1 int
		EmbedField2 int
		Base
	}

	base := Base{}
	embeded := Embed{}
	embeded.BaseField1 = 1
	embeded.BaseField2 = 2
	embeded.EmbedField1 = 3
	embeded.EmbedField2 = 4

	copier.Copy(&base, &embeded)

	if base.BaseField1 != 1 {
		t.Error("Embedded fields not copied")
	}
}

type structSameName1 struct {
	A string
	B int64
	C time.Time
}

type structSameName2 struct {
	A string
	B time.Time
	C int64
}

func TestCopyFieldsWithSameNameButDifferentTypes(t *testing.T) {
	obj1 := structSameName1{A: "123", B: 2, C: time.Now()}
	obj2 := &structSameName2{}
	err := copier.Copy(obj2, &obj1)
	if err != nil {
		t.Error("Should not raise error")
	}

	if obj2.A != obj1.A {
		t.Errorf("Field A should be copied")
	}
}

type ScannerValue struct {
	V int
}

func (s *ScannerValue) Scan(src interface{}) error {
	return errors.New("I failed")
}

type ScannerStruct struct {
	V *ScannerValue
}

type ScannerStructTo struct {
	V *ScannerValue
}

func TestScanner(t *testing.T) {
	s := &ScannerStruct{
		V: &ScannerValue{
			V: 12,
		},
	}

	s2 := &ScannerStructTo{}

	err := copier.Copy(s2, s)
	if err != nil {
		t.Error("Should not raise error")
	}

	if s.V.V != s2.V.V {
		t.Errorf("Field V should be copied")
	}
}

func TestNonNilCopy(t *testing.T) {
	type V struct {
		Value string
	}
	type S struct {
		Version int
		V []V
	}

	source := S{}
	source.V = append(source.V, V{
		Value: "Test",
	})

	dest := S{
		Version: 1,
	}
	err := copier.Copy(&dest, &source)
	if err != nil {
		t.Error("Should not raise error")
	}
	if len(dest.V) == 0 || dest.V[0].Value != "Test" {
		t.Error("Copy failed.")
	}
	if dest.Version != 1 {
		t.Error("Version failed.")
	}
}

func TestAllPossibleTypes(t *testing.T)  {
	type Str struct {
		V int
	}
	type V struct {
		Float float64
		Int int64
		IntPointer *int
		String string
		StringPointer *string
		IntArray []int64
		StringArray []string
		StrValue Str
		StrPointer *Str
	}
	dest := V{}
	dest.Float = 100
	dest.Int = 1
	v := 10
	dest.IntPointer = &v
	dest.String = "Man"
	s := "Wow"
	dest.StringPointer = &s
	dest.IntArray = []int64{1, 2, 3}
	dest.StringArray = []string{"a", "b", "c"}
	dest.StrValue = Str{ V: 10 }
	dest.StrPointer = &Str{ V: 11 }

	destCopy := dest


	source := V{}
	err := copier.Copy(&dest, &source)
	if err != nil {
		t.Error("Should not raise error")
	}

	// test a == b
	equalityTest := func(a V, b V) {
		assert.Equal(t, a.Float, b.Float)
		assert.Equal(t, a.Int, b.Int)
		assert.Equal(t, *a.IntPointer, *b.IntPointer)
		assert.Equal(t, a.String, b.String)
		assert.Equal(t, *a.StringPointer, *b.StringPointer, )
		assert.Equal(t, a.IntArray, b.IntArray)
		assert.Equal(t, a.StringArray, b.StringArray)
		assert.Equal(t, a.StrValue, b.StrValue)
		assert.Equal(t, *a.StrPointer, *b.StrPointer)
	}

	// now this copy shouldn't copy anything.
	equalityTest(destCopy, dest)

	source.Float = 200
	copier.Copy(&dest, &source)
	destCopy.Float = 200
	equalityTest(destCopy, dest)

	source.Int = 110
	copier.Copy(&dest, &source)
	destCopy.Int = 110
	equalityTest(destCopy, dest)

	vp := 99
	source.IntPointer = &vp
	copier.Copy(&dest, &source)
	cp := 99
	destCopy.IntPointer = &cp
	equalityTest(destCopy, dest)

	source.String = "Check"
	copier.Copy(&dest, &source)
	destCopy.String = "Check"
	equalityTest(destCopy, dest)

	source.IntArray = []int64{4, 5, 6}
	copier.Copy(&dest, &source)
	destCopy.IntArray = []int64{4, 5, 6}
	equalityTest(destCopy, dest)

	source.StringArray = []string{"d", "e", "f"}
	copier.Copy(&dest, &source)
	destCopy.StringArray = []string{"d", "e", "f"}
	equalityTest(destCopy, dest)


	//source.StrValue = Str{ V: 89 }
	//copier.Copy(&dest, &source)
	//destCopy.StrValue = Str{ V: 89 }
	//equalityTest(destCopy, dest)

	source.StrPointer = &Str{ V: 89 }
	copier.Copy(&dest, &source)
	destCopy.StrPointer = &Str{ V: 89 }
	equalityTest(destCopy, dest)

}