package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"testing"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Calculator started! Available operations:")
	fmt.Println("  +  : Addition")
	fmt.Println("  -  : Subtraction")
	fmt.Println("  *  : Multiplication")
	fmt.Println("  /  : Division")
	fmt.Println("  ^  : Power")
	fmt.Println("  √  : Square root (use only one number)")

	for {
		fmt.Print("\nEnter calculation (e.g., 5 + 3 or √ 16) or 'q' to quit: ")
		scanner.Scan()
		input := scanner.Text()

		if strings.ToLower(input) == "q" {
			fmt.Println("Goodbye!")
			break
		}

		result, err := calculate(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Result: %v\n", result)
	}
}

func calculate(input string) (float64, error) {
	// Split the input into parts
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid input format: please use 'number operator number' or '√ number'")
	}

	// Handle square root separately as it only needs one number
	if parts[0] == "√" {
		if len(parts) != 2 {
			return 0, fmt.Errorf("square root operation requires exactly one number")
		}
		num, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number for square root: %v", err)
		}
		if num < 0 {
			return 0, fmt.Errorf("cannot calculate square root of a negative number")
		}
		return math.Sqrt(num), nil
	}

	// For all other operations, we need exactly 3 parts
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid input format: please use 'number operator number'")
	}

	// Parse first number
	num1, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid first number: %v", err)
	}

	// Get operator
	operator := parts[1]

	// Parse second number
	num2, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid second number: %v", err)
	}

	// Perform calculation
	switch operator {
	case "+":
		return num1 + num2, nil
	case "-":
		return num1 - num2, nil
	case "*":
		return num1 * num2, nil
	case "/":
		if num2 == 0 {
			return 0, fmt.Errorf("division by zero is not allowed")
		}
		return num1 / num2, nil
	case "^":
		return math.Pow(num1, num2), nil
	default:
		return 0, fmt.Errorf("invalid operator: must be +, -, *, /, ^, or √")
	}
}

// Test functions
func TestCalculate(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"2 + 3", 5, false},
		{"10 - 5", 5, false},
		{"4 * 3", 12, false},
		{"15 / 3", 5, false},
		{"2 ^ 3", 8, false},
		{"√ 16", 4, false},
		{"2 + a", 0, true},
		{"15 / 0", 0, true},
		{"√ -4", 0, true},
		{"3 $ 4", 0, true},
		{"", 0, true},
		{"1", 0, true},
		{"1 +", 0, true},
		{"+ 1", 0, true},
		{"√ 2 3", 0, true},
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		// Check error cases
		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}

		// For non-error cases, check the result
		if !test.hasError && result != test.expected {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestSpecialCases(t *testing.T) {
	// Test precision for floating point operations
	result, err := calculate("2.5 + 2.5")
	if err != nil {
		t.Errorf("Unexpected error for floating point addition: %v", err)
	}
	if result != 5.0 {
		t.Errorf("Floating point addition failed: expected 5.0, got %f", result)
	}

	// Test power with decimal numbers
	result, err = calculate("2.5 ^ 2")
	if err != nil {
		t.Errorf("Unexpected error for power operation: %v", err)
	}
	expected := 6.25
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("Power operation failed: expected %f, got %f", expected, result)
	}

	// Test square root precision
	result, err = calculate("√ 2")
	if err != nil {
		t.Errorf("Unexpected error for square root: %v", err)
	}
	expected = math.Sqrt(2)
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("Square root operation failed: expected %f, got %f", expected, result)
	}
}

func TestLargeNumbers(t *testing.T) {
	// Test addition of large numbers
	result, err := calculate("9000000000 + 1000000000")
	if err != nil {
		t.Errorf("Unexpected error for large number addition: %v", err)
	}
	expected := 10000000000.0
	if result != expected {
		t.Errorf("Large number addition failed: expected %f, got %f", expected, result)
	}

	// Test multiplication of large numbers
	result, err = calculate("100000 * 50000")
	if err != nil {
		t.Errorf("Unexpected error for large number multiplication: %v", err)
	}
	expected = 5000000000.0
	if result != expected {
		t.Errorf("Large number multiplication failed: expected %f, got %f", expected, result)
	}

	// Test power operation with large numbers
	result, err = calculate("2 ^ 30")
	if err != nil {
		t.Errorf("Unexpected error for large power operation: %v", err)
	}
	expected = 1073741824.0
	if result != expected {
		t.Errorf("Large power operation failed: expected %f, got %f", expected, result)
	}
}
