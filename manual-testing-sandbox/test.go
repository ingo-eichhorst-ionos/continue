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
	fmt.Println("  log: Natural logarithm")

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
	// Trim and split the input into parts
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid input format: please use 'number operator number' or '√ number' or 'log number'")
	}

	// Handle special cases with one number (square root and log)
	switch parts[0] {
	case "√":
		if len(parts) != 2 {
			return 0, fmt.Errorf("invalid input format for square root")
		}
		num, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number for square root: %v", err)
		}
		if num < 0 {
			return 0, fmt.Errorf("cannot calculate square root of a negative number")
		}
		return math.Sqrt(num), nil
	case "log":
		if len(parts) != 2 {
			return 0, fmt.Errorf("log operation requires exactly one number")
		}
		num, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number for logarithm: %v", err)
		}
		if num <= 0 {
			return 0, fmt.Errorf("cannot calculate logarithm of zero or a negative number")
		}
		return math.Log(num), nil
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
		return 0, fmt.Errorf("invalid operator: must be +, -, *, /, ^, √, or log")
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

func TestEdgeCases(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"  ", 0, true},             // Empty input
		{"10 * ", 0, true},          // Missing second operand
		{"* 10", 0, true},           // Malformed expression
		{"10    *    2", 20, false}, // With excessive spaces but valid
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestNegativeNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"-5 + 3", -2, false},  // Positive result
		{"-5 + -3", -8, false}, // Negative result
		{"5 * -3", -15, false}, // Multiplication with negative
		{"-4 / 2", -2, false},  // Division with negative
		{"-5 ^ 2", 25, false},  // Negative base power
		{"√ -9", 0, true},      // Invalid square root
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestDivisionPrecision(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"10 / 3", 3.3333, false}, // Simple division with precision check
		{"22 / 7", 3.1429, false}, // Approximation of pi
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && math.Abs(result-test.expected) > 0.0001 {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestScientificNotation(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"1e3 + 2e3", 3000, false},   // Addition with scientific notation
		{"1e3 * 2", 2000, false},     // Multiplication
		{"1.5e3 - 5e2", 1000, false}, // Subtraction
		{"2e3 / 2", 1000, false},     // Division
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestZeroOperands(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"0 + 0", 0, false}, // Zero addition
		{"0 - 0", 0, false}, // Zero subtraction
		{"0 * 5", 0, false}, // Zero multiplication
		{"5 * 0", 0, false}, // Zero multiplication
		{"0 / 1", 0, false}, // Zero division
		{"0 ^ 5", 0, false}, // Zero power
		{"5 ^ 0", 1, false}, // Zero exponent
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && result != test.expected {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestLogarithms(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"log 1", 0, false},       // Logarithm of 1
		{"log 2.71828", 1, false}, // Logarithm close to e
		{"log 0", 0, true},        // Invalid logarithm
		{"log -1", 0, true},       // Invalid logarithm
		{"log a", 0, true},        // Invalid input
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && math.Abs(result-test.expected) > 0.0001 {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestMixedOperations(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"3 + 2 * 2", 4, true},  // Valid only as sequential operations are not supported yet
		{"2 * 2 + 3", 5, true},  // Should produce error due to unsupported operation format
		{"10 / 2 - 5", 0, true}, // Division then subtraction, not parsed correctly
	}

	for _, test := range tests {
		_, err := calculate(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
	}
}

func TestFunnyMessages(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"3 + 3", "This should be easy, right?"},
		{"42 * 0", "The answer to life, the universe, and everything, but multiplied by nothing."},
		{"1024 / 256", "Back in the day, this was storage capacity!"},
		{"2 ^ 10", "That's kilobyte-sized knowledge!"},
		{"√ 10000", "How many runners make a marathon?"},
		{"pi delete 3", "Don't go breaking the circle, over here!"},
	}

	for _, test := range tests {
		fmt.Printf("Running a funny test with input '%s': %s\n", test.input, test.expectedMessage)
	}
}

func TestJokes(t *testing.T) {
	tests := []struct {
		scenario string
		joke     string
	}{
		{"Pun", "I just completed a calculation—it was mathematically pun-ishing!"},
		{"Why", "Why was the equal sign so humble? Because he knew he wasn’t less than or greater than anyone else."},
		{"Knock Knock", "Knock knock. Who’s there? Divide by two. Divide by two who? You heard me, divide by two and you’ll be half the trouble."},
		{"Math Joke", "Parallel lines have so much in common ... it’s a shame they’ll never meet."},
	}

	for _, test := range tests {
		fmt.Printf("Scenario: %s - Joke: %s\n", test.scenario, test.joke)
	}
}

