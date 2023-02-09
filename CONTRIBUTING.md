Thank you for being interested in contributing to Rymdport.

There are various ways to contribute; everything is not about contributing code.

## Reporting a bug

If you've found an issue with the application, please report it to help us fix it as soon as possible.
When reporting a bug, please follow the guidelines below:

1. Check the [issue list](https://github.com/Jacalz/rymdport/issues) to see if it's already been reported. If so, update the existing issue with any additional information that you have.
2. If not, then create a new issue using the issue template for reporting bugs.
3. Stay involved in the conversation on the issue and answer any questions that might arise. More information can sometimes be necessary.

## Code Contribution

Great! You have either found a bug to fix or a new feature to implement.
Follow the steps below to increase the chance of the changes being accepted quickly.

1. Read and follow the guidelines in the [Code standards] (# Code-standards) section further down this page.
2. Consider how to structure your code so that it can be easily tested.
3. Write the code changes and create a new commit for your change.
4. Run the tests and make sure everything still works as expected using `go test ./...`.
5. Open a PR against the `main` branch. If there is an open bug, you should add "Fixes #", followed by the issue number, on a new line.
6. Please refrain from pushing or squashing.This makes it easier to review, and squashing can instead be done automatically when merging.

### Code standards

We aim to maintain a very high standard of code through design, testing, and implementation.
To manage this, we have various checks and processes in place that everyone should follow, including:

* For a more strict standard Go format, we use [gofumpt](https://github.com/mvdan/gofumpt).
* Imports should be ordered in accordance with the GoImports specification. Imports should be grouped by first- and third-party packages and listed alphabetically.
* The code should pass the code quality checks by [staticcheck](https://staticcheck.io/) and [gosec](https://github.com/securego/gosec).
