/* parse files with format:
Key = value 

Or 
Key = value value \
      value
*/
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func myScannerSplit(data []byte, atEOF bool) (advance int, token []byte, err error) {
   isComment := false
   for i := 0; i < len(data); i++ {
      if ((data[i] == '=') || (!isComment && data[i] == '\n') )  {
         return i + 1, data[:i], nil
      }
      if ((data[i] == '\\') || (data[i] == '#')) { // ignore until end of line
         isComment = true
      }
      if (data[i] == '\n') { // ignore until end of line
         data[i] = ' '
         isComment = false
      }
      if isComment {
         data[i] = ' '
      }
   }
   return 0, data, bufio.ErrFinalToken
}

func haveDuplicateKeys(values []string) []string {
   var keysDuplicate []string
   for i, v1 := range values {
      for _, v2 := range values [i+1:] {
         if v1 == v2 {
            keysDuplicate = append(keysDuplicate, v1)
         }
      }
   }
   return keysDuplicate
}

func main() {
   var keys       []string
   var lastKey    string
   var cleanValue string
   var file_str    string
   
   keyValues := make(map[string][]string)
   isKey := true
   fileName := "binarios.dep"
   
   if len(os.Args) > 1 {
      fileName = os.Args[1]      
   } 
   fmt.Printf("Processing file (%s)\n", fileName)
   f, err := os.Open(fileName)
   check(err)

   scanner := bufio.NewScanner(f)
   myBuffer := make([]byte, 500000)
   scanner.Buffer(myBuffer, 500000)
   scanner.Split(myScannerSplit)
   
   i := 0 
	for scanner.Scan() {
      i = i + 1
      file_str = ""
      file_str = strings.Trim(scanner.Text(), "\n\r\t ")
      if isKey && (file_str != "") {

         keys = append(keys, file_str)
         lastKey = file_str
         isKey = false
      } else if !isKey {
         keyVal := strings.Split(file_str, " ")
         for _, values := range keyVal {
            cleanValue = ""
            cleanValue = strings.Trim(values, "\t ")
            if len(cleanValue) > 0 {
               keyValues[cleanValue] = append(keyValues[cleanValue], lastKey)
            }
         }      
         isKey = true
      }
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}   
   f.Close()
   fmt.Println(keys) // Println will add back the final '\n'
   
   for k, v := range keyValues { 
      fmt.Printf("key[%s] value%s\n", k, v) 
   }
   
   for k, v := range keyValues { 
      keysDuplicate := haveDuplicateKeys(v)
      if len(keysDuplicate) > 0 {
         fmt.Printf("key with duplicates: [%s] value%s\n", k, keysDuplicate) 
      }
   }   
}