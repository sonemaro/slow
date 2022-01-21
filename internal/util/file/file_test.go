package file

//
//import (
//	"github.com/spf13/afero"
//	"testing"
//)
//
//var exampleContent = `
//2022-01-20 03:14:07.475 +0330 [971] LOG:  database system was shut down at 2022-01-20 02:21:28 +0330
//2022-01-20 03:14:07.481 +0330 [965] LOG:  database system is ready to accept connections
//2022-01-20 03:56:04.025 +0330 [5628] postgres@postgres ERROR:  function pq_sleep(integer) does not exist at character 8
//2022-01-20 03:56:04.025 +0330 [5628] postgres@postgres HINT:  No function matches the given name and argument types. You might need to add explicit type casts.
//2022-01-20 03:56:04.025 +0330 [5628] postgres@postgres STATEMENT:  select pq_sleep(2);
//2022-01-20 03:56:39.668 +0330 [5628] postgres@postgres ERROR:  function pq_sleep(integer) does not exist at character 8
//2022-01-20 03:56:39.668 +0330 [5628] postgres@postgres HINT:  No function matches the given name and argument types. You might need to add explicit type casts.
//2022-01-20 03:56:39.668 +0330 [5628] postgres@postgres STATEMENT:  select pq_sleep(2);
//2022-01-20 03:56:49.136 +0330 [5628] postgres@postgres LOG:  duration: 2016.449 ms  statement: select pg_sleep(2);
//2022-01-20 03:58:09.683 +0330 [5628] postgres@postgres LOG:  duration: 2002.358 ms  statement: select pg_sleep(2);
//2022-01-20 04:37:47.992 +0330 [11532] soroush@soroush FATAL:  role "soroush" does not exist
//2022-01-20 06:49:03.969 +0330 [965] LOG:  received fast shutdown request
//2022-01-20 06:49:03.972 +0330 [965] LOG:  aborting any active transactions
//2022-01-20 06:49:03.973 +0330 [965] LOG:  background worker "logical replication launcher" (PID 977) exited with exit code 1
//2022-01-20 06:49:03.973 +0330 [972] LOG:  shutting down
//2022-01-20 06:49:04.009 +0330 [965] LOG:  database system is shut down
//2022-01-20 00:39:37.004 +0330 [73383] postgres@postgres LOG:  duration: 10020.428 ms  statement: SELECT pg_sleep(10);
//2022-01-20 00:42:10.572 +0330 [73383] postgres@postgres LOG:  duration: 2002.590 ms  statement: SELECT pg_sleep(2);
//2022-01-20 00:48:52.624 +0330 [73383] postgres@postgres LOG:  duration: 2002.398 ms  statement: SELECT pg_sleep(2);
//2022-01-20 01:39:15.411 +0330 [73383] postgres@postgres LOG:  duration: 2002.336 ms  statement: SELECT pg_sleep(2);
//2022-01-20 01:43:38.882 +0330 [73383] postgres@postgres LOG:  duration: 2002.335 ms  statement: SELECT pg_sleep(2);
//2022-01-20 01:43:41.446 +0330 [73383] postgres@postgres LOG:  duration: 2002.332 ms  statement: SELECT pg_sleep(2);
//2022-01-20 02:21:27.993 +0330 [961] LOG:  received fast shutdown request
//2022-01-20 02:21:27.996 +0330 [961] LOG:  aborting any active transactions
//2022-01-20 02:21:27.998 +0330 [961] LOG:  background worker "logical replication launcher" (PID 976) exited with exit code 1
//2022-01-20 02:21:27.998 +0330 [971] LOG:  shutting down
//2022-01-20 02:21:28.028 +0330 [961] LOG:  database system is shut down
//`
//
//func TestGetLastLineWithSeek(t *testing.T) {
//	fs := afero.NewMemMapFs()
//	afs := &afero.Afero{Fs: fs}
//	f, err := afs.TempFile("/tmp", "test_file")
//	if err != nil {
//		t.Fatal(err)
//	}
//	_, err = f.WriteString(exampleContent)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	if got := GetLastLineWithSeek("/tmp/test_file"); got != tt.want {
//		t.Errorf("GetLastLineWithSeek() = %v, want %v", got, tt.want)
//	}
//}