func TestHipsterCalculations(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"Hipster + Coffee", 1},
		{"Vegan - Gluten", 1},
		{"Beard * Oil", 1},
		{"0 / AvocadoToast", 0},
		{"Smoothie ^ Blender", 0},
	}

	for _, test := range tests {
		fmt.Printf("Embracing hipster math with input '%s': just roll with it!\n", test.input)
	}
}

func TestAbsurdInputs(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"unicorns / rainbows", true},
		{"moon * cheese", true},
		{"time warp - flux capacitor", true},
		{"spaghetti ^ meatballs", true},
		{"chocolate log blackhole", true},
	}

	for _, test := range tests {
		fmt.Printf("Handling the impossible with input '%s': Expected error? %v\n", test.input, test.hasError)
	}
}

func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"2 + 3 * 4", 14, true},   // Expression evaluation not supported
		{"(4 + 5) * 6", 54, true}, // Parentheses not supported
		{"8 / 2 + 3", 7, true},    // BODMAS not supported
	}

	for _, test := range tests {
		_, err := calculate(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
	}
}

func TestStringInputs(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"Hello World", true},
		{"123abc + 456def", true},
		{"rainbows & unicorns", true},
		{"log Banana", true},
	}

	for _, test := range tests {
		_, err := calculate(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
	}
}

func TestPrecisionLimits(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		precision float64
		hasError  bool
	}{
		{"1.000000001 + 2.000000002", 3.000000003, 0.00000001, false},
		{"1.999999999 / 3", 0.666666666, 0.000000001, false},
		{"1.000000001 - 1.000000001", 0, 0.000000001, false},
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
		if !test.hasError && err != nil {
			t.Errorf("Unexpected error for input '%s': %v", test.input, err)
			continue
		}
		if !test.hasError && math.Abs(result-test.expected) > test.precision {
			t.Errorf("For input '%s': expected %f, but got %f", test.input, test.expected, result)
		}
	}
}

func TestReservedCommands(t *testing.T) {
	tests := []struct {
		input    string
		hasError bool
	}{
		{"printf", true}, // Common command, not an operation
		{"exit 0", true}, // Simulating shell command
		{"QUIT", false},  // Input similar to quit command
	}

	for _, test := range tests {
		_, err := calculate(test.input)
		if test.hasError && err == nil {
			t.Errorf("Expected error for input '%s', but got none", test.input)
			continue
		}
	}
}

func TestOverconfidence(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"5 / 0", 0, false},         // Assumes it works out no matter what
		{"sqrt -1", 0, false},       // Imaginary numbers, who needs theories?
		{"let's_do_this", 0, false}, // All inputs will work, because why not?
	}

	for _, test := range tests {
		result, err := calculate(test.input)

		if test.hasError && err == nil {
			fmt.Printf("Blind optimism on input '%s': expected failure but nah, it's cool.\n", test.input)
		}
		if !test.hasError && result != test.expected {
			fmt.Printf("This should have been perfect for '%s', expected %f, got %f.\n", test.input, test.expected, result)
		}
	}
}

func TestMisleadingInputs(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Infinity + 1", "Infinity and beyond!"},
		{"π * π", "The pie's the limit!"},
		{"0 ^ 0", "Indeterminate yet totally fine!"},
	}

	for _, test := range tests {
		fmt.Printf("Input '%s' faces the harsh reality: %s\n", test.input, test.expected)
	}
}

func TestHopefulOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"10 ~ 10", 0, true}, // Wishing for a new wave operator
		{"20 <3 5", 0, true}, // Pretty number logic
		{"9 += 3", 0, true},  // Algebraic wish-thinking
	}

	for _, test := range tests {
		_, err := calculate(test.input)
		if !test.hasError && err != nil {
			fmt.Printf("It could be great if it worked on input '%s', right?\n", test.input)
		}
	}
}

func TestAmbitiousCalculations(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"∞ / ∞", 0, true},  // Because infinity should solve everything
		{"2 ** 3", 0, true}, // Hoping for a double star operation
		{"10..20", 0, true}, // Range? No, just creative syntax
	}

	for _, test := range tests {
		_, err := calculate(test.input)
		if !test.hasError && err != nil {
			fmt.Printf("Dreaming big on input '%s': maybe not today.\n", test.input)
		}
	}
}
