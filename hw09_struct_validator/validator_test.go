package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
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

	BadRegexp struct {
		Dummy string `validate:"regexp:["`
	}

	IntLength struct {
		Dummy int `validate:"len:34"`
	}

	StrangeStruct struct {
		Dummy string `validate:"contains:some"`
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
		{
			in:          Response{Code: 502},
			expectedErr: ValidationErrors{ValidationError{Field: "Code", Err: ErrInConstraint}},
		},
		{
			in:          Response{Code: 500},
			expectedErr: nil,
		},
		{
			in:          Token{},
			expectedErr: nil,
		},
		{
			in: User{},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLenConstraint},
				ValidationError{Field: "Age", Err: ErrMinConstraint},
				ValidationError{Field: "Email", Err: ErrRegexpConstraint},
				ValidationError{Field: "Role", Err: ErrInConstraint},
			},
		},
		{
			in: User{Phones: []string{""}},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLenConstraint},
				ValidationError{Field: "Age", Err: ErrMinConstraint},
				ValidationError{Field: "Email", Err: ErrRegexpConstraint},
				ValidationError{Field: "Role", Err: ErrInConstraint},
				ValidationError{Field: "Phones[0]", Err: ErrLenConstraint},
			},
		},
		{
			in: User{Phones: []string{"79054443377", "1234567"}},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: ErrLenConstraint},
				ValidationError{Field: "Age", Err: ErrMinConstraint},
				ValidationError{Field: "Email", Err: ErrRegexpConstraint},
				ValidationError{Field: "Role", Err: ErrInConstraint},
				ValidationError{Field: "Phones[1]", Err: ErrLenConstraint},
			},
		},
		{
			in: User{ID: "308f06f2-5dc2-461e-a029-360ce6b4b24f", Email: "sample@test.com"},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrMinConstraint},
				ValidationError{Field: "Role", Err: ErrInConstraint},
			},
		},
		{
			in: User{
				ID:    "308f06f2-5dc2-461e-a029-360ce6b4b24f",
				Email: "sample@test.com",
				Role:  "stuff",
				meta:  json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrMinConstraint},
			},
		},
		{
			in: User{ID: "308f06f2-5dc2-461e-a029-360ce6b4b24f", Email: "sample@test.com", Role: "stuff", Age: 123},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Age", Err: ErrMaxConstraint},
			},
		},
		{
			in:          User{ID: "308f06f2-5dc2-461e-a029-360ce6b4b24f", Email: "sample@test.com", Role: "stuff", Age: 34},
			expectedErr: nil,
		},
		{
			in:          BadRegexp{},
			expectedErr: ErrBadValidateParam,
		},
		{
			in:          "I am not a struct",
			expectedErr: ErrNotStruct,
		},
		{
			in:          IntLength{},
			expectedErr: ValidationErrors{ValidationError{Field: "Dummy", Err: ErrWrongArgumentType}},
		},
		{
			in:          StrangeStruct{Dummy: "fff"},
			expectedErr: ErrUnknownValidateType,
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
				var valid ValidationErrors
				if errors.As(tt.expectedErr, &valid) {
					var realErr ValidationErrors
					require.ErrorAs(t, err, &realErr)
					require.Equal(t, len(valid), len(realErr))
					for i, expected := range valid {
						require.Equalf(t, expected.Field, realErr[i].Field, "Field at pos %d must be equal", i)
						require.ErrorIsf(t, realErr[i].Err, expected.Err, "Error at pos %d must be of type %T", i, expected.Err)
					}
				} else {
					require.ErrorIs(t, err, tt.expectedErr)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}
