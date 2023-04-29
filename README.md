# SiteWatchdog

SiteWatchdog is a GitHub Action that checks the status of a list of websites and writes the results to a markdown file in your repository. The action is written in Go and uses a YAML file to specify the sites to check.

## Usage

To use SiteWatchdog, you need to add it as a step in your GitHub Actions workflow. Here's an example workflow that uses SiteWatchdog:

```yaml
name: SiteWatchdog
on: 
  push:
    branches:
      - main
  schedule:
    - cron:  '*/10 * * * *'
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@main

      - name: Run SiteWatchDog
        uses: qjoly/SiteWatchdog@main
        env:
          SITES_YAML_PATH: .github/sites.yaml
          README_TEMPLATE_PATH: .github/README.md.tmpl
        with:
          show_output: false

      - name: Push changes if needed
        run: |
          git add README.md
          git config user.email 41898282+github-actions[bot]@users.noreply.github.com
          git config user.name "GitHub Actions[Bot]"
          # Push only if README.md has been modified
          if [[ $(git status --porcelain "README.md" | wc -l) -gt 0 ]]; then git commit -m "[BOT] Update README.md"; git push; fi

```

This workflow runs on every push to the `main` branch and checks the status of the sites specified in the `sites.yaml` file. The `README.md` file is generated using the template in `README.md.tmpl`. The `show_output` input controls whether the results are printed to the console or not.

## Configuration

SiteWatchdog uses a YAML file to specify the sites to check. The file should contain a list of site objects, each with a `name` and a `url` field:

```yaml
sites:
- name: My website
  url: https://example.com
- name: My blog
  url: https://blog.example.com
```

By default, SiteWatchdog will look for `sites.yaml` and `README.md.tmpl` in the root of your repository. If you need to use a different path, you can set the `SITES_YAML_PATH` and `README_TEMPLATE_PATH` environment variables in your workflow.

You can use [README-example.md.tmpl](https://github.com/QJoly/SiteWatchdog/blob/main/README-example.md) as a starting point for your README.md file.

## Permissions

In order for SiteWatchdog to write the results to the `README.md` file, you need to grant it read and write permissions to your repository. To do this, you can add the following permissions to the `permissions` section of your workflow:

![Permissions](https://github.com/QJoly/SiteWatchdog/blob/main/.github/readwriteperm.png?raw=true)

This will allow SiteWatchdog to create and modify files in your repository.

## License

SiteWatchdog is released under the MIT License. See `LICENSE` for details.

# Example : 

| Website                 | Status                |
| ----------------------- | --------------------- |
| [Perdu](https://www.perdu.com) | :green_square: |
| [MarchePas](https://www.nemarchepascom) | :red_square: |
| [GitHub](https://www.github.com) | :green_square: |