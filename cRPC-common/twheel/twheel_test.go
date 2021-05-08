package twheel

import (
	"testing"
	"time"
)

var timeout = 5

func TestNew(t *testing.T) {

	ch := make(chan interface{}, 20)
	go func() {
		i := 0
		for msg := range ch {
			i++
			t.Logf("i = %d, %v", i, msg)
		}
	}()

	_, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
}

func testNew(ch chan interface{}) (*TimeWheel, error) {

	return New(20, ch)
}

type testIns struct {
}

func TestTimeWheel_Add(t *testing.T) {

	var addTime time.Time
	ch := make(chan interface{}, 20)
	go func() {
		i := 0
		for msg := range ch {
			i++
			t.Logf("i = %d, %v", i, msg)
			diff := int(time.Now().Sub(addTime).Seconds())
			if diff != timeout {
				t.Error("timeout was wrong")
			}
		}
	}()

	tw, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
	tw.Start()
	err = tw.Add(new(testIns), 5)
	if err != nil {
		t.Fatalf("add ins err %v", err)
	}
	addTime = time.Now()
	time.Sleep(10 * time.Second)
}

func TestTimeWheel_Refresh(t *testing.T) {

	ch := make(chan interface{}, 20)
	go func() {
		for range ch {
			t.Error("err")
		}
	}()

	tw, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
	tw.Start()
	ins := new(testIns)
	err = tw.Add(ins, 6)
	if err != nil {
		t.Fatalf("add ins err %v", err)
	}
	go func() {
		ticker := time.Tick(4 * time.Second)
		for {
			<-ticker
			err = tw.Refresh(ins, 6)
			if err != nil {
				t.Errorf("refresh ins err %v", err)
			}
			t.Log("refresh")
		}
	}()
	select {}
}

func TestTimeWheel_Delete(t *testing.T) {

	ch := make(chan interface{}, 20)
	go func() {
		for range ch {
			t.Error("timeout was wrong")
		}
	}()

	tw, err := testNew(ch)
	if err != nil {
		t.Fatalf("create tw err %v", err)
	}
	tw.Start()
	ins := new(testIns)
	err = tw.Add(ins, 5)
	if err != nil {
		t.Fatalf("add ins err %v", err)
	}
	tw.Delete(ins)
	time.Sleep(10 * time.Second)
}
