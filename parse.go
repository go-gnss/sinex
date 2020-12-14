package sinex

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// TODO: What should this struct contain?
type File struct {
	Version  float32             // Should header fields be struct attributes?
	Comments []string            // TODO: Comments with line numbers?
	Blocks   map[string][]string // What should this actually be?
}

func Parse(r io.Reader) (file File, err error) {
	br := bufio.NewReader(r)

	file = File{Comments: []string{}, Blocks: map[string][]string{}}

	// First line must be header
	if err = parseHeader(br, &file); err != nil {
		return file, err
	}

	line, err := readLine(br, &file)
	for ; err == nil; line, err = readLine(br, &file) {
		switch line[0] {
		case '+':
			if err = parseBlock(line[1:], br, &file); err != nil {
				break
			}
		case '%':
			// %ENDSNX<EOF> - Not sure why the LF is present
			if line != "%ENDSNX\n" {
				err = fmt.Errorf("invalid trailer line")
			}
			// TODO: Also check if the next read yields io.EOF?
			return file, err
		}
	}

	// TODO: Check if mandatory blocks are found
	return file, err
}

func readLine(br *bufio.Reader, file *File) (line string, err error) {
	line, err = br.ReadString('\n')
	if err != nil {
		return line, err
	}

	if line[0] == '*' {
		file.Comments = append(file.Comments, line[1:])
		return readLine(br, file)
	}

	return line, err
}

// TODO: Return header type or something
// "%=SNX " + F4.2 + " " + A3 + " " + I2:I3:I5 + I2:I3:I5 + " " + A1 + " " + I5 + " " + A1 + 6(" " + A1)
func parseHeader(br *bufio.Reader, file *File) error {
	header, err := br.ReadString('\n')
	if err != nil {
		return err
	}

	fields := strings.Split(header, " ")
	if len(fields) < 10 {
		return fmt.Errorf("invalid header line")
	}

	if fields[0] != "%=SNX" {
		return fmt.Errorf("invalid header line")
	}

	version, err := strconv.ParseFloat(fields[1], 32)
	file.Version = float32(version)

	return err
}

// The following blocks are defined:
//  FILE/REFERENCE
//  FILE/COMMENT
//  INPUT/HISTORY
//  INPUT/FILES
//  INPUT/ACKNOWLEDGEMENTS
//  NUTATION/DATA
//  PRECESSION/DATA
//  SOURCE/ID
//  SITE/ID
//  SITE/DATA
//  SITE/RECEIVER
//  SITE/ANTENNA
//  SITE/GPS_PHASE_CENTER
//  SITE/GAL_PHASE_CENTER
//  SITE/ECCENTRICITY
//  SATELLITE/ID
//  SATELLITE/PHASE_CENTER
//  BIAS/EPOCHS
//  SOLUTION/EPOCHS
//  SOLUTION/STATISTICS
//  SOLUTION/ESTIMATE
//  SOLUTION/APRIORI
//  SOLUTION/MATRIX_ESTIMATE {p} {type}
//  SOLUTION/MATRIX_APRIORI {p} {type}
//  SOLUTION/NORMAL_EQUATION_VECTOR
//  SOLUTION/NORMAL_EQUATION_MATRIX {p}
//    Where: {p} L or U
//           {type} CORR or COVA or INFO
func parseBlock(block string, br *bufio.Reader, file *File) error {
	// TODO: Check if block type is valid
	file.Blocks[block] = []string{} // TODO: Check if block redeclared

	for {
		line, err := readLine(br, file)
		if err != nil {
			return err // Wrap error in some context
		}

		if line == "-"+block {
			return nil
		}

		file.Blocks[block] = append(file.Blocks[block], line)
	}
}
