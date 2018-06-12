# Buzz
A Simple CLI tool to calculate an account's payment volume (USD) over time. 

## Setup
First, you will need a [Currencylayer](https://currencylayer.com/) api key.
Once you get this, add it as an environment variable:
```
CURRENCY_LAYER_API_KEY
```
Then you have two options to install this:

1. clone this repo, `cd buzz`, `go install`, 'dep ensure' (if you don't have dep, `brew install dep`)
2. `go get github.com/robertdurst/buzz`

Finally, query some account payment history!

Currently there is only one function:
```
buzz query [stellar_address] [output_csv_filename]
```

## Current State
* Only option is to output a CSV
* Only supports XLM and fiat based tokens
* No data processing

## Future Work
* Post process data, aggregating data by day/month/year
* Use an alternative currency, such as EUR, instead of USD

## Wish List
* Fully functional web app
* Programmatically generate graphs
* Advanaced statistical analysis