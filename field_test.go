package cronfab

import "testing"

func TestCrontabField_String(t *testing.T) {
	cf := CrontabField([][3]int{})
	s := cf.String()
	if s != "" {
		t.Errorf("unexpected value")
	}
	cf = CrontabField([][3]int{{1, 2, 1}})
	s = cf.String()
	if s != "1-2/1" {
		t.Errorf("unexpected value")
	}
	cf = CrontabField([][3]int{{1, 2, 1}, {3, 20, 1}})
	s = cf.String()
	if s != "1-2/1,3-20/1" {
		t.Errorf("unexpected value")
	}
}

func TestCrontabField_Len(t *testing.T) {
	cf := CrontabField([][3]int{})
	l := cf.Len()
	if l != 0 {
		t.Errorf("unexpected value")
	}
	cf = CrontabField([][3]int{{1, 2, 1}})
	l = cf.Len()
	if l != 1 {
		t.Errorf("unexpected value")
	}
	cf = CrontabField([][3]int{{1, 2, 1}, {3, 20, 1}})
	l = cf.Len()
	if l != 2 {
		t.Errorf("unexpected value")
	}
}

func TestCrontabField_SetConstraint(t *testing.T) {
	cf := CrontabField([][3]int{{1, 2, 1}})
	if cf.GetConstraint(0) != [3]int{1, 2, 1} {
		t.Errorf("unexpected value")
	}
	cf.SetConstraint(0, [3]int{3, 20, 2})
	if cf.GetConstraint(0) != [3]int{3, 20, 2} {
		t.Errorf("unexpected value")
	}
}

func TestCrontabField_Validate(t *testing.T) {
	cf := CrontabField([][3]int{{1, 2, 1}, {3, 4, 1}})
	err := cf.Validate()
	if err != nil {
		t.Errorf("err: %v", err)
	}
	cf = CrontabField([][3]int{{1, 2, 1}, {2, 4, 1}})
	err = cf.Validate()
	if err == nil {
		t.Errorf("unexpected value")
	}
}
