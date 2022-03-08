# Triggering a Test Flow Run

## Automatically triggering a run

Each time you upload a new build of your app to Waldo, a run of all your
currently enabled test flows is triggered _automatically_.

A run of this kind is tagged as `auto`.

## Manually triggering a run

You can also _manually_ trigger a run of all your currently enabled test flows
from the Waldo web app. Clicking on the `+` button drops down a menu containing
a `Run Tests` item for you to select.

A run of this kind is tagged as `manual`.

## CI-triggered runs

Lastly, you can trigger a run of of one or more test flows for your app with
Waldo CLI. (See [README.md][readme] for installation instructions.) This is
convenient when you want to trigger a run via CI. Simply add the following to
your CI script:

```bash
$ waldo trigger --upload_token 0123456789abcdef0123456789abcdef
```

> **Important:** Make sure you replace the fake upload token value shown above
> with the _real_ value for your Waldo app.

You can also use an environment variable to provide the upload token to Waldo
CLI:

```bash
$ export WALDO_UPLOAD_TOKEN=0123456789abcdef0123456789abcdef
$ waldo trigger
```

A run of this kind is tagged as `ci-trigger`.

### Advanced Usage

Whereas only the upload token is _required_ to successfully trigger a run on
Waldo, there are a couple other _non-required_ options recognized by Waldo CLI
that you may find useful:

- `--rule_name <value>` — This option instructs Waldo CLI to only run the test
  flows permitted by the named rule.
- `--verbose` — If you specify this option, Waldo CLI prints additional debug
  information. This can shed more light on why your trigger reqeuest is
  failing.

---

For further details about CI scripts and the Waldo upload token, please refer
to the `Documentation` section in the Waldo web app.

[readme]:   https://github.com/waldoapp/waldo-go-cli/blob/master/README.md
