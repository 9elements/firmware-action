# Testing with pytest
To run pytest and all tests, firstly `cd` into this directory with `cd .dagger-ci`.

Then just run:
```
$ pytest
```

## Coverage
Pytest can report test coverage, simply use following arguments `--cov --cov-report=term-missing`.

## Run slow test
Some tests are time-consuming, for example actual docker container building.

To run slow tests, simply use following argument `--runslow`.

## Live log
To see live log from tests, for example docker-building with dagger, simply use following arguments `--log-cli-level NOTSET  --show-capture no`.

To ommit the summary at the end, use argument `--no-summary`.

## Log verbosity
By default it is `DEBUG`, you can reduce it by using argument `--log-cli-level=INFO`.

## HTML report
You can also get more info with HTML report by using argument `--cov-report=html` and then opening `firmware-action/.dagger-ci/htmlcov/index.html` in browser.

## Recap
```
$ pytest --cov --cov-report=term-missing --runslow --log-cli-level NOTSET --show-capture no --log-cli-level=INFO --cov-report=html
```

