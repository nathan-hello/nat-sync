// package tests
//
// import (
//
//	"fmt"
//	"io"
//	"testing"
//	"time"
//
// )
//
//	func TestConnection(t *testing.T) {
//		reader, writer := io.Pipe()
//		initAll(reader)
//		b, err := writer.Write([]byte("seek --hours=10"))
//		if err != nil {
//			t.Fatal(err)
//		}
//		fmt.Printf("bits written: %d\n", b)
//		time.Sleep(5 * time.Second)
//		reader.Close()
//	}
package tests
