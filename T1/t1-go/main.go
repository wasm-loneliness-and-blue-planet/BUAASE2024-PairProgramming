/*
Function name: bocchiShutUp() or bocchi_shut_up() etc., choose the appropriate naming format according to your programming language of choice;
Parameters:
An i32 type numeric flag, 1 or 2;
an i32 type array seq, containing a number of two-digit numbers between [11,16] or [21,26];
an i32 type number size, is the number of elements of the array seq.
Return value: a number of type i32.
Behavior:
Check flag:
If flag is 1, count the frequency of occurrence of each digit in seq with ten digits of 1, e.g., 11 occurs 5 times, 12 occurs 4 times, etc.
If flag is 2, similarly, count the number of occurrences of each number in the seq with a tens digit of 2;
The most frequent number in the count is "👻 ghost":
If there is only one "👻 ghost", return that number;
If there is more than one "👻 ghost", return 10.

Translated with DeepL.com (free version)
*/

package main

// import (
//   "fmt"
// )

// Declare a main function, this is the entrypoint into our go module
// That will be run. In our example, we won't need this
func main() {
	/*
		array := [11]int{21, 21, 21, 21, 21, 21, 26, 26, 26, 26, 26}
		res := bocchiShutUp(2, array[:], 11)
		fmt.Println(res)
	*/
}

//export bocchiShutUp
func bocchiShutUp(flag int, seq []int, size int) int {
	var counts [6]int
	var offset int
	if flag == 1 {
		offset = 11
	} else if flag == 2 {
		offset = 21
	}
	for i := 0; i < size; i++ {
		number := seq[i]
		if number < offset || number > offset+5 {
			continue
		}
		counts[number-offset]++
	}

	var maxValue int = -1
	var maxIndex int = 0
	var multipleMax bool = false
	for i := 0; i < 6; i++ {
		if counts[i] > maxValue {
			multipleMax = false
			maxValue = counts[i]
			maxIndex = i
		} else if counts[i] == maxValue {
			multipleMax = true
		}
	}
	if multipleMax {
		return offset - 1
	}
	return maxIndex + offset
}
