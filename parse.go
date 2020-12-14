package sinex

import (
	"bufio"
	"fmt"
	"io"
)

// TODO: What should this struct contain?
type File struct {
	Version  float32             // Should header fields be struct attributes?
	Comments []string            // TODO: Comments with line numbers?
	Blocks   map[string][]string // What should this actually be?
}

func Parse(r io.Reader) (file File, err error) {
	br := bufio.NewReader(r)
	if err = parseHeader(br); err != nil {
		// First line must be header
		return file, err
	}

	file = File{Comments: []string{}, Blocks: map[string][]string{}}

	char, err := br.ReadByte()
	for ; err == nil; char, err = br.ReadByte() {
		switch char {
		case '*':
			// TODO: Store comment line
		case '+':
			if err = parseBlock(br, &file); err != nil {
				break
			}
		case '%':
			// %ENDSNX<EOF>
			line, _, err := br.ReadLine()
			if err == nil && string(line) != "ENDSNX" {
				err = fmt.Errorf("invalid trailer line")
			}
			// TODO: Also check if the next read yields io.EOF?
			return file, err
		}
	}

	// TODO: Check if mandatory blocks are found
	return file, err
}

// TODO: Return header type or something
// "%=SNX " + F4.2 + " " + A3 + " " + I2:I3:I5 + I2:I3:I5 + " " + A1 + " " + I5 + " " + A1 + 6(" " + A1)
func parseHeader(br *bufio.Reader) error {
	header, err := br.ReadString('\n')
	if err != nil {
		return err
	}

	// TODO: Actually parse header line
	if header[:5] != "%=SNX" {
		return fmt.Errorf("invalid header line")
	}

	return nil
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
func parseBlock(br *bufio.Reader, file *File) error {
	block, err := br.ReadString('\n')
	if err != nil {
		return err
	}

	// TODO: Check if block type is valid
	file.Blocks[block] = []string{} // TODO: Check if block redeclared

	for {
		line, err := parseLine(br, file)
		if err != nil {
			return err // Wrap error in some context
		}

		if line == "-"+block {
			return nil
		}

		file.Blocks[block] = append(file.Blocks[block], line)
	}
}

func parseLine(br *bufio.Reader, file *File) (line string, err error) {
	line, err = br.ReadString('\n')
	if err != nil {
		return line, err
	}

	if line[0] == '*' {
		file.Comments = append(file.Comments, line[1:])
		return parseLine(br, file)
	}

	return line, err
}
