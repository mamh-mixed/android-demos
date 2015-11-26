package settle

// import (
// 	"bufio"
// 	"io"
// 	"os"
// )

// func readSftpTxt(fn string) ([][]string, error) {
// 	f, err := os.Open(fn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	buf := bufio.NewReader(f)
// 	for {
// 		line, err := buf.ReadString('\n')
// 		line = strings.TrimSpace(line)
// 		if err != nil {
// 			if err == io.EOF {
// 				return nil
// 			}
// 			return err
// 		}
// 	}
// 	return nil
// }
