# Overview

This tool will identify the size of all directories recursively and then print which directories are distinctly larger or smaller.


# Rationale

I wrote this tool to help me identify which directories of images have not yet been processed (collections of images relating to specific events or topics that have already been processed will have been moved to a different tree). I wanted to identify directories that had more than the others without using static thresholds. I might want to check one month at a time, and some months are busier than others.


# Design

It simply calculates the mean, the standard deviation, and determines which directories are more than one standard deviations away.


# Install

To get:

```
$ go get github.com/dsoprea/go-find-dense-paths/command/go-find-dense-paths
```

To build:

```
$ go build github.com/dsoprea/go-find-dense-paths/command/go-find-dense-paths
```


# Usage

There is full command-line help. You must only provide a path. By defaut, you'll get a list of polarity, path, and count, where polarity represents whether that path is less than or more than the mean by more than one standard deviation. You may pass --big or --small to only return those results (in this case, polarity will be omitted). You may pass --sigma to use a different standard deviation threshold (as an integer). You may pass --json to print JSON-formatted results instead of CSV-formatted results.

Directories are scanned recursively. Parent counts do not include the counts of their children. Results are always sorted by path.


# Example

```
$ go-find-dense-albums --path ./ --big
20181214_Detroit_Canada_NewYork_Nantucket-Phone,2202
20181214_Detroit_Canada_NewYork_Nantucket-Sony,1871
20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10480201,1056
20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10580202,2742
20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10680203,992
20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10780204,1730
20190211_Bogota/Sony/DCIM/10290210,1179
20190211_Bogota/Sony/DCIM/10390208,1075
```

```
$ go-find-dense-albums --path ./ --big --json
[
  {
    "path": "20181214_Detroit_Canada_NewYork_Nantucket-Phone",
    "count": 2202
  },
  {
    "path": "20181214_Detroit_Canada_NewYork_Nantucket-Sony",
    "count": 1871
  },
  {
    "path": "20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10480201",
    "count": 1056
  },
  {
    "path": "20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10580202",
    "count": 2742
  },
  {
    "path": "20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10680203",
    "count": 992
  },
  {
    "path": "20190131_Lisbon_Brussels_Netherlands_Lux_Germany/Images/Sony_RX100-VI/DCIM/10780204",
    "count": 1730
  },
  {
    "path": "20190211_Bogota/Sony/DCIM/10290210",
    "count": 1179
  },
  {
    "path": "20190211_Bogota/Sony/DCIM/10390208",
    "count": 1075
  },
]
```
