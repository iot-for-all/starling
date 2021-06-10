
## Development Environment Setup ##

### Golang ###
Follow the instructions to [install Go](https://golang.org/doc/install). Pick the appropriate package to install the
latest 1.15.x release of Go. This will give you access to the Go toolchain and compiler.

- If you are on Windows, use the MSI to install. It will set the necessary environment variables.
- If you installed via the tarball, you will need to add a GOROOT environment variable pointing to the
  folder where you installed Go (typically /usr/local/go on linux-based systems)
- You should also check to make sure that you can access the Go compiler and tools. They are available at $GOROOT/bin
  (or $GOROOT\bin) and should be added to your path if they are not already. You can verify this by running the following:
    - Max/Linux: `which go`
    - Windows (CMD): `where go`

### Prometheus ###
Starling simulator metrics are available from the [http://localhost:6002/metrics](http://localhost:6002/metrics) endpoint.
Promethus can be configured to scrape these metrics periodically (15 secs, by default) and store in its timeseries
database locally for analysis.

1. __Prometheus Install:__ Download [Prometheus](https://prometheus.io/download/) and unzip into a folder.
2. __Configure Prometheus:__ Copy over the [setup/prometheus.yml](setup/prometheus.yml) file into the install folder.
   You can run the prometheus.exe executable in the install folder. By default prometheus is avaialble
   at [http://localhost:9090](http://localhost:9090)

### Grafana ###
Timeseries metrics from Prometheus can be analyzed in a graphical dashboard tool called Grafana.

1. __Grafana Install:__ Download [Grafana](https://grafana.com/grafana/download) and unzip into a folder.
   Run the `Grafana_Install_Dir/bin/grafana-server.exe` file. Grafana can be accessed from
   [http://localhost:3000/](http://localhost:3000/). Once you can access the site, login admin/admin or admin/ [no password].
   ![picture alt](assets/grafana-setup1.png "Grafana Setup")

2. __Data Source Setup:__ Add data source (second tile on the homepage), select Prometheus with the following parameters:
   ![picture alt](assets/grafana-setup2.png "Grafana Setup")
   ![picture alt](assets/grafana-setup3.png "Grafana Setup")

Parameter  | Value
-----------|-------------
Name       | Prometheus
Url        | http://localhost:9090
Access     | Server (default)

Click _Save & Test_ button.
![picture alt](assets/grafana-setup4.png "Grafana Setup")

3. __Dashboard Setup:__ Select  `+ -> Import` menu item on the left navbar and select `Upload JSON File`.
   Select [setup/grafana-dashboard.json](../setup/grafana-dashboard.json) file. Once you start running the Starling server,
   Prometheus and Grafana will come alive.
   ![picture alt](assets/grafana-setup5.png "Grafana Setup")
   ![picture alt](assets/grafana-setup6.png "Grafana Setup")
   ![picture alt](assets/grafana-setup7.png "Grafana Setup")
   You should see the following dashboard. Once we start the Starling server, this will start showing simulation data.
   ![picture alt](assets/grafana-setup8.png "Grafana Setup")


[Back to contents](../Readme.md)| Previous: [Configuring and running simulations](configure.md) | Next: [FAQ - Frequently asked questions](faq.md)
---------------------------------|-------------------------------------------------------|------------------------------------
