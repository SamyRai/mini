package strategy

import (
	"fmt"
)

// ValidationContext manages validation strategies and execution
// This follows the Context pattern from Strategy pattern with full type safety
type ValidationContext[T any] struct {
	strategies map[string]ValidationStrategy[T]
}

// NewValidationContext creates a new validation context with type safety
func NewValidationContext[T any]() *ValidationContext[T] {
	return &ValidationContext[T]{
		strategies: make(map[string]ValidationStrategy[T]),
	}
}

// RegisterStrategy registers a validation strategy
func (c *ValidationContext[T]) RegisterStrategy(strategy ValidationStrategy[T]) {
	c.strategies[strategy.GetName()] = strategy
}

// Validate validates a value using the specified strategy
func (c *ValidationContext[T]) Validate(field string, value T, strategyName string) error {
	strategy, exists := c.strategies[strategyName]
	if !exists {
		return fmt.Errorf("validation strategy '%s' not found", strategyName)
	}

	return strategy.Validate(field, value)
}

// ValidateAll validates a value using all registered strategies
func (c *ValidationContext[T]) ValidateAll(field string, value T) error {
	for _, strategy := range c.strategies {
		if err := strategy.Validate(field, value); err != nil {
			return err
		}
	}
	return nil
}

// GetAvailableStrategies returns the names of available strategies
func (c *ValidationContext[T]) GetAvailableStrategies() []string {
	strategies := make([]string, 0, len(c.strategies))
	for name := range c.strategies {
		strategies = append(strategies, name)
	}
	return strategies
}

// HasStrategy checks if a strategy is registered
func (c *ValidationContext[T]) HasStrategy(strategyName string) bool {
	_, exists := c.strategies[strategyName]
	return exists
}

// GetStrategyCount returns the number of registered strategies
func (c *ValidationContext[T]) GetStrategyCount() int {
	return len(c.strategies)
}
