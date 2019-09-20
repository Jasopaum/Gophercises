package primitive

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type Mode uint8

const (
	ModeCombo Mode = iota
	ModeTriangle
	ModeRect
	ModeEllipse
	ModeCircle
	ModeRotatedRect
	ModeBeziers
	ModeRotatedEllipse
	ModePolygon
)

func WithMode(mode Mode) func() []string {
	return func() []string {
		return []string{"-m", fmt.Sprintf("%d", mode)}
	}
}

func Transform(image io.Reader, numShapes int, opts ...func() []string) (io.Reader, error) {
	inTmpFilename := "tmp_in.*.png"
	outTmpFilename := "tmp_out.*.png"
	inTmpFile, err := ioutil.TempFile("/tmp/", inTmpFilename)
	if err != nil {
		return nil, errors.New("Could not create in tmp file.")
	}
	defer os.Remove(inTmpFile.Name()) // clean up
	outTmpFile, err := ioutil.TempFile("/tmp/", outTmpFilename)
	if err != nil {
		return nil, errors.New("Could not create out tmp file.")
	}
	defer os.Remove(outTmpFile.Name()) // clean up

	// Write what is in reader in a file
	if _, err := io.Copy(inTmpFile, image); err != nil {
		return nil, errors.New("Could not write image to tmp file.")
	}

	var args []string
	for _, opt := range opts {
		args = append(args, opt()...)
	}

	err = primitive(inTmpFile.Name(), outTmpFile.Name(), numShapes, args)
	if err != nil {
		return nil, err
	}

	b := bytes.NewBuffer(nil)
	_, err = io.Copy(b, outTmpFile)
	if err != nil {
		return nil, errors.New("Could not write result to buffer.")
	}
	return b, nil
}

func primitive(inputFile, outputFile string, numShapes int, otherArgs []string) error {
	fmt.Println("Creating image...")
	args := fmt.Sprintf("-i %s -o %s -n %d", inputFile, outputFile, numShapes)
	allArgs := append(strings.Fields(args), otherArgs...)
	cmd := exec.Command("primitive", allArgs...)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New("Could not create image from primitive.")
	}
	fmt.Printf("%s\n", stdoutStderr)
	return nil
}
