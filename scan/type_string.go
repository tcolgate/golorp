// Code generated by "stringer -type Type"; DO NOT EDIT

package scan

import "fmt"

const _Type_name = "EOFErrorNewlineCommentNumberAtomFunctorAtomSpecialAtomVariableUnboundLeftBrackRightBrackBarEmptyListLeftParenRightParenStopCommaSemiColon"

var _Type_index = [...]uint8{0, 3, 8, 15, 22, 28, 32, 43, 54, 62, 69, 78, 88, 91, 100, 109, 119, 123, 128, 137}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return fmt.Sprintf("Type(%d)", i)
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
