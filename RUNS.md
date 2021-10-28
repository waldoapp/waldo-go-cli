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

## CI triggered runs

Lastly, you can trigger a run of all your currently enabled test flows with a
direct call to the Waldo API. This is convenient when you want to trigger a run
via CI. Simply add the following to your CI script:

```bash
UPLOAD_TOKEN=0123456789abcdef0123456789abcdef

curl -X POST -H "Authorization: Upload-Token ${UPLOAD_TOKEN}" https://api.waldo.io/suites
```

> **Note:** You _must_ replace the fake upload token value shown above with the
> real value for your Waldo application.

A run of this kind is tagged as `ci-trigger`.

For further details about CI scripts and the Waldo upload token, please refer
to the `Documentation` section in the Waldo web app.
