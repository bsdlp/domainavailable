# domainavailable

this is a quick script a wrote that looks for domains that are available for registration. It takes in a `names.csv`, `suffix.csv`, and `prefix.csv` which are newline delimited lists of words, generates a list of all permutations, and queries AWS route53domain to see which ones are available for registration. Domains available for registration are printed to `stdout`.

realistically speaking it's probably simpler just generating a list and then pasting it into a registrar that supports looking up a list in bulk - but this was slightly more fun.
