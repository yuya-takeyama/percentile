# percentime

Executes a command n times and shows percentile of input numbers.

This Project is based on
- [yuya-takeyama/ntimes](https://github.com/yuya-takeyama/ntimes)
- [yuya-takeyama/percentile](https://github.com/yuya-takeyama/percentile)

## Installation

```
$ go get github.com/fluktuid/percentime
```

## Usage

``` bash
$ percentime 100 -- curl -s -o /dev/null -w "%{time_total}" google.com
50%:	0.061504
66%:	0.063227
75%:	0.064522
80%:	0.065158
90%:	0.066899
95%:	0.07377
98%:	0.07851
99%:	0.082946
100%:	0.096252
```

## License

The MIT License
