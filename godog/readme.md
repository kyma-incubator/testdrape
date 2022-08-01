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
			testing.NewContext(s).Register(`^(\d+) equals to (\d+)$`, shouldBeEqual)
		}
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
```
