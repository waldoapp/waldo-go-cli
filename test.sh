#!/bin/bash

# export WALDO_UPLOAD_TOKEN=ffcdb81ee5119a500968f4b6c25e4106

waldo ~/.wdo/apps/SnoopSim/SnoopSim.app                     \
      --git_branch bogosity                                 \
      --git_commit 0123456789abcdef0123456789abcdef01234567 \
      --upload_token ffcdb81ee5119a500968f4b6c25e4106       \
      --variant_name bogus                                  \
      --verbose
