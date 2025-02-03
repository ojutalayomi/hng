package utils

import (
	"io"
	"math"
	"net/http"
)

func IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}

	// Check only up to sqrt(n)
	for i := 5; i <= int(math.Sqrt(float64(n))); i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func IsPerfect(n int) bool {
	if n < 2 {
		return false
	}

	sum := 1 // 1 is always a divisor
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			sum += i
			if i != n/i {
				sum += n / i
			}
		}
	}

	return sum == n
}

func IsArmstrong(n int) bool {
	temp, sum, digits := n, 0, 0

	// Count number of digits
	for temp > 0 {
		digits++
		temp /= 10
	}

	temp = n // Reset temp

	// Compute sum of each digit raised to the power of digits
	for temp > 0 {
		digit := temp % 10
		sum += int(math.Pow(float64(digit), float64(digits)))
		temp /= 10
	}

	return sum == n
}

func IsEven(n int) bool {
	return n%2 == 0
}

func DigitalSum(n int) int {
	sum := 0
	n = Abs(n) // Handle negative numbers

	for n > 0 {
		sum += n % 10 // Extract last digit and add to sum
		n /= 10       // Remove last digit
	}

	return sum
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func FetchAPI(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() // Ensure response body is closed

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
