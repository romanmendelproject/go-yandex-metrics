package utils

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseURLUpdate(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    URLParams
		wantErr bool
	}{
		{
			name: "Good url gauge",
			args: args{
				url: "/update/gauge/test/0.1",
			},
			want: URLParams{
				MetricType:  "gauge",
				MetricName:  "test",
				MetricValue: "0.1",
			},
			wantErr: false,
		},
		{
			name: "Good url count",
			args: args{
				url: "/update/gauge/test/1",
			},
			want: URLParams{
				MetricType:  "gauge",
				MetricName:  "test",
				MetricValue: "1",
			},
			wantErr: false,
		},
		{
			name: "Bad (too many arguments)",
			args: args{
				url: "/update/gauge/test/test/1",
			},
			want:    URLParams{},
			wantErr: true,
		},
		{
			name: "Bad (too few arguments)",
			args: args{
				url: "/update/gauge/1",
			},
			want:    URLParams{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURLUpdate(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseURLValue(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		want    URLParams
		wantErr bool
	}{
		{
			name: "Good url gauge",
			args: args{
				url: "/update/gauge/test",
			},
			want: URLParams{
				MetricType: "gauge",
				MetricName: "test",
			},
			wantErr: false,
		},
		{
			name: "Good url count",
			args: args{
				url: "/update/gauge/test",
			},
			want: URLParams{
				MetricType: "gauge",
				MetricName: "test",
			},
			wantErr: false,
		},
		{
			name: "Bad (too many arguments)",
			args: args{
				url: "/update/gauge/test/test/1",
			},
			want:    URLParams{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURLValue(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURLValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseURLValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFloatPtr(t *testing.T) {
	value := 5.1
	address := &value
	valuePtr := GetFloatPtr(value)

	assert.Equal(t, address, valuePtr)
}

// TestISinTrustedNetwork tests the ISinTrustedNetwork function.
func TestISinTrustedNetwork(t *testing.T) {
	tests := []struct {
		name     string
		checkIP  string
		cidr     string
		expected bool
	}{
		{
			name:     "IP in CIDR range",
			checkIP:  "192.168.1.10",
			cidr:     "192.168.1.0/24",
			expected: true,
		},
		{
			name:     "IP not in CIDR range",
			checkIP:  "192.168.2.10",
			cidr:     "192.168.1.0/24",
			expected: false,
		},
		{
			name:     "Invalid CIDR format",
			checkIP:  "192.168.1.10",
			cidr:     "192.168.1.0/33", // Invalid CIDR
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ISinTrustedNetwork(tt.checkIP, tt.cidr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestStringToInt tests the StringToInt function.
func TestStringToInt(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{
			name:     "Valid positive integer",
			input:    "123",
			expected: 123,
		},
		{
			name:     "Valid negative integer",
			input:    "-456",
			expected: -456,
		},
		{
			name:     "Zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "Invalid string",
			input:    "abc",
			expected: 0, // Expecting 0 on failure
		},
		{
			name:     "Empty string",
			input:    "",
			expected: 0, // Expecting 0 on failure
		},
		{
			name:     "Valid large integer",
			input:    "2147483647", // Max int32 value
			expected: 2147483647,
		},
		{
			name:     "Valid small integer",
			input:    "-2147483648", // Min int32 value
			expected: -2147483648,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StringToInt(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestUnPointer tests the UnPointer function for both int64 and float64 types.
func TestUnPointer(t *testing.T) {
	// Test cases for int64
	t.Run("int64 tests", func(t *testing.T) {
		var intVal int64 = 42
		var nilIntVal *int64 = nil
		tests := []struct {
			name     string
			input    *int64
			expected int64
		}{
			{"Non-nil pointer", &intVal, 42},
			{"Nil pointer", nilIntVal, 0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := UnPointer(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
	// Test cases for float64
	t.Run("float64 tests", func(t *testing.T) {
		var floatVal = 3.14
		var nilFloatVal *float64 = nil
		tests := []struct {
			name     string
			input    *float64
			expected float64
		}{
			{"Non-nil pointer", &floatVal, 3.14},
			{"Nil pointer", nilFloatVal, 0.0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := UnPointer(tt.input)
				assert.Equal(t, tt.expected, result)
			})
		}
	})
}

// TestToPointer tests the ToPointer function for both int64 and float64 types.
func TestToPointer(t *testing.T) {
	// Test cases for int64
	t.Run("int64 tests", func(t *testing.T) {
		var intVal int64 = 42
		result := ToPointer(intVal)
		assert.NotNil(t, result)         // Ensure the result is not nil
		assert.Equal(t, &intVal, result) // Ensure the pointer points to the original value
	})
	// Test cases for float64
	t.Run("float64 tests", func(t *testing.T) {
		var floatVal = 3.14
		result := ToPointer(floatVal)
		assert.NotNil(t, result)           // Ensure the result is not nil
		assert.Equal(t, &floatVal, result) // Ensure the pointer points to the original value
	})
}
