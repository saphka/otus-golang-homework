package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in:          App{Version: "1234"},
			expectedErr: ValidationErrors{ValidationError{Field: "Version", Err: ErrLenConstraint}},
		},
		{
			in:          App{Version: "12345"},
			expectedErr: nil,
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			t.Parallel()

			err := Validate(tt.in)
			if tt.expectedErr != nil {
				if valid, ok := tt.expectedErr.(ValidationErrors); ok {
					var realErr ValidationErrors
					require.ErrorAs(t, err, &realErr)
					require.Equal(t, len(valid), len(realErr))
					for i, expected := range valid {
						require.Equalf(t, expected.Field, realErr[i].Field, "Field at pos %d must be equal", i)
						require.ErrorAsf(t, realErr[i].Err, &expected.Err, "Error at pos %d must be of type %T", i, expected.Err)
					}
				} else {
					require.ErrorAs(t, err, &tt.expectedErr)
				}
			} else {
				require.NoError(t, err)
			}

		})
	}
}
