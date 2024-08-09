This is an application that builds quizzes. Can be used in conferences or other events.

It's still work in progress but already usable.

![screenshot](https://github.com/jimmykarily/quizmaker/blob/master/images/screenshot.png?raw=true)

All you need it to create a yaml file with the possible questions (see the test
file as an example: [test questions.yaml](tests/assets/question_pool.yaml))

Then you need to generate a secret that will sign the cookies. E.g. with:

```bash
export QUIZMAKER_COOKIE_SECRET=$(openssl rand -base64 32)
```

then run the application with golang:

```bash
go run . -question-pool questions.yaml
```

NOTE: This application started as part of the [Kairos.io](https://kairos.io/) team hackweek.

TODO:

- improve the README
- make it configurable so other teams can use their own logo and text
- create an easy way to collect results
- create an easy deployment method (kustomization / helm chart / other)
