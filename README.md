# SolarPal 
Simple CRUD app for interacting with an [NREL API](https://developer.nrel.gov/docs/solar/pvwatts/v8/) to get PV power estimates using user provided parameters about their custom solar panel array.

>Please note that the SQLite driver used in this project, namely [go-sqlite3](https://github.com/mattn/go-sqlite3/tree/master), requires GCC to build for the first time and having CGo enabled. The steps to get this done can be found in the link above or for simply downloading the TDM-GCC toolchain for example you can go to this [link](https://jmeubank.github.io/tdm-gcc/) also provided in the docs for the driver.
