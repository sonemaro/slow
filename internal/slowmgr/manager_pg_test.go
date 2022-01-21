package slowmgr

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

var examplePgLog = `
2022-01-21 04:37:43.526 +0330 [986] LOG:  database system was shut down at 2022-01-21 03:00:58 +0330
2022-01-21 04:37:43.536 +0330 [984] LOG:  database system is ready to accept connections
2022-01-21 08:01:27.020 +0330 [53986] postgres@postgres LOG:  duration: 5025.394 ms  statement: select pg_sleep(5);
2022-01-21 08:01:32.326 +0330 [53986] postgres@postgres ERROR:  syntax error at or near "pg_sleep" at character 8
2022-01-21 08:01:32.326 +0330 [53986] postgres@postgres STATEMENT:  insert pg_sleep(3);
2022-01-21 08:03:40.091 +0330 [53986] postgres@postgres ERROR:  syntax error at or near "pg_sleep" at character 8
2022-01-21 08:03:40.091 +0330 [53986] postgres@postgres STATEMENT:  insert pg_sleep(2);
2022-01-21 08:03:48.625 +0330 [53986] postgres@postgres LOG:  duration: 2002.342 ms  statement: select pg_sleep(2);
2022-01-21 08:04:00.018 +0330 [53986] postgres@postgres LOG:  duration: 4004.359 ms  statement: select pg_sleep(4);
2022-01-21 08:05:29.514 +0330 [53986] postgres@postgres LOG:  duration: 1207.829 ms  statement: select pg_sleep(1.2);
`

var file afero.File

func setupManager() (*ManagerPg, error) {
	fs := afero.NewOsFs()
	afs := &afero.Afero{Fs: fs}
	f, err := afs.TempFile("", "log")
	//defer os.Remove(f.Name())
	if err != nil {
		return nil, err
	}
	file = f
	//f.WriteString()

	stor := NewStoreDefault(DefaultStoreOptions{
		MaxEntries: 1000,
	})
	parsr := NewParserPg()
	mgr := NewManagerPg(f.Name(), ManagerPgOptions{
		Store:  stor,
		Parser: parsr,
	})
	return mgr, nil
}

func start() (*ManagerPg, error) {
	mgr, err := setupManager()
	if err != nil {
		return nil, err
	}
	err = mgr.Start()
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

func fakeSlow() *Slow {
	fkSl := func() *Slow {
		return &Slow{
			Query:     "SELECT * FROM users WHERE id = 47",
			Operation: "select",
			Duration:  2002.12,
		}
	}
	sl := Slow{}
	err := faker.FakeData(&sl)
	if err != nil {
		return fkSl()
	}
	return &sl
}

func TestManagerPg_Start(t *testing.T) {
	mgr, err := start()
	if err != nil {
		t.Fatal(err)
	}

	assert.True(t, mgr.started)
}

func TestManagerPg_Save(t *testing.T) {
	mgr, err := start()
	if err != nil {
		t.Fatal(mgr)
	}
	sl := fakeSlow()
	sl.Operation = "select"
	err = mgr.Save(sl)
	if err != nil {
		t.Fatal(err)
	}
	rows := mgr.Filter("select", 1, 1)
	assert.NotZero(t, len(rows))
	assert.EqualValues(t, rows[0], sl)
}

func TestManagerPg_Sort(t *testing.T) {
	mgr, err := start()
	if err != nil {
		t.Fatal(mgr)
	}
	var slows [10]*Slow
	for i := 0; i < 10; i++ {
		slows[i] = fakeSlow()
		err := mgr.Save(slows[i])
		if err != nil {
			t.Fatal(err)
		}
	}
	srt := mgr.Sort(1, 10)
	for i, sorted := range srt {
		for _, sj := range srt[i+1:] {
			assert.LessOrEqualf(t,
				sorted.Duration,
				sj.Duration,
				fmt.Sprintf("i:%d | 1st:%v | 2nd:%v", i, sorted, sj),
			)
		}
	}
}

func TestManagerPg_Filter(t *testing.T) {
	mgr, err := start()
	if err != nil {
		t.Fatal(mgr)
	}
	sl := fakeSlow()
	sl.Operation = "select"
	err = mgr.Save(sl)
	if err != nil {
		t.Fatal(err)
	}
	rows := mgr.Filter("select", 1, 1)
	assert.NotZero(t, len(rows))
	assert.EqualValues(t, rows[0], sl)
}

func TestAll(t *testing.T) {
	mgr, err := start()
	fmt.Println(file.Name())
	if err != nil {
		t.Fatal(mgr)
	}

	for _, l := range strings.Split(examplePgLog, "\n") {
		if l == "" {
			continue
		}
		_, err := file.WriteString(l + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}
	time.Sleep(1 * time.Second)
	data := mgr.Filter("select", 1, 10)
	assert.Len(t, data, 4, "invalid length")
}
