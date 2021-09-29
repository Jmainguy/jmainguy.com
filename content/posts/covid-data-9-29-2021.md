---
title: "Covid Data as of 9-29-2021"
date: 2021-09-29T13:33:06-04:00
draft: false
---

I recently created a [stock price operator for Kubernetes](https://github.com/Jmainguy/stockop), to gather the current price of stock for a list of tickers the user defines. Doing so, I learned the same [API](https://finnhub.io/docs/api/) I was using for the stock, also had some Covid-19 data in it as well.

The code was pretty quick to work out thanks to the great API docs and examples.
```/bin/bash
package main

import (
	"context"
	"fmt"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func main() {
        cfg := finnhub.NewConfiguration()
        cfg.AddDefaultHeader("X-Finnhub-Token", "c50beiaad3ic9bdl2tf0")
        finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

        // Covid-19
        covid19, _, _ := finnhubClient.Covid19(context.Background()).Execute()
        fmt.Println("State, TotalCases, TotalDeaths, MortalityRate, Updated")
        for _, covidInfo := range covid19 {
                percentageOFCaseToDeath := *covidInfo.Death / *covidInfo.Case * 100
                fmt.Printf(
                        "%+v, %.0f, %.0f, %.2f, %+v\n",
                        *covidInfo.State, *covidInfo.Case, *covidInfo.Death,
                        percentageOFCaseToDeath, *covidInfo.Updated,
                )
        }
}
```

Using the code above, I could get the data into a csv base, which I could then manipulate easily with the `sort` command.

```/bin/bash
$ ./convid | grep -v 'Washington, D.C' | sort -t, -k4 -h | tail -5
New Jersey, 1150246, 27380, 2.38, 2021-09-29 19:02:03
Grand Princess, 103, 3, 2.91, 2021-04-20 17:00:47
Navajo Nation, 30366, 1262, 4.16, 2021-04-20 17:00:47
Veteran Affair, 250417, 11558, 4.62, 2021-04-20 17:00:47
Northern Mariana Islands, 31, 2, 6.45, 2021-04-20 17:00:47
```

Here we can see that the Northern Mariana Islands (I have no idea where this is, but I believe they have a trench) has the highest rate, however the sample size is too small to really take it any futher.

The next two highest groups are the VA and the Navajo Nation, pretty interesting to see such a high mortality rate for the Navajo and Veterans, would love to understand why this is.

It is cool the Navajo Nation has data like this, why dont more first nation tribes have this kind of data? Something for me to look into another time.

Ultimately, I think I would need more data sources, ones that update more often than Finnhub here, also would like to be able to see a states total population, vaccine counts, and range of data overtime, to see the valleys and peaks in the data.
