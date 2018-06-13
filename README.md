# Buzz
A Simple CLI tool to calculate an account's payment volume (USD) over time and output results.

## Output
Currently, there are two output options availible.
### Print to Screen
Quick and easy... see the magic right before your eyes.
![terminal output](https://imgur.com/crks1Nu.png)

### Save to CSV
Export as CSV and let your excel sheet wizards do their magic.
![imported to sheets](https://imgur.com/WpqQV6C.png)

## Setup
First, you will need a [Currencylayer](https://currencylayer.com/) api key.
Once you get this, add it as an environment variable:
```
CURRENCY_LAYER_API_KEY
```
Then you have two options to install this:

1. clone this repo, `cd buzz`, `go install`, `dep ensure` (if you don't have dep, `brew install dep`)
2. `go get github.com/robertdurst/buzz`

Finally, query some account payment history!

Currently there is only one function:
```
buzz query [stellar_address] [output_csv_filename] [flags]
```

Flags
```
--aggregate [day | month | none] // Time interval aggregating of data [DEFAULT: none]
--output [csv | terminal] // Print to screen or save to csv [DEFAULT: terminal]
```

## Current State
* Only supports XLM and fiat based tokens

## Future Work
* Use an alternative currency, such as EUR, instead of USD (**requires paid subscription**)

## Wish List
* Fully functional web app
* Programmatically generate graphs
* Advanaced statistical analysis