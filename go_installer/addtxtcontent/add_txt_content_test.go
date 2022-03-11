package addtxtcontent_test

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	. "github.com/flocknroll/azga/go_installer/addtxtcontent"
	"github.com/flocknroll/azga/go_installer/msfstools"
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

	if !CheckContent(src.Name(), dest.Name(), 1) {
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

	if CheckContent(src.Name(), dest.Name(), 1) {
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

	s, e, l, f := CheckDelimitedSection(src.Name(), "# Start", "# End")

	if !f || s != 2 || e != 6 || l != 6 {
		t.Fail()
	}
}

func BenchmarkCheckContent(b *testing.B) {
	src, _ := ioutil.TempFile(os.TempDir(), "src*.txt")
	defer os.Remove(src.Name())

	io.WriteString(src, `# AZGA DATA START
	BENIC,50.401443,2.790833,LF,0
	ADIBO,50.188332,2.819167,LF,0
	TWENT,50.545833,3.496944,LF,0
	ZEBLO,50.419723,2.969444,LF,0
	MPEPE,50.404167,3.275000,LF,0
	LYKII,50.342220,3.675000,LF,0
	CDVLL,50.312527,3.148083,LF,0
	LUDOV,49.912498,2.840000,LF,0
	KNAKI,49.514805,3.033277,LF,0
	DIDRL,49.268612,3.683056,LF,0
	YVOUN,49.808334,3.318056,LF,0
	AZGAR,50.006668,4.586666,LF,0
	LUKYN,50.354999,3.793333,LF,0
	MONFX,49.958610,4.093889,LF,0
	DITES,49.815834,3.767392,LF,0
	BENTR,49.569637,4.002861,LF,0
	PYRAX,49.596390,4.357500,LF,0
	IF33R,49.955841,3.631867,LF,1
	IF15L,50.249561,3.367461,LF,1
	IF33L,49.948059,3.633391,LF,1
	IF15R,50.249817,3.361817,LF,1
	# AZGA DATA END`)
	src.Close()

	path, _ := msfstools.GetPackageFolderPath()
	dest := filepath.Join(path, "Community/aerosoft-crj/Data/NavData/Waypoints.txt")
	CheckContent(src.Name(), dest, 1)
}
