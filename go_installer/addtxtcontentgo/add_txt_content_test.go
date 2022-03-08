package addtxtcontentgo_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/flocknroll/azga/go_installer/addtxtcontentgo"
	. "github.com/flocknroll/azga/go_installer/addtxtcontentgo"
)

func TestAddContentFound(t *testing.T) {
	src, _ := ioutil.TempFile(os.TempDir(), "src*.txt")
	dest, _ := ioutil.TempFile(os.TempDir(), "dest*.txt")
	defer os.Remove(src.Name())
	defer os.Remove(dest.Name())

	io.WriteString(src, `
D
E
F`)

	io.WriteString(dest, `
A
B
C

D
E
F`)
	src.Close()
	dest.Close()

	if !CheckContent(src.Name(), dest.Name()) {
		t.Fail()
	}
}

func TestAddContentNotFound(t *testing.T) {
	src, _ := ioutil.TempFile(os.TempDir(), "src*.txt")
	dest, _ := ioutil.TempFile(os.TempDir(), "dest*.txt")
	defer os.Remove(src.Name())
	defer os.Remove(dest.Name())

	io.WriteString(src, `
X
Y
Z`)

	io.WriteString(dest, `
A
B
C

D
E
F`)
	src.Close()
	dest.Close()

	if CheckContent(src.Name(), dest.Name()) {
		t.Fail()
	}
}

func TestCheckDelimitedSection(t *testing.T) {
	src, _ := ioutil.TempFile(os.TempDir(), "src*.txt")
	defer os.Remove(src.Name())

	io.WriteString(src, `
# Start
A
B
C
# End`)
	src.Close()

	s, e, f := addtxtcontentgo.CheckDelimitedSection(src.Name(), "# Start", "# End")

	if !f || s != 2 || e != 6 {
		t.Fail()
	}
}
