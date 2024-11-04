# cmt

`cmt` (short for commit) is a command-line tool designed to help developers write consistent and standardized commit messages across different commit conventions.

![demo](demo.gif)

## Conventions

Open source contributors often face a unique challenge when working across multiple projects: navigating the diverse landscape of commit conventions. Each repository may adhere to its own set of rules for structuring commit messages, from [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) to [Gitmoji](https://gitmoji.dev/), or even custom formats. This inconsistency can lead to confusion, reduced productivity, and potential errors in version control.

`cmt` aims to alleviate this pain point by providing a flexible, easy-to-use tool that adapts to different commit conventions. By standardizing the commit process across projects, `cmt` helps contributors maintain consistency and focus on what truly matters - their code contributions.

Supported conventions at the time of writing:

- [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
- [Gitmoji](https://gitmoji.dev/)

## Install

```bash
# Install in the current directory
curl -sSL https://raw.githubusercontent.com/segersniels/cmt/master/scripts/install.sh | bash
# Install in /usr/local/bin
curl -sSL https://raw.githubusercontent.com/segersniels/cmt/master/scripts/install.sh | sudo bash -s /usr/local/bin
```

## Usage

### New project

First initialize `cmt` in your project:

```
cmt init
```

This will create a `.cmtrc.json` file with your preferred settings.

After that simple create a commit:

```
cmt commit
```

or use the shorthand:

```
cmt c
```

Follow the interactive prompts to construct your commit.

### Existing project

You can choose to use `cmt` in an existing project without adding a new configuration file.
If no `.cmtrc.json` file is found, `cmt` will attempt to determine the commit convention from the last commit message.

## Configuration

`cmt` uses a `.cmtrc.json` file to store configuration. You can edit this file manually or run `cmt init` to set up your preferences.
