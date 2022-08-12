# Godog step adapter

This step adapter allows to switch from the [godog](https://github.com/cucumber/godog) step definition format:

```go
func shouldBeEqualDefault(arg1 int, arg2 int) error {
	if arg1 != arg2 {
		return errors.New("strings are not equal")
	}

	return nil
}

func TestFeatures(t *gotesting.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			s.Step(`^(\d+) equals to (\d+)$`, shouldBeEqualDefault)
		}
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
```

to the go testing format:
```go
import ("github.com/kyma-incubator/testdrape/godog/testing")

func shouldBeEqual(t *testing.T, arg1 int, arg2 int) {
	assert.Equal(t, arg1, arg2)
}

func TestFeaturesWrapped(t *gotesting.T) {
	tt := tester{}

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			testing.NewContext(s).Register(`^{int} equals to {int}$`, shouldBeEqual)
		}
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
```

## Predefined types
### int
```go

// Scenario: Integers comparison
//    Then 45 should be equal to 45

// ...
testing.NewContext(s).Register(`^{int} should be equal to {int}$`, shouldBeEqual)
// ...

func shouldBeEqual(t *testing.T, arg1, arg2 int) {
    assert.Equal(t, arg1, arg2)
}
```
### word
```go

// Scenario: Words comparison
//    Then one should be equal to one

// ...
testing.NewContext(s).Register(`^{word} should be equal to {word}$`, shouldBeEqual)
// ...

func shouldBeEqual(t *testing.T, arg1, arg2 string) {
    assert.Equal(t, arg1, arg2)
}
```
### string
```go

// Scenario: Words comparison
//    Then "another one" should be equal to 'another one'

// ...
testing.NewContext(s).Register(`^{string} should be equal to {string}$`, shouldBeEqual)
// ...

func shouldBeEqual(t *testing.T, arg1, arg2 string) {
    assert.Equal(t, arg1, arg2)
}
```
### Custom Tokens
```go

// Scenario: Custom tokens usage
//    Then there should be 25 pods

// ...
testing.NewContext(s).Int('Pods Count').Register(`^there should be {Pods Count} pods$`, shouldBeEqual)
// ...

func shouldBeEqual(t *testing.T, podsCount int) {
    assert.Equal(t, arg1, arg2)
}
```
